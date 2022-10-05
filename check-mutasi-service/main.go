package main

import (
	"flag"
	"sync"

	"github.com/check-mutasi-service/config"
	"github.com/check-mutasi-service/src"
	dotenv "github.com/joho/godotenv"
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
