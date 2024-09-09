package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mapper",
	Short: "Mapper is a stream mapping service",
	Long:  `Mapper is a service that maps streams to their origin servers and provides an API for querying this information.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
