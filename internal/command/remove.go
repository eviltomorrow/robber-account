package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/eviltomorrow/robber-account/pkg/client"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove acccount service",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("[Error] Missing args for remove, eg. [uuid]\r\n")
		}

		stub, close, err := client.NewClientForAccount()
		if err != nil {
			log.Fatalf("[Error] Create account grpc client failure, nest error: %v\r\n", err)
		}
		defer close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = stub.Remove(ctx, &wrapperspb.StringValue{
			Value: args[0],
		})
		if err != nil {
			log.Fatalf("[Error] Remove account failure, nest error: %v\r\n", err)
		}
		fmt.Printf("Remove account success, uuid: %v\r\n", args[0])
	},
}

func init() {
	removeCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.toml", "robber-account's config file")

	rootCmd.AddCommand(removeCmd)
}
