package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/robber-account/pkg/client"
	"github.com/eviltomorrow/robber-account/pkg/pb"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create acccount service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			log.Fatalf("[Error] Missing args for create, eg. [seq: nick_name, email, phone]\r\n")
		}

		var (
			nickName = args[0]
			email    = args[1]
			phone    = args[2]
		)
		stub, close, err := client.NewClientForAccount()
		if err != nil {
			log.Fatalf("[Error] Create account grpc client failure, nest error: %v\r\n", err)
		}
		defer close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		reply, err := stub.Create(ctx, &pb.User{
			NickName: nickName,
			Email:    email,
			Phone:    phone,
		})
		if err != nil {
			log.Fatalf("[Error] Create account failure, nest error: %v\r\n", err)
		}
		fmt.Printf("Create account success, uuid: %v\r\n", reply.Value)
	},
}

func init() {
	createCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.toml", "robber-account's config file")

	rootCmd.AddCommand(createCmd)
}
