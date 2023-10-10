package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/koenno/standard-deviation-service/client"
	"github.com/koenno/standard-deviation-service/client/randomorg"
	"github.com/koenno/standard-deviation-service/random"
	"github.com/koenno/standard-deviation-service/server"
	"github.com/koenno/standard-deviation-service/service"
)

func main() {
	reqSender := client.New()
	respParser := randomorg.NewBodyParser()
	reqFactory := randomorg.NewRequestFactory()

	generator := random.NewRandom(reqSender, respParser, reqFactory)

	calculator := service.NewStdDevService()

	srv := server.NewRandomServer(generator, calculator)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		srv.Stop()
	}()

	srv.Run()
}
