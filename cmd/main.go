package main

import (
	"os"

	"github.com/MyNameIsWhaaat/algo-learning/config"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/handler"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/server"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil{
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err:= godotenv.Load(); err != nil{
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := config.NewPostgresDB(config.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil{
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil{
		logrus.Fatalf("error occures while running http server: %s", err.Error())
	}
}

func initConfig() error{
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}