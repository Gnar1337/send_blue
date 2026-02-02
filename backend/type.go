package main

import (
	"database/sql"
	"sync"
	"time"
)

// ///////////////// QUEUE //////////////////
type MessageQueue struct {
	items    []MessageQueueItem
	mu       sync.Mutex
	sendTime time.Duration
}

// Enqueue adds an item to the end
func (q *MessageQueue) Enqueue(item MessageQueueItem) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *MessageQueue) ChangeSendTime(newSendTime time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.sendTime = newSendTime * time.Second
}

// Dequeue removes and returns the front item
func (q *MessageQueue) Dequeue() (MessageQueueItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return MessageQueueItem{}, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// ///////////////// DB TYPES ///////////////////
type Client struct {
	ID              string             `gorm:"type:uuid;primaryKey;column:uid" json:"uid"`
	Name            string             `json:"name"`
	MessagesSent    int                `json:"messagesSent"`
	Leads           []Lead             `gorm:"foreignKey:ClientUID;references:ID" json:"leads"`
	MessagesInQueue []MessageQueueItem `gorm:"foreignKey:FromClientID;references:ID" json:"messageQueue"`
	AllMessagesSent []MessageQueueItem `gorm:"foreignKey:FromClientID;references:ID" json:"allMessagesSent"`
}
type Lead struct {
	LeadNumber       string       `json:"leadNumber"`
	ClientUID        string       `gorm:"type:uuid;column:client_uid" json:"clientUid"`
	MessagesReceived int          `json:"messagesReceived"`
	LastContacted    sql.NullTime `json:"lastContacted"`
}

type MessageQueueItem struct {
	MsgUID              string                `gorm:"type:uuid;primaryKey;column:msg_uid" json:"uid"`
	MessageBody         string                `gorm:"column:message_body" json:"messageBody"`
	FromClientID        string                `gorm:"type:uuid;column:from_client_id" json:"fromClientId"`
	ToClientLead        string                `gorm:"column:to_client_lead" json:"toClientLead"`
	ScheduledSendTime   time.Time             `gorm:"column:scheduled_send_time" json:"scheduledSendTime"`
	TimeSent            sql.NullTime          `gorm:"column:time_sent" json:"timeSent"`
	Status              string                `gorm:"column:status" json:"status"`
	Archived            bool                  `gorm:"column:archived" json:"archived"`
	MessageEventHistory []MessageEventHistory `gorm:"foreignKey:MsgUID;references:MsgUID" json:"messageEventHistory"`
}

type MessageEventHistory struct {
	MsgUID     string    `gorm:"type:uuid;column:msg_uid" json:"msgUid"`
	TimeStamp  time.Time `gorm:"column:event_time" json:"time_stamp"`
	PrevStatus string    `gorm:"column:prev_status" json:"prev_status"`
	CurrStatus string    `gorm:"column:curr_status" json:"curr_status"`
}
