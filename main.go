package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

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
	Logs  []Log
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
}

func startNodes(x int) []Node {
	var servers []string
	var logs []Log

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

	for {
		for _, node := range nodes {
			node.Run()
		}
	}
}
