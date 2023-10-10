package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koenno/standard-deviation-service/client"
	"github.com/koenno/standard-deviation-service/client/randomorg"
	"github.com/koenno/standard-deviation-service/random"
	"github.com/koenno/standard-deviation-service/server"
	"github.com/koenno/standard-deviation-service/service"
	"golang.org/x/time/rate"
)

func main() {
	reqsPerSec := flag.Int("reqs", 10, "number of requests per second")
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	rateLimiter := rate.NewLimiter(rate.Every(time.Second), *reqsPerSec)
	reqSender := client.New(rateLimiter)
	respParser := randomorg.NewBodyParser()
	reqFactory := randomorg.NewRequestFactory()

	generator := random.NewRandom(reqSender, respParser, reqFactory)

	calculator := service.NewStdDevService()

	srv := server.NewRandomServer(generator, calculator, *port)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		srv.Stop()
	}()

	srv.Run()
}
