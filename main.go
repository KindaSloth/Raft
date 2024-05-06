package main

import (
	context "context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct{}

type Config struct {
	// randomized to be between 150ms and 300ms
	ElectionTimeout  int
	HeartbeatTimeout int
	// Servers address
	Servers []string
}

type Node struct {
	Id     int
	Config Config
	// Follower | Candidate | Leader
	State string
	Logs  []*Log
	// current node term
	Term int
	// how many votes node has received
	VoteCount int
	// pointer to node id
	VotedFor   *int
	GrpcServer *grpc.Server
}

// Node State Machine (get/run current state)
func (node *Node) Run() {
	currentState := fmt.Sprintf("Node ID: %s, Current State: %s", strconv.Itoa(node.Id), node.State)
	fmt.Println(currentState)
	time.Sleep(5 * time.Second)
}

func (s *server) AppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	log.Printf("Received AppendEntries request")
	return &AppendEntriesResponse{Term: 1, Success: true}, nil
}

func (s *server) RequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	log.Printf("Received RequestVote request")
	return &RequestVoteResponse{Term: 1, VoteGranted: false}, nil
}

func (s *server) mustEmbedUnimplementedRaftServer() {}

// just testing servers
func (node *Node) SendMessageTest() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := NewRaftClient(conn)

	res, err := client.AppendEntries(context.Background(), &AppendEntriesRequest{LeaderTerm: 1, LeaderId: 1, Entries: node.Logs})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	log.Printf("Response from server: %s", res.String())
}

func startNodes(x int) []Node {
	var servers []string
	var logs []*Log

	config := Config{
		ElectionTimeout:  150,
		HeartbeatTimeout: 300,
		Servers:          servers,
	}

	var nodes []Node

	for i := 1; i <= x; i++ {
		port := ":808" + strconv.Itoa(i)
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("Listening on %s", port)

		s := grpc.NewServer()
		RegisterRaftServer(s, &server{})
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		servers = append(servers, port)

		node := Node{
			Id:         i,
			Config:     config,
			State:      "Follower",
			Logs:       logs,
			Term:       0,
			VoteCount:  0,
			VotedFor:   nil,
			GrpcServer: s,
		}

		nodes = append(nodes, node)
	}

	return nodes
}

func main() {
	nodes := startNodes(4)

	for _, node := range nodes {
		go func(node Node) {
			for {
				node.Run()
			}
		}(node)
	}

	select {}
}
