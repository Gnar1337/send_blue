package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Gnar1337/send_blue/tree/main/proto"
)

type server struct {
	pb.UnimplementedMessageServiceServer
}

// SendAndTrack implements the gRPC service
func (s *server) SendAndTrack(req *pb.Message, stream pb.MessageService_SendAndTrackServer) error {
	log.Printf("Incoming iMessage Request [UID: %s] for Lead: %s", req.Uid, req.Lead)

	// 1. Immediate Status: QUEUED
	// (We send this back as soon as the request hits the server)
	if err := stream.Send(&pb.StatusUpdate{Uid: req.Uid, Status: "QUEUED"}); err != nil {
		return err
	}

	// Define the remaining lifecycle
	lifecycle := []string{"ACCEPTED", "SENT", "DELIVERED", "RECEIVED"}

	for _, status := range lifecycle {
		// Simulate processing time/network lag
		time.Sleep(2 * time.Second)

		log.Printf("[%s] Updating Status: %s", req.Uid, status)

		// 2. PUSH the update to the client
		err := stream.Send(&pb.StatusUpdate{
			Uid:    req.Uid,
			Status: status,
		})

		if err != nil {
			log.Printf("Stream error (Client probably disconnected): %v", err)
			return err // Ending the function closes the stream
		}
	}

	log.Printf("[%s] Lifecycle Complete", req.Uid)
	return nil // Closing the stream successfully
}

func main() {
	// Start listening on the standard gRPC port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, &server{})

	fmt.Println("🚀 iMessage Gateway (gRPC Server) running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
