package db

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fpt-corp/fptshop/text"
)

type Promotion struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Url         string `db:"url" json:"url"`
	Image       string `db:"image" json:"image_url"`
	Description string `db:"description" json:"description"`
	SentCount   int    `db:"sent_count" json:"sent_count"`
	ReadCount   int    `db:"read_count" json:"read_count"`
}

type PromotionSend struct {
	IDs         string `db:"promotion_ids"`
	ID          int    `db:"promotion_id"`
	RecipientID string `db:"recipient_id"`
	SentAt      string `db:"sent_at"`
	ReadAt      string `db:"read_at"`
	ClickedAt   string `db:"clicked_at"`
}

type PromotionSchedule struct {
	ID          int    `db:"schedule_id"`
	PromotionID string `db:"promotion_id"`
	Message     string `db:"message"`
	DoneAt      string `db:"done_at"`
}

type PromotionCriteria struct {
	Title string
	Ids   string
	Id    int
}

func (p Promotion) GetTitle() string       { return p.Title }
func (p Promotion) GetUrl() string         { return p.Url }
func (p Promotion) GetImage() string       { return p.Image }
func (p Promotion) GetDescription() string { return p.Description }

func QueryPromotion(criteria PromotionCriteria) ([]Promotion, error) {
	query := `SELECT p.id, p.title, p.url, p.image, p.description, COUNT(ps.sent_at) AS sent_count, COUNT(ps.read_at) AS read_count
		FROM promotion p
		LEFT JOIN promotion_send ps ON p.id = ps.promotion_id
		WHERE is_active = 1`
	if criteria.Title != "" {
		criteria.Title = "%" + criteria.Title + "%"
		query += ` AND title like :title`
	}
	if criteria.Ids != "" {
		err := validateIds(criteria.Ids)
		if err != nil {
			return nil, err
		}
		query += fmt.Sprintf(` AND id IN (%s)`, criteria.Ids)
	}
	if criteria.Id > 0 {
		query += ` AND id = :id`
	}
	query += ` GROUP BY p.id`

	rows, err := connection.NamedQuery(query, criteria)
	if err != nil {
		return nil, err
	}

	// Initialize using 'make' instead of 'var' for getting '[]' instead of 'null' in json.Marshal
	promotions := make([]Promotion, 0)
	for rows.Next() {
		var promotion Promotion
		if err = rows.StructScan(&promotion); err != nil {
			log.Error(err)
			continue
		}
		promotions = append(promotions, promotion)
	}
	return promotions, nil
}

func NewPromotion(promotion *Promotion) error {
	query := `INSERT INTO promotion (title, url, image, description, is_active)
	 		VALUES (:title, :url, :image, :description, 1)`
	_, err := connection.NamedExec(query, promotion)
	return err
}

func DeletePromotion(ids string) error {
	err := validateIds(ids)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE promotion SET is_active = 0 WHERE id IN (%s)", ids)
	_, err = connection.Exec(query)
	return err
}

func UpdatePromotion(promotion *Promotion) error {
	query := `UPDATE promotion SET title = :title, url = :url, image = :image, description = :description WHERE id = :id`
	_, err := connection.NamedExec(query, promotion)
	return err
}

func NewPromotionSchedule(promotionIds []int, senderFilter string, message string) error {
	query := "INSERT INTO promotion_schedule (promotion_id, message) VALUES (?, ?)"
	result, err := connection.Exec(query, text.FromIntSlice(promotionIds), message)
	if err != nil {
		return err
	}
	scheduleId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	var recipients []Sender
	if senderFilter == "*" {
		var err error
		recipients, err = QuerySender()
		if err != nil {
			return err
		}
	}
	query = "INSERT IGNORE INTO promotion_send (promotion_id, recipient_id, schedule_id) VALUES "
	count := 0
	var placeHolders []string
	var params []interface{}
	for _, promotionId := range promotionIds {
		for _, recipient := range recipients {
			count++
			placeHolders = append(placeHolders, "(?, ?, ?)")
			params = append(params, promotionId, recipient.ID, scheduleId)
			if count > 100 { // Insert 100 rows at a time
				_, err := connection.Exec(query+strings.Join(placeHolders, ","), params...)
				if err != nil {
					return err
				}
				count = 0
				params, placeHolders = nil, nil
			}
		}
	}
	if len(placeHolders) > 0 { // Insert left-over rows
		_, err := connection.Exec(query+strings.Join(placeHolders, ","), params...)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetActivePromotionSchedule() (PromotionSchedule, error) {
	query := "SELECT schedule_id, promotion_id, message FROM promotion_schedule WHERE done_at IS NULL LIMIT 1"
	schedules := []PromotionSchedule{}

	err := connection.Select(&schedules, query)
	if err != nil || len(schedules) == 0 {
		return PromotionSchedule{}, err
	}
	return schedules[0], nil
}

func QueryPromotionsToSend(scheduleId int, batchSize int) ([]PromotionSend, error) {
	query := `SELECT recipient_id, GROUP_CONCAT(promotion_id) AS 'promotion_ids'
			FROM promotion_send
			WHERE schedule_id = ? AND sent_at IS NULL
			GROUP BY recipient_id LIMIT ?`
	promotionsToSend := []PromotionSend{}
	err := connection.Select(&promotionsToSend, query, scheduleId, batchSize)
	return promotionsToSend, err
}

func UpdatePromotionSent(promotionIds string, recipient string) error {
	// Internal call so don't need to use place holders here
	query := fmt.Sprintf(`UPDATE promotion_send SET sent_at = NOW()
			WHERE promotion_id IN (%s) AND recipient_id = '%s'`, promotionIds, recipient)
	_, err := connection.Exec(query)
	return err
}

func UpdatePromotionRead(recipient string, time time.Time) error {
	query := `UPDATE promotion_send SET read_at = CONVERT_TZ(?, '+00:00', 'SYSTEM')
			WHERE sent_at IS NOT NULL AND read_at IS NULL AND recipient_id = ?`
	timeStr := time.Format("2006-01-02 15:04:05")
	_, err := connection.Exec(query, timeStr, recipient)
	return err
}

func UpdatePromotionScheduleDone(scheduleId int) error {
	query := `UPDATE promotion_schedule SET done_at = NOW() WHERE schedule_id = ?`
	_, err := connection.Exec(query, scheduleId)
	return err
}
