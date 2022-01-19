package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lemon-mint/envaddr"
	"github.com/lemon-mint/godotenv"
	"github.com/valyala/fasthttp"
)

func main() {
	godotenv.Load()

	ln, err := net.Listen("tcp", envaddr.Get(":16702"))
	if err != nil {
		log.Fatalln(err)
	}

	server := fasthttp.Server{}
	server.Name = "presigned"
	server.ReadTimeout = time.Minute
	server.WriteTimeout = time.Minute
	server.IdleTimeout = time.Minute

	SignalChan := make(chan os.Signal, 1)
	signal.Notify(SignalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Serve(ln); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	<-SignalChan
	log.Println("Shutting down the server...")
	err = server.Shutdown()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Server gracefully stopped")
}
