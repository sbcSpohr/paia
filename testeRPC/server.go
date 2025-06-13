package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Player struct {
	X, Y int
}

type GameServer struct {
	mu      sync.Mutex
	players map[string]Player
}

type RegisterArgs struct {
	Name string
}

type MoveArgs struct {
	Name  string
	DeltaX int
	DeltaY int
}

type StateReply struct {
	Players map[string]Player
}

func (s *GameServer) RegisterPlayer(args RegisterArgs, reply *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.players[args.Name]; !exists {
		s.players[args.Name] = Player{X: 0, Y: 0}
		*reply = "Player registered!"
	} else {
		*reply = "Player already exists."
	}
	return nil
}

func (s *GameServer) Move(args MoveArgs, reply *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, exists := s.players[args.Name]; exists {
		p.X += args.DeltaX
		p.Y += args.DeltaY
		s.players[args.Name] = p
		*reply = fmt.Sprintf("Moved to (%d, %d)", p.X, p.Y)
	} else {
		*reply = "Player not found."
	}
	return nil
}

func (s *GameServer) GetState(_ struct{}, reply *StateReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	reply.Players = s.players
	return nil
}

func main() {
	server := &GameServer{
		players: make(map[string]Player),
	}
	rpc.Register(server)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	fmt.Println("Servidor RPC escutando na porta 1234...")
	rpc.Accept(listener)
}
