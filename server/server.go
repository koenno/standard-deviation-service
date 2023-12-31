package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/koenno/standard-deviation-service/service"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

//go:generate mockery --name=RandomIntegerGenerator --case underscore --with-expecter
type RandomIntegerGenerator interface {
	Integers(ctx context.Context, quantity int) ([]int, error)
}

//go:generate mockery --name=StdDevCalculator --case underscore --with-expecter
type StdDevCalculator interface {
	Calculate(input <-chan []int) <-chan service.StdDevResult
}

type RandomServer struct {
	srv        http.Server
	generator  RandomIntegerGenerator
	calculator StdDevCalculator
	port       int
}

func NewRandomServer(generator RandomIntegerGenerator, calculator StdDevCalculator, port int) *RandomServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(validationMiddleware)

	s := &RandomServer{
		srv: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
		generator:  generator,
		calculator: calculator,
		port:       port,
	}

	r.Route("/random", func(r chi.Router) {
		r.Use(validationMiddleware)
		r.Get("/mean", s.Mean)
	})

	return s
}

func (s *RandomServer) Run() {
	slog.Info("server is running", "timestamp", time.Now(), "port", s.port)
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "timestamp", time.Now(), "error", err)
	}
}

func (s *RandomServer) Stop() {
	err := s.srv.Shutdown(context.Background())
	if err != nil {
		slog.Error("server shutting down error", "timestamp", time.Now(), "error", err)
	}
}

func (s *RandomServer) Mean(w http.ResponseWriter, r *http.Request) {
	requests, _ := paramPositiveInt(r, "requests")
	length, _ := paramPositiveInt(r, "length")

	res, err := s.doMean(r.Context(), requests, length)
	if err != nil {
		slog.Error("mean calculation", "timestamp", time.Now(), "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("failed to encode the payload", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *RandomServer) doMean(ctx context.Context, requests, numbers int) ([]service.StdDevResult, error) {
	pipe := make(chan []int, requests)

	resultPipe := s.calculator.Calculate(pipe)

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < requests; i++ {
		g.Go(func() error {
			randomInts, err := s.generator.Integers(ctx, numbers)
			if err != nil {
				return err
			}
			pipe <- randomInts
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate standard deviation: %v", err)
	}
	close(pipe)

	var res []service.StdDevResult
	for singleRes := range resultPipe {
		res = append(res, singleRes)
	}
	return res, nil
}
