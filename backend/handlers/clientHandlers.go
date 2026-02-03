package handlers

import (
	"fmt"
	"log"
	"send-blue-backend/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConn type holds the db connection
type DBConn struct {
	Conn *gorm.DB
}

// InitDB starts db connction
func InitDB() *DBConn {
	// Init db

	dsn := "host=db user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.Exec("SELECT 1") // Simple query to verify connection
	return &DBConn{
		Conn: db,
	}
}

// GetClients gets the list of clients efoqr the ui
func (db *DBConn) GetClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		var clients []types.Client
		rows, err := db.Conn.Raw("SELECT uid::text, name FROM clients").Rows()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var client types.Client
			if err := rows.Scan(&client.ID, &client.Name); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			clients = append(clients, client)
		}
		c.JSON(200, gin.H{
			"clients": clients,
		})
	}
}

// GetClientLeads gets the list of a client's leads
func (db *DBConn) GetClientLeads() gin.HandlerFunc {
	return func(c *gin.Context) {
		var leads []types.Lead
		clientId := c.Query("client_id")
		rows, err := db.Conn.Raw("SELECT lead_number FROM client_leads WHERE client_uid = $1", clientId).Rows()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var lead types.Lead
			if err := rows.Scan(&lead.LeadNumber); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			leads = append(leads, lead)
		}
		c.JSON(200, gin.H{
			"leads": leads,
		})
	}
}

// GetClientsQueue gets the list of a client's queued messages
func (db *DBConn) GetClientsQueue() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msgs []types.MessageQueueItem
		fromClientId := c.Query("client_id")
		// _ = db.Table("clients").Find(&clients)
		rows, err := db.Conn.Raw("SELECT msg_uid::text, message_body, from_client_id::text, to_client_lead, scheduled_send_time, time_sent, status FROM message_queue WHERE status = 'QUEUED' AND from_client_id::text = $1 AND NOT archived ORDER BY scheduled_send_time ASC", fromClientId).Rows()
		if err != nil {

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var msg types.MessageQueueItem
			if err := rows.Scan(&msg.MsgUID, &msg.MessageBody, &msg.FromClientID, &msg.ToClientLead, &msg.ScheduledSendTime, &msg.TimeSent, &msg.Status); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			msgs = append(msgs, msg)
		}
		c.JSON(200, gin.H{
			"messages": msgs,
		})
	}
}

// GetMessageHistory gets the list of a message's event history
func (db *DBConn) GetMessageHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msgHistory []types.MessageEventHistory
		msgUid := c.Query("msg_uid")
		// _ = db.Table("clients").Find(&clients)
		rows, err := db.Conn.Raw("SELECT msg_uid::text, curr_status, prev_status, event_time FROM message_event_history WHERE msg_uid::text = $1 ORDER BY event_time ASC", msgUid).Rows()
		if err != nil {

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var msg types.MessageEventHistory
			if err := rows.Scan(&msg.MsgUID, &msg.CurrStatus, &msg.PrevStatus, &msg.TimeStamp); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			msgHistory = append(msgHistory, msg)
		}
		c.JSON(200, gin.H{
			"messagesHistory": msgHistory,
		})
	}
}

// ClientGetData gets the messages sent and queued for the given client
func (db *DBConn) ClientGetData() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId := c.Query("client_id")
		client := types.Client{}
		query := `SELECT lead_number, client_uid, messages_received, last_contacted 
					FROM client_leads 
					WHERE client_uid = $1;`
		err := db.Conn.Raw(query, clientId).Scan(&client.Leads).Error
		if err != nil {
			log.Println("Error scanning clientleads:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		query = `SELECT msg_uid, message_body, from_client_id, to_client_lead, 
       				scheduled_send_time, time_sent, status, archived 
					FROM message_queue 
					WHERE from_client_id = $1 AND status = 'QUEUED';   `
		err = db.Conn.Raw(query, clientId).Scan(&client.MessagesInQueue).Error
		if err != nil {
			log.Println("Error scanning messagequeue:", err)
		}
		query = `SELECT msg_uid, message_body, from_client_id, to_client_lead, 
					scheduled_send_time, time_sent, status, archived 
					FROM message_queue 
					WHERE from_client_id = $1 AND status != 'QUEUED';`
		err = db.Conn.Raw(query, clientId).Scan(&client.AllMessagesSent).Error
		if err != nil {
			log.Println("Error scanning allmessagessent:", err)
		}
		client.ID = clientId
		err = db.Conn.Raw(`SELECT name FROM clients WHERE uid::text = $1`, clientId).Scan(&client.Name).Error
		if err != nil {
			log.Println("Error scanning name:", err)
		}
		fmt.Println("executed query")
		c.JSON(200, gin.H{"client": client})
	}
}

// ScheduleMessage schedules a messages for the given client
func (db *DBConn) ScheduleMessage(messageQueue *types.MessageQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msgToQueue types.MessageQueueItem
		if err := c.BindJSON(&msgToQueue); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Insert into DB and get the generated UUID
		result := db.Conn.Raw("INSERT INTO message_queue (message_body, from_client_id, to_client_lead, scheduled_send_time, status, archived) VALUES ($1, $2, $3, $4, $5, $6) RETURNING msg_uid",
			msgToQueue.MessageBody,
			msgToQueue.FromClientID,
			msgToQueue.ToClientLead,
			msgToQueue.ScheduledSendTime,
			msgToQueue.Status,
			false,
		).Scan(&msgToQueue.MsgUID)

		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}
		// Enqueue
		messageQueue.Enqueue(msgToQueue)
		c.JSON(201, gin.H{"status": "received", "data": msgToQueue})
	}
}

// SetMessageQueueInterval sets the messageQueues send interval
func (db *DBConn) SetMessageQueueInterval(messageQueue *types.MessageQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		seconds := c.Query("seconds")
		interval, err := strconv.Atoi(seconds)
		if err != nil {
			log.Fatal("Failed to parse QUEUE_INTERVAL environment variable:", err)
		}
		fmt.Println("Queue interval set to:", interval, "seconds")
		messageQueue.ChangeSendTime(time.Duration(interval) * time.Second)
		c.JSON(200, gin.H{"status": "received", "data": interval})
	}
}
