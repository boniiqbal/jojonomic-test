package main

import (
	"flag"
	"sync"

	dotenv "github.com/joho/godotenv"
	"github.com/topup-storage-service/config"
	"github.com/topup-storage-service/src"
)

func main() {
	flag.Parse()

	err := dotenv.Load("./config/.env")
	if err != nil {
		panic(".env is not loaded properly")
	}

	cfg := config.NewConfig()
	server := src.InitServer(cfg)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run()
	}()

	// Wait All services to end
	wg.Wait()

}
