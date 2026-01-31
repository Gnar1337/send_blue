package main

import (
	"context"
	"io"
	"log"

	// This matches the module name in your proto/go.mod
	pb "github.com/Gnar1337/send_blue/proto"
	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendiMessageAndListen(m MessageQueueItem, db *gorm.DB) {
	// 1. Establish connection to the Gateway (Simulator)
	// In Docker, "gateway" matches the service name in docker-compose.yml
	conn, err := grpc.NewClient("gateway:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessageServiceClient(conn)

	// 2. Prepare the Message payload
	// msg := &pb.Message{
	// 	Uid:     "uuid-12345",
	// 	Body:    "Testing the iMessage Status Pipeline",
	// 	Message: "Metadata here",
	// 	Client:  "Internal-App",
	// 	Lead:    "John Doe",
	// 	Status:  "QUEUED",
	// }
	msg := &pb.Message{
		Uid:     m.MsgUID,
		Body:    m.MessageBody,
		Message: "Metadata here",
		Client:  m.FromClientID,
		Lead:    m.ToClientLead,
		Status:  "QUEUED",
	}

	// 3. Call the streaming function
	// We use context.Background() but in prod, you'd use a timeout context
	stream, err := client.SendAndTrack(context.Background(), msg)
	if err != nil {
		log.Fatalf("Could not start stream: %v", err)
	}

	log.Printf("Message %s initiated. Waiting for lifecycle updates...", msg.Uid)

	// 4. LISTEN for updates until the stream is closed
	for {
		update, err := stream.Recv()

		// io.EOF means the server (Gateway) has finished sending all 5 statuses
		if err == io.EOF {
			log.Println("✅ All updates received. Lifecycle complete.")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream update: %v", err)
		}

		// This is where you would update your Database (Postgres/Redis)
		// with the new status (ACCEPTED, SENT, etc.)
		log.Printf("📨 [STATUS UPDATE] Message ID: %s | New Status: %s", update.Uid, update.Status)
		if update.Status == "RECEIVED" {
			db.Exec("UPDATE message_queue SET status=$1, archived=$3 WHERE msg_uid::text=$2",
				update.Status,
				update.Uid,
				true,
			)
		} else if update.Status == "SENT" {
			db.Exec("UPDATE message_queue SET status=$1, time_sent=NOW() WHERE msg_uid::text=$2",
				update.Status,
				update.Uid,
			)
		} else {
			db.Exec("UPDATE message_queue SET status=$1 WHERE msg_uid::text=$2",
				update.Status,
				update.Uid,
			)
		}
	}
}
