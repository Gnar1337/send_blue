package main

import (
	"database/sql"
	"sync"
	"time"
)

/////////////////// QUEUE //////////////////

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

// func main() {
//     q := &SafeQueue{}
//     q.Enqueue("Data A")
//     q.Enqueue("Data B")

//     val, ok := q.Dequeue()
//     if ok {
//         fmt.Println("Dequeued:", val) // Output: Data A
//     }
// }
// type MessageQueue []MessageQueueItem

// func (mq *MessageQueue) Enqueue(item MessageQueueItem) {
// 	*mq = append(*mq, item)
// }
// func (mq *MessageQueue) Dequeue() (MessageQueueItem, bool) {
// 	if len(*mq) == 0 {
// 		return MessageQueueItem{}, false
// 	}
// 	item := (*mq)[0]
// 	*mq = (*mq)[1:]
// 	return item, true
// }

// ///////////////// DB TYPES ///////////////////
type Client struct {
	ID           string `gorm:"type:uuid;primaryKey;column:uid" json:"uid"`
	Name         string `json:"name"`
	MessagesSent int    `json:"messagesSent"`
}
type Lead struct {
	LeadNumber       string       `json:"leadNumber"`
	ClientUID        string       `gorm:"type:uuid;column:client_uid" json:"clientUid"`
	MessagesReceived int          `json:"messagesReceived"`
	LastContacted    sql.NullTime `json:"lastContacted"`
}

type MessageQueueItem struct {
	MsgUID            string       `gorm:"type:uuid;primaryKey;column:msg_uid" json:"uid"`
	MessageBody       string       `gorm:"column:message_body" json:"messageBody"`
	FromClientID      string       `gorm:"type:uuid;column:from_client_id" json:"fromClientId"`
	ToClientLead      string       `gorm:"column:to_client_lead" json:"toClientLead"`
	ScheduledSendTime time.Time    `gorm:"column:scheduled_send_time" json:"scheduledSendTime"`
	TimeSent          sql.NullTime `gorm:"column:time_sent" json:"timeSent"`
	Status            string       `gorm:"column:status" json:"status"`
	Archived          bool         `gorm:"column:archived" json:"archived"`
}
