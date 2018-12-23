package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "coordinator",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {

			// Do Stuff Here
		},
	}
)

var (
	PostgresPort  int
	RaftPort      int
	JoinAddr      string
	DataDirectory string
)

func init() {
	rootCmd.Flags().IntVar(&PostgresPort, "postgres-port", 5433, "Port that will accept PostgreSQL connections.")
	rootCmd.Flags().IntVar(&RaftPort, "raft-port", 5434, "Port that will be used for Noah's raft protocol.")
	rootCmd.Flags().StringVar(&DataDirectory, "data-dir", "data", "Directory that will be used for Noah's key value store.")
	rootCmd.Flags().StringVar(&JoinAddr, "join", "", "The address of another node in a cluster to join.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
