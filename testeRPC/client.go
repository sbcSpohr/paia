package main

import (
	"fmt"
	"net/rpc"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}

	var reply string
	err = client.Call("GameServer.RegisterPlayer", RegisterArgs{Name: "breno"}, &reply)
	fmt.Println("->", reply)

	for i := 0; i < 5; i++ {
		err = client.Call("GameServer.Move", MoveArgs{Name: "breno", DeltaX: 1, DeltaY: 0}, &reply)
		fmt.Println("->", reply)

		var state StateReply
		err = client.Call("GameServer.GetState", struct{}{}, &state)
		fmt.Println("Estado atual dos jogadores:", state.Players)

		time.Sleep(1 * time.Second)
	}
}

type RegisterArgs struct {
	Name string
}

type MoveArgs struct {
	Name   string
	DeltaX int
	DeltaY int
}

type Player struct {
	X, Y int
}

type StateReply struct {
	Players map[string]Player
}
