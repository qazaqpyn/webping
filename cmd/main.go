package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qazaqpyn/webping"
	"github.com/qazaqpyn/webping/domain/websites"
	handler "github.com/qazaqpyn/webping/internal/handler/http"
	"github.com/qazaqpyn/webping/internal/repository"
	"github.com/qazaqpyn/webping/internal/repository/mongo"
	"github.com/qazaqpyn/webping/internal/service"
	"github.com/qazaqpyn/webping/pkg/ping"
	"github.com/spf13/viper"
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

	// server
	db, err := mongo.NewMongodb(viper.GetString("mongo.name"))
	if err != nil {
		log.Fatal(err)
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	hadlers := handler.NewHandler(services, web)

	// start pinging
	ping := ping.NewPingHandler(services.Results)

	ping.StartPing()
	go ping.GenerateJobs(urls)
	go ping.ProcessResults()

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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
