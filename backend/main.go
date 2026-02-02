package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var messageQueue *MessageQueue
var db *gorm.DB

func main() {
	// Init db
	dsn := "host=db user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.Exec("SELECT 1") // Simple query to verify connection
	// // END Init db

	//init queue
	interval, err := strconv.Atoi(os.Getenv("QUEUE_INTERVAL")) // Load .env file
	if err != nil {
		log.Fatal("Failed to parse QUEUE_INTERVAL environment variable:", err)
	}
	messageQueue = &MessageQueue{
		items:    make([]MessageQueueItem, 0),
		sendTime: time.Duration(interval) * time.Second,
	}
	// END init queue

	// Initialize Gin router
	r := gin.Default()
	// Configure CORS for http://localhost:5173
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	// Optional: Allow credentials (cookies, auth headers)
	config.AllowCredentials = true
	r.Use(cors.New(config))

	//ENDPOINTS
	// GET Clients endpoint
	r.GET("/clients", func(c *gin.Context) {
		var clients []Client
		rows, err := db.Raw("SELECT uid::text, name FROM clients").Rows()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var client Client
			if err := rows.Scan(&client.ID, &client.Name); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			clients = append(clients, client)
		}
		c.JSON(200, gin.H{
			"clients": clients,
		})
	})

	r.GET("/clients/leads", func(c *gin.Context) {
		var leads []Lead
		clientId := c.Query("client_id")
		rows, err := db.Raw("SELECT lead_number FROM client_leads WHERE client_uid = $1", clientId).Rows()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var lead Lead
			if err := rows.Scan(&lead.LeadNumber); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			leads = append(leads, lead)
		}
		c.JSON(200, gin.H{
			"leads": leads,
		})
	})
	r.GET("/clients/scheduled", func(c *gin.Context) {
		var msgs []MessageQueueItem
		fromClientId := c.Query("client_id")
		// _ = db.Table("clients").Find(&clients)
		rows, err := db.Raw("SELECT msg_uid::text, message_body, from_client_id::text, to_client_lead, scheduled_send_time, time_sent, status FROM message_queue WHERE status = 'QUEUED' AND from_client_id::text = $1 AND NOT archived ORDER BY scheduled_send_time ASC", fromClientId).Rows()
		if err != nil {

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var msg MessageQueueItem
			if err := rows.Scan(&msg.MsgUID, &msg.MessageBody, &msg.FromClientID, &msg.ToClientLead, &msg.ScheduledSendTime, &msg.TimeSent, &msg.Status); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			msgs = append(msgs, msg)
		}
		c.JSON(200, gin.H{
			"messages": msgs,
		})
	})

	r.GET("/client/data", func(c *gin.Context) {
		clientId := c.Query("client_id")
		client := Client{}
		// leads := []Lead{}
		query := `SELECT lead_number, client_uid, messages_received, last_contacted 
					FROM client_leads 
					WHERE client_uid = $1;`
		err := db.Raw(query, clientId).Scan(&client.Leads).Error
		if err != nil {
			log.Println("Error scanning clientleads:", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		query = `SELECT msg_uid, message_body, from_client_id, to_client_lead, 
       				scheduled_send_time, time_sent, status, archived 
					FROM message_queue 
					WHERE from_client_id = $1 AND status = 'QUEUED';   `
		err = db.Raw(query, clientId).Scan(&client.MessagesInQueue).Error
		if err != nil {
			log.Println("Error scanning messagequeue:", err)
		}
		query = `SELECT msg_uid, message_body, from_client_id, to_client_lead, 
					scheduled_send_time, time_sent, status, archived 
					FROM message_queue 
					WHERE from_client_id = $1 AND status != 'QUEUED';`
		err = db.Raw(query, clientId).Scan(&client.AllMessagesSent).Error
		if err != nil {
			log.Println("Error scanning allmessagessent:", err)
		}
		client.ID = clientId
		err = db.Raw(`SELECT name FROM clients WHERE uid::text = $1`, clientId).Scan(&client.Name).Error
		if err != nil {
			log.Println("Error scanning name:", err)
		}
		fmt.Println("executed query")
		c.JSON(200, gin.H{"client": client})
	})

	// POST endpoint
	r.POST("/client/schedule_message", func(c *gin.Context) {
		var msgToQueue MessageQueueItem
		if err := c.BindJSON(&msgToQueue); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Insert into DB and get the generated UUID
		result := db.Raw("INSERT INTO message_queue (message_body, from_client_id, to_client_lead, scheduled_send_time, status, archived) VALUES ($1, $2, $3, $4, $5, $6) RETURNING msg_uid",
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
	})
	r.GET("/gateway/interval", func(c *gin.Context) {
		seconds := c.Query("seconds")
		interval, err := strconv.Atoi(seconds)
		if err != nil {
			log.Fatal("Failed to parse QUEUE_INTERVAL environment variable:", err)
		}
		fmt.Println("Queue interval set to:", interval, "seconds")
		messageQueue.ChangeSendTime(time.Duration(interval))
		c.JSON(200, gin.H{"status": "received", "data": messageQueue.sendTime})
	})

	go populateLeads(db)
	go delayDumpQueue(db)
	// Start server
	r.Run(":8080")
}

func populateLeads(db *gorm.DB) {
	rows, err := db.Raw("SELECT uid::text FROM clients").Rows()
	if err != nil {
		log.Println("Error getting client rows:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var clientId string
		if err := rows.Scan(&clientId); err != nil {
			log.Println("Error scanning client ID:", err)
			continue
		}
		for i := 1; i <= 5; i++ {
			num := 3373230000 + i
			leadNumber := "+1" + fmt.Sprint(num)
			result := db.Exec("INSERT INTO client_leads (lead_number, client_uid) VALUES ($1, $2)", leadNumber, clientId)
			if result.Error != nil {
				log.Println("Error inserting lead:", result.Error)
			}
		}
	}
}

// Delays the sends all messages
func delayDumpQueue(db *gorm.DB) {
	for {
		time.Sleep(messageQueue.sendTime)
		sendMessagesFromQueue(db)
	}
}

// Sends messages from the queue
func sendMessagesFromQueue(db *gorm.DB) {
	for {
		item, ok := messageQueue.Dequeue()
		if !ok {
			return
		}
		go SendiMessageAndListen(item, db)
	}
}
