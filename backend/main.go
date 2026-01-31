package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var messageQueue *MessageQueue

func main() {
	// Init db
	dsn := "host=db user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.Exec("SELECT 1") // Simple query to verify connection
	// END Init db

	//init queue
	messageQueue = &MessageQueue{
		items: make([]MessageQueueItem, 0),
	}
	// _ = messageQueue
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
		// _ = db.Table("clients").Find(&clients)
		rows, err := db.Raw("SELECT uid::text, name FROM clients").Rows()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var client Client
			if err := rows.Scan(&client.ID, &client.Name); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			clients = append(clients, client)
		}

		// result := db.Table("clients").Find(&clients)
		c.JSON(200, gin.H{
			"clients": clients,
		})
	})

	r.GET("/clients/leads", func(c *gin.Context) {
		var leads []Lead
		clientId := c.Query("client_id")
		// _ = db.Table("clients").Find(&clients)
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

		// result := db.Table("clients").Find(&clients)
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
			fmt.Println(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var msg MessageQueueItem
			if err := rows.Scan(&msg.MsgUID, &msg.MessageBody, &msg.FromClientID, &msg.ToClientLead, &msg.ScheduledSendTime, &msg.TimeSent, &msg.Status); err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			msgs = append(msgs, msg)
		}

		// result := db.Table("clients").Find(&clients)
		c.JSON(200, gin.H{
			"messages": msgs,
		})
	})

	// POST endpoint
	r.POST("/client/schedule_message", func(c *gin.Context) {
		var msgToQueue MessageQueueItem
		if err := c.BindJSON(&msgToQueue); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// Insert into DB
		result := db.Exec("INSERT INTO message_queue (message_body, from_client_id, to_client_lead, scheduled_send_time, status, archived) VALUES ($1, $2, $3, $4, $5, $6)",
			msgToQueue.MessageBody,
			msgToQueue.FromClientID,
			msgToQueue.ToClientLead,
			msgToQueue.ScheduledSendTime,
			msgToQueue.Status,
			false,
		)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}
		// Enqueue
		messageQueue.Enqueue(msgToQueue)
		c.JSON(201, gin.H{"status": "received", "data": msgToQueue})
	})
	go populateLeads(db)
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
