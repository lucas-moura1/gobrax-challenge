package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucas-moura1/gobrax-challenge/config"
	"github.com/lucas-moura1/gobrax-challenge/entity"
	"github.com/lucas-moura1/gobrax-challenge/handler"
	"github.com/lucas-moura1/gobrax-challenge/repository"
	"github.com/lucas-moura1/gobrax-challenge/usecase"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viper.AutomaticEnv()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log := logger.Sugar()

	db, err := config.LoadDatabase()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Driver{}, &entity.Vehicle{})

	driverRepository := repository.NewDriverRepository(log, db)
	driverUsecase := usecase.NewDriverUsecase(log, driverRepository)
	driverHandler := handler.DriverHandler{
		DriverUsecase: driverUsecase,
	}

	http.HandleFunc("GET /drivers", driverHandler.GetAll)
	http.HandleFunc("GET /drivers/{id}", driverHandler.GetById)
	http.HandleFunc("POST /drivers", driverHandler.Create)
	http.HandleFunc("POST /drivers/{id}/vehicle", driverHandler.AddVehicle)
	http.HandleFunc("PUT /drivers/{id}", driverHandler.Update)
	http.HandleFunc("DELETE /drivers/{id}", driverHandler.Delete)

	vehicleRepository := repository.NewVehicleRepository(log, db)
	vehicleUsecase := usecase.NewVehicleUsecase(vehicleRepository)
	vehicleHandler := handler.VehicleHandler{
		VehicleUsecase: vehicleUsecase,
	}

	http.HandleFunc("GET /vehicles", vehicleHandler.GetAll)
	http.HandleFunc("GET /vehicles/{id}", vehicleHandler.GetById)
	http.HandleFunc("PUT /vehicles/{id}", vehicleHandler.Update)
	http.HandleFunc("DELETE /vehicles/{id}", vehicleHandler.Delete)

	server := &http.Server{Addr: fmt.Sprintf(":%s", viper.GetString("PORT"))}

	go func() {
		log.Info("Server started at :8080")
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Info("Server gracefully stopped!")
}
