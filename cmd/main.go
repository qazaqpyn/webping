package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qazaqpyn/webping"
	"github.com/qazaqpyn/webping/domain/websites"
	handler "github.com/qazaqpyn/webping/internal/handler/http"
	"github.com/qazaqpyn/webping/internal/repository"
	"github.com/qazaqpyn/webping/internal/repository/mongo"
	"github.com/qazaqpyn/webping/internal/service"
	"github.com/qazaqpyn/webping/pkg/workerpool"
	"github.com/spf13/viper"
)

const (
	INTERVAL        = time.Second * 1
	REQUEST_TIMEOUT = time.Second * 10
	WORKER_COUNT    = 10
)

func main() {
	// init config
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
		return
	}

	// initilize websites list
	web, err := websites.NewWebsites()
	if err != nil {
		log.Fatal(err)
		return
	}
	urls := web.GetWebsites()

	results := make(chan workerpool.Result)
	workerPool := workerpool.NewPool(WORKER_COUNT, REQUEST_TIMEOUT, results)

	// run workers that range over jobs channel
	workerPool.Init()

	// generate jobs
	go generateJobs(workerPool, urls)

	//processResults runnnig in background
	go processResults(results)

	// server
	db, err := mongo.NewMongodb(viper.GetString("mongo.name"))
	if err != nil {
		log.Fatal(err)
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	hadlers := handler.NewHandler(services, web)

	srv := new(webping.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), hadlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	// workerPool.Stop()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}

func generateJobs(workerPool *workerpool.Pool, urls []string) {
	for {
		for _, url := range urls {
			workerPool.Push(workerpool.Job{URL: url})
		}
		time.Sleep(INTERVAL)
	}
}

func processResults(results chan workerpool.Result) {
	go func() {
		for results := range results {

			if results.Error != nil {
				log.Println(results.Error)
				continue
			}
		}
	}()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
