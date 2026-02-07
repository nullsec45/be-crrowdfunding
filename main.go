package main

import (
	"crowdfunding-api/config"
	"fmt"
	"crowdfunding-api/user"
	"crowdfunding-api/auth"
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
	authService := auth.NewService(cfg)

	userHandler := handler.NewUserHandler(userService, authService)
	
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
}