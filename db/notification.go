package db

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// Table: notification
type ProductNotification struct {
	ProductID   string `db:"product_id" json:"id"`
	ProductName string `db:"product_name" json:"product_name"`
	SKU			string `db:"sku" json:"sku"`
	Url 		string `db:"url" json:"url"`
	ImageUrl 	string `db:"image_url" json:"image_url"`
	SubscriberCount int `db:"subscriber_count" json:"subscriptions"`
}

// Table: notification_subscribers
type Subscriber struct {
	ProductID   string `db:"product_id"`
	SenderID    string `db:"sender_id"`
	SenderName  string `db:"sender_name"`
}

func SubscribeProductNotification(p ProductNotification, s Subscriber) error {
	tx, err := connection.Beginx()
	if err != nil {
		return err
	}

	query := `INSERT IGNORE INTO notification (product_id, product_name, sku, url, image_url)
			VALUES (:product_id, :product_name, :sku, :url, :image_url);`
	_, err = tx.NamedExec(query, p)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `INSERT INTO notification_subscribers (product_id, sender_id, sender_name)
			VALUES (:product_id, :sender_id, :sender_name)
			ON DUPLICATE KEY UPDATE subscribed_at = NOW();`
	_, err = tx.NamedExec(query, s)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func QueryProductNotifications() ([]ProductNotification, error)  {
	query := `SELECT n.product_id, n.product_name, n.url, n.image_url, n.sku, count(s.sender_id) as subscriber_count
			FROM notification n JOIN notification_subscribers s ON n.product_id = s.product_id
			WHERE n.is_active = 1
			GROUP BY n.product_id, n.product_name, n.url, n.image_url, n.sku`
	notifications := []ProductNotification{}
	err := connection.Select(&notifications, query)
	return notifications, err
}

func GetProductNotification(id int) (ProductNotification, error) {
	query := `SELECT product_id, product_name, url, image_url, sku
			FROM notification
			WHERE product_id = ?`
	notification := ProductNotification{}
	err := connection.Get(&notification, query, id)
	return notification, err
}

func DeleteProductNotifications(ids string) error {
	err := validateIds(ids)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE notification SET is_active = 0 WHERE product_id IN (%s)", ids)
	_, err = connection.Exec(query)
	return err
}

func QuerySubscribers(productId int) ([]Subscriber, error) {
	query := `SELECT sender_id, sender_name FROM notification_subscribers
			WHERE product_id = ?`
	rows, err := connection.Queryx(query, productId)
	if err != nil {
		return nil, err
	}

	var subscribers []Subscriber
	for rows.Next() {
		subscriber := Subscriber{}
		if err = rows.StructScan(&subscriber); err != nil {
			log.Error(err)
			continue
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}
