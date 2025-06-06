package main

import (
	"context"
	"fmt"

	"github.com/jurgen-kluft/go-home/sensor-server/hollywood/actor"
)

type message struct {
	data string
}

type foo struct{}

func newFoo() actor.Receiver {
	return &foo{}
}

func (f *foo) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		fmt.Println("actor started")
	case actor.Stopped:
		fmt.Println("actor stopped")
	case *message:
		fmt.Println("actor has received", msg.data)
	}
}

func main() {
	engine, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	pid := engine.Spawn(newFoo, "my_actor")
	for i := 0; i < 5; i++ {
		engine.Send(pid, &message{data: "hello world!"})
	}

	ctx := engine.PoisonCtx(context.Background(), pid)
	<-ctx.Done()
}
