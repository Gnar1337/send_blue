package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"send-blue-backend/handlers"
	"send-blue-backend/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var MessageQueue *types.MessageQueue
var DB *handlers.DBConn

func main() {
	// Init db
	DB = handlers.InitDB()

	//init queue
	interval, err := strconv.Atoi(os.Getenv("QUEUE_INTERVAL")) // Load .env file
	if err != nil {
		log.Fatal("Failed to parse QUEUE_INTERVAL environment variable:", err)
	}
	MessageQueue = types.InitQueue(time.Duration(interval), DB.Conn)

	// Initialize Gin router
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	//Make ENDPOINTS
	r.GET("/clients", DB.GetClients())
	r.GET("/clients/leads", DB.GetClientLeads())
	r.GET("/clients/scheduled", DB.GetClientsQueue())
	r.GET("/client/data", DB.ClientGetData())
	r.POST("/client/schedule_message", DB.ScheduleMessage(MessageQueue))
	//interval change
	r.GET("/gateway/interval", DB.SetMessageQueueInterval(MessageQueue))

	go populateLeads(DB.Conn)
	go delayDumpQueue(DB.Conn)
	// Start server
	r.Run(":8080")
}

// //////////////// UTILITIES ////////////////////////
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
		time.Sleep(MessageQueue.GetSendTime())
		sendMessagesFromQueue(db)
	}
}

// Sends messages from the queue
func sendMessagesFromQueue(db *gorm.DB) {
	for {
		item, ok := MessageQueue.Dequeue()
		if !ok {
			return
		}
		go SendiMessageAndListen(item, db)
	}
}
