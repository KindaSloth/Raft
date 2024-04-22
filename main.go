package main

import "fmt"

type Config struct {
	// randomized to be between 150ms and 300ms
	ElectionTimeout  int64
	HeartbeatTimeout int64
	Servers          []string
}

type Node struct {
	Id     int64
	Config Config
	// Follower | Candidate | Leader
	State string
	Logs  []Log
	// current node term
	Term int64
	// how many votes node has received
	VoteCount int64
	// pointer to node id
	VotedFor int64
}

func main() {
	fmt.Println("Raft")
}
