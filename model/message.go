package model

import "time"

type MessageType int

type Message struct {
	Id       int64       `json:"id"`
	UserId   int64       `json:"user_id"`
	Sender   string      `json:"sender"`
	Type     MessageType `json:"type"`
	Content  string      `json:"content"`
	SendTime time.Time   `json:"send_time"`
	Expire   int64       `json:"expire"`
}
