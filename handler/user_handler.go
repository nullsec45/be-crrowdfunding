package handler

import (
	"crowdfunding-api/user"
	"crowdfunding-api/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"crowdfunding-api/helper"
	"fmt"
	"path/filepath"
	"time"
	"mime"
	"os"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler{
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context){
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "Error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "Error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "Error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token);

	response := helper.APIResponse("Account successfully registered", http.StatusOK, "Success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context){
	var input user.LoginInput
	
	err := c.ShouldBindJSON(&input)			

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Login  failed", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors" : err.Error()}

		response := helper.APIResponse("Login  failed", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)

	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "Error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)


	response := helper.APIResponse("Successfully Login", http.StatusOK, "Success", formatter)
	c.JSON(http.StatusOK, response)
}

func  (h *userHandler) CheckEmailAvailability(c *gin.Context){
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors" : "Server error"}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "Error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available" : isEmailAvailable,
	}

	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email is available"	
	}else {
		metaMessage = "Email has been registered"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "Success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
    file, err := c.FormFile("avatar")
    if err != nil {
        data := gin.H{"is_uploaded": false}
        response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
        c.JSON(http.StatusBadRequest, response)
        return
    }

	contentType := file.Header.Get("Content-Type")
	extensions, err := mime.ExtensionsByType(contentType)
    
    if err != nil || len(extensions) == 0 {
        response := helper.APIResponse("Unknown file type", http.StatusBadRequest, "error", nil)
        c.JSON(http.StatusBadRequest, response)
        return
    }
	
	currentUser := c.MustGet("currentUser").(user.User)
    userID := currentUser.ID
	extension := filepath.Ext(file.Filename)


	userExist, err := h.userService.GetUserByID(userID)
    if err != nil {
        response := helper.APIResponse("User not found", http.StatusNotFound, "error", nil)
        c.JSON(http.StatusNotFound, response)
        return
    }

	oldFilePath := userExist.AvatarFileName

	
    path := fmt.Sprintf("images/%d-%d%s", userID, time.Now().Unix(), extension)

    if err := c.SaveUploadedFile(file, path); err != nil {
        data := gin.H{"is_uploaded": false}
        response := helper.APIResponse("Failed to save file", http.StatusBadRequest, "error", data)
        c.JSON(http.StatusBadRequest, response)
        return
    }

   
    _, err = h.userService.SaveAvatar(userID, path)

    if err != nil {
        data := gin.H{"is_uploaded": false}
        response := helper.APIResponse("Failed to update database", http.StatusBadRequest, "error", data)
        c.JSON(http.StatusBadRequest, response)
        return
    }

	if oldFilePath != "" {
		if _, err := os.Stat(oldFilePath); err == nil {
			err := os.Remove(oldFilePath)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
        		return
			}
		}
	}

    data := gin.H{"is_uploaded": true}
    response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
    c.JSON(http.StatusOK, response)
}