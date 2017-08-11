package main

import (
	"os"
	"log"
	"github.com/diablowu/nbridge/command"
)

func main() {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)

	app := command.NewApp()
	app.Init(os.Args)
	app.Run()
}
