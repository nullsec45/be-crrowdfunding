package main

import (
	// "log"
	"crowdfunding-api/config"
	"fmt"
	"crowdfunding-api/user"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"

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


	userInput := user.RegisterUserInput{}
	// userInput.Name="Test"
	// userInput.Email="fajar@gmail.com"
	// userInput.Occupation="Anak Band"
	// userInput.Password="password"
	userService.RegisterUser(userInput)
}