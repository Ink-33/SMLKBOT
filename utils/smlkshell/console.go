package smlkshell

import (
	"fmt"
	"log"
	"os"
)

//ConsoleHandler : 处理console invoke
func ConsoleHandler() {
	go func() {
		for {
			var command string
			fmt.Scanln(&command)
			switch command {
			case "exit":
				log.Println("Program stop.")
				os.Exit(0)
			case "status":
				log.Println(GetStatus())
			}
		}
	}()
}
