package main

import (
	// "log"
	"crowdfunding-api/config"
	"fmt"
	"crowdfunding-api/user"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"github.com/gin-gonic/gin"
	"crowdfunding-api/handler"

)
var cfgFile string

func main(){
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}else{
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	cfg := config.NewConfig()
	db, err := cfg.ConnectionMysql()

	if err != nil {
		log.Error().Err(err).Msgf("[Main-cfg.ConnectionMySQL] Failed to connect to database : %v", err)
		return
	}

	userRepository := user.NewRepository(db.DB)
	userService := user.NewService(userRepository)	




	userHandler := handler.NewUserHandler(userService)
	
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	router.Run()
}