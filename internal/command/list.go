package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eviltomorrow/robber-account/pkg/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all account",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		stub, close, err := client.NewClientForAccount()
		if err != nil {
			log.Fatalf("[Error] Create account grpc client failure, nest error: %v\r\n", err)
		}
		defer close()

		resp, err := stub.FindAll(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Fatalf("[Error] List all account failure, nest error: %v\r\n", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Nick Name", "Email", "Phone"})
		var count int
		for {
			user, err := resp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("[Error] Recv account failure, nest error: %v\r\n", err)
			}
			var v = []string{user.Uuid, user.NickName, user.Email, user.Phone}
			table.Append(v)
			count++
		}
		if count != 0 {
			table.Render()
		} else {
			fmt.Println("Empty")
		}

	},
}

func init() {
	listCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.toml", "robber-account's config file")

	rootCmd.AddCommand(listCmd)
}
