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
	"crowdfunding-api/campaign"
	"strings"
	"crowdfunding-api/helper"
	"net/http"
	"github.com/dgrijalva/jwt-go"
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
	campaignRepository := campaign.NewRepository(db.DB)
	userService := user.NewService(userRepository)	
	authService := auth.NewService(cfg)
	campaignService := campaign.NewService(campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	
	router := gin.Default()
	router.Static("/images/","./images")
	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)

	err = router.Run()

	if err != nil {
		panic(err)
	}
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""

		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString=arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

