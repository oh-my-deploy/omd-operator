package main

import (
	cli "github.com/oh-my-deploy/omd-operator/cmd/internal"
	"github.com/oh-my-deploy/omd-operator/internal/utils"
	"log"
	"os"
)

func main() {
	rootCmd := cli.CreateRootCmd()
	utils.RegisterSubCommands(rootCmd, cli.InitVersionCmd, cli.InitGenerateCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
