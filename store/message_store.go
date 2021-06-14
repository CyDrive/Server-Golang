package store

import (
	"github.com/yah01/CyDrive/model"
	"time"
)

type MessageStore interface {
	GetMessageByTime(userId int64, count int64, time time.Time) []model.Message
	PutMessage(message model.Message)
}

type MessageMemStore struct {
	messageMap map[int64][]model.Message
}

func (store MessageMemStore) GetMessageByTime(userId int64, count int64, time time.Time) []model.Message {
	messageList, ok := store.messageMap[userId]
	if ok == false {
		return []model.Message{}
	}
	left := 0
	right := len(messageList) - 1
	if messageList[left].SendTime.After(time) {
		return []model.Message{}
	}
	for left < right {
		mid := (left + right) / 2
		if time.After(messageList[mid+1].SendTime) {
			left = mid + 1
		} else {
			right = mid
		}
	}
	if left-int(count)+1 < 0 {
		return messageList[0:left]
	} else {
		return messageList[left-int(count)+1 : left]
	}
}

func (store MessageMemStore) PutMessage(message model.Message) {
	userId := message.UserId
	_, ok := store.messageMap[userId]
	if ok == false {
		store.messageMap[userId] = []model.Message{}
	}
	store.messageMap[userId] = append(store.messageMap[userId], message)
}

var messageStore MessageStore = MessageMemStore{}

func GetMsgMemStore() MessageStore { return messageStore }
