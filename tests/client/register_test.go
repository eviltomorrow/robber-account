package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/eviltomorrow/robber-account/pkg/client"
	"github.com/eviltomorrow/robber-account/pkg/pb"
)

func TestCreate(t *testing.T) {
	stub, close, err := client.NewClientForAccount()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repley, err := stub.Create(ctx, &pb.User{
		NickName: "shepard",
		Email:    "eviltomorrow@163.com",
		Phone:    "132514628460",
	})
	if err != nil {
		log.Fatalf("Create error: %v", err)
	}
	fmt.Println(repley.Value)
}
