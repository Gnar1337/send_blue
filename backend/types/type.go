package types

import (
	"database/sql"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ///////////////// QUEUE //////////////////
// MessageQueue holds messages to queue
type MessageQueue struct {
	items    []MessageQueueItem
	mu       sync.Mutex
	sendTime time.Duration
}

// InitQueue initializes the queue
func InitQueue(newSendTime time.Duration, db *gorm.DB) *MessageQueue {
	MsgQueue := MessageQueue{
		items:    make([]MessageQueueItem, 0),
		mu:       sync.Mutex{},
		sendTime: newSendTime * time.Second,
	}
	MsgQueue.GetCachedQueue(db)
	return &MsgQueue
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

func (q *MessageQueue) GetSendTime() time.Duration {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.sendTime
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

// GetCachedQueue retrieves queue from DB although in prodw would use Redis
func (q *MessageQueue) GetCachedQueue(db *gorm.DB) bool {
	rows, err := db.Raw("SELECT msg_uid::text, message_body, from_client_id::text, to_client_lead, scheduled_send_time, time_sent, status FROM message_queue WHERE status = 'QUEUED' AND NOT archived ORDER BY scheduled_send_time ASC").Rows()
	if err != nil {
		return false
	}
	for rows.Next() {
		var msg MessageQueueItem
		if err := rows.Scan(&msg.MsgUID, &msg.MessageBody, &msg.FromClientID, &msg.ToClientLead, &msg.ScheduledSendTime, &msg.TimeSent, &msg.Status); err != nil {
			return false
		}
		q.Enqueue(msg)
	}
	return true
}

// ///////////////// DB TYPES ///////////////////
// DBConn holds the connection setup for scaling
type DBConn struct {
	Conn *gorm.DB
}

// Client holds various client data
type Client struct {
	ID              string             `gorm:"type:uuid;primaryKey;column:uid" json:"uid"`
	Name            string             `json:"name"`
	MessagesSent    int                `json:"messagesSent"`
	Leads           []Lead             `gorm:"foreignKey:ClientUID;references:ID" json:"leads"`
	MessagesInQueue []MessageQueueItem `gorm:"foreignKey:FromClientID;references:ID" json:"messageQueue"`
	AllMessagesSent []MessageQueueItem `gorm:"foreignKey:FromClientID;references:ID" json:"allMessagesSent"`
}

// Lead for lead information
type Lead struct {
	LeadNumber       string       `json:"leadNumber"`
	ClientUID        string       `gorm:"type:uuid;column:client_uid" json:"clientUid"`
	MessagesReceived int          `json:"messagesReceived"`
	LastContacted    sql.NullTime `json:"lastContacted"`
}

// MessageQueueItem is the actualmessage data
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

// MessageEventHistory is for tracking a messages history
type MessageEventHistory struct {
	MsgUID     string    `gorm:"type:uuid;column:msg_uid" json:"msgUid"`
	TimeStamp  time.Time `gorm:"column:event_time" json:"time_stamp"`
	PrevStatus string    `gorm:"column:prev_status" json:"prev_status"`
	CurrStatus string    `gorm:"column:curr_status" json:"curr_status"`
}
