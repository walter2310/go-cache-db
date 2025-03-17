package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/walter2310/basic-cache-db/internal/cache"
	"github.com/walter2310/basic-cache-db/internal/commands"
)

func main() {
	c := cache.NewCache()

	err := c.LoadFromDisk("internal/data/cache.json")
	if err != nil {
		fmt.Println("[ERROR] Failed to load cache:", err)
	}

	// Limpieza de llaves expiradas cada 20 segundos
	c.CleanUpExpiredKeys(time.Second * 20)

	// Capturar señales para guardar antes de salir
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		fmt.Println("\n[INFO] Saving cache before exit...")
		err := c.SaveToDisk("internal/data/cache.json")
		if err != nil {
			fmt.Println("[ERROR] Failed to save cache:", err)
		}
		os.Exit(0)
	}()

	// Loop de comandos
	for {
		input, err := commands.GetInput()
		if err != nil {
			log.Fatal(err)
		}

		if input == "cls" {
			break
		}

		commands.ExecuteCommands(input, c)
	}

	fmt.Println("[INFO] Saving cache before closing...")
	err = c.SaveToDisk("internal/data/cache.json")
	if err != nil {
		fmt.Println("[ERROR] Failed to save cache:", err)
	}

	log.Println("Closing...")
}
