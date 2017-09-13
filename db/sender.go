package db

import (
	log "github.com/Sirupsen/logrus"
)

type Sender struct {
	ID       string `db:"sender_id"`
	FullName string `db:"full_name"`
	Gender   string `db:"gender"`
	Phone    string `db:"phone"`
}

func NewSender(sender Sender) error {
	query := `INSERT INTO senders (sender_id, full_name, gender)
	 		VALUES (:sender_id, :full_name, :gender)
			ON DUPLICATE KEY UPDATE last_send = NOW();`
	_, err := connection.NamedExec(query, sender)
	return err
}

func GetSenderByID(id string) (Sender, error) {
	query := `SELECT sender_id, full_name, gender, phone FROM senders WHERE sender_id = ?`
	sender := Sender{}
	err := connection.Get(&sender, query, id)
	return sender, err
}

func UpdateSender(sender Sender) error {
	query := `UPDATE senders SET sender_id = ? `
	params := []interface{}{sender.ID}

	if sender.FullName != "" {
		query += `, full_name = ? `
		params = append(params, sender.FullName)
	}
	if sender.Phone != "" {
		query += `, phone = ? `
		params = append(params, sender.Phone)
	}

	query += ` WHERE sender_id = ?`
	params = append(params, sender.ID)
	_, err := connection.Exec(query, params...)
	return err
}

func QuerySender() ([]Sender, error) {
	query := `SELECT sender_id, full_name, gender FROM senders`
	rows, err := connection.Queryx(query)
	if err != nil {
		return nil, err
	}

	var senders []Sender
	for rows.Next() {
		sender := Sender{}
		if err = rows.StructScan(&sender); err != nil {
			log.Error(err)
			continue
		}
		senders = append(senders, sender)
	}

	return senders, nil
}
