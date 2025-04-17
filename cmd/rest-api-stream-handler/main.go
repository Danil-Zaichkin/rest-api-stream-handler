package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/api"
	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/entity"
	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/repository"
	"github.com/Danil-Zaichkin/rest-api-stream-handler/internal/usecase"
	"github.com/redis/go-redis/v9"
)

type HttpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	memoryRepo := repository.NewInMemoryRepository()
	dbRepo := repository.NewDBRepo(rdb)

	calcUC := usecase.NewCalculatorUsecase(dbRepo, memoryRepo)

	calcHanler := api.New(calcUC)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: calcHanler.Handler(),
	}

	go func() {
		log.Printf("server listening at %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	ctx := context.Background()
	AddShutdownHook(ctx, srv, dbRepo, memoryRepo)
}

type DBRepository interface {
	SaveStreamsContexts(ctx context.Context, streamsCtx map[string]*entity.StreamContext) error
}

type InMemoryRepo interface {
	GetStreamsContexts() map[string]*entity.StreamContext
}

func AddShutdownHook(ctx context.Context, s HttpServer, dbr DBRepository, mr InMemoryRepo) {
	log.Printf("listening signals...")
	c := make(chan os.Signal, 1)
	signal.Notify(
		c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
	)

	<-c

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("can't shutdown server: %v", err)
	}
	fmt.Println(mr.GetStreamsContexts())
	err := dbr.SaveStreamsContexts(ctx, mr.GetStreamsContexts())
	if err != nil {
		log.Fatalf("can't save contexts: %v", err)
	}
}
