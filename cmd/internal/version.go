package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func InitVersionCmd() *cobra.Command {
	return createVersionCmd()
}

func createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of omd-generate-plugin",
		Long:  "All software has versions. This is omd-generate-plugin's",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("omd-generate-plugin of Cli version - LOCAL ")
		},
	}
}
