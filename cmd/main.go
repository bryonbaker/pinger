package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/bryonbaker/pinger/pkg/protoc"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	//fmt.Printf("Received message: %s\n", req.GetMessage())
	return &pb.PingResponse{}, nil
}

func (s *server) mustEmbedUnimplementedPingerServer() {
	//var _ pb.PingerServer = (*unimplementedPingerServer)(nil)
}

const (
	reportInterval = 10 * time.Second
)

type pingResult struct {
	success bool
}

var hadSuccess = false

func main() {
	listenPort := flag.String("listen-port", "50051", "port to listen on")
	remoteURL := flag.String("remote-url", ":50052", "URL of the remote server")
	siteName := flag.String("site-name", "site-name-undefined", "Name of the site that is running")
	flag.Parse()

	fmt.Printf("Initialising.\n")
	fmt.Printf("\tSite name: %s\n", *siteName)
	fmt.Printf("\tListening on port: %s\n", *listenPort)
	fmt.Printf("\tRemote server url: %s\n", *remoteURL)

	// Start the server
	go func() {
		lis, err := net.Listen("tcp", ":"+*listenPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterPingerServer(s, &server{})

		log.Printf("Server started, listening on :%s\n", *listenPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(500 * time.Millisecond)

	// Start the client
	conn, err := grpc.Dial(*remoteURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPingerClient(conn)

	var successCount uint64 = 0
	var failureCount uint64 = 0
	results := make(chan *pingResult)

	go func() {
		ticker := time.NewTicker(reportInterval)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("Summary: Success: %d, Failure: %d\n", successCount, failureCount)
		}
	}()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		go func() {
			pingMsg := fmt.Sprintf("Hello from: %s", *siteName)

			_, err := c.Ping(ctx, &pb.PingRequest{Message: pingMsg})
			if err != nil {
				log.Println("Ping failed:", err)
				results <- &pingResult{success: false}
			} else {
				hadSuccess = true
				successCount++
				results <- &pingResult{success: true}
			}
		}()

		select {
		case result := <-results:
			if !result.success {
				if hadSuccess {
					failureCount++
				} else {
					fmt.Printf("Waiting for initial synchronisation\n")
				}
			}
		case <-ctx.Done():
			log.Println("Ping timed out")
			if hadSuccess {
				failureCount++
			} else {
				fmt.Printf("Timeout - waiting for initial synchronisation\n")
			}
		}

		time.Sleep(time.Second)
	}
}
