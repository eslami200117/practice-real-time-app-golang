package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"rest.gtld.test/realTimeApp/app/model"
	"rest.gtld.test/realTimeApp/app/usecases"
)


type userHandler struct {
	userUsecaseImp *usecases.UserUsecaseImp
}

func NewUserHanlder(userUsecase * usecases.UserUsecaseImp) *userHandler {
	return &userHandler{
		userUsecaseImp: userUsecase,
	}
}

func (u *userHandler) HandleLogin(c *gin.Context) {
	var json model.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if ok := u.userUsecaseImp.AuthenticateUser(c, &json); ok {
	
		generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":  json.Username,
			"exp": time.Now().Add(time.Hour * 2).Unix(),
		})
	
		token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
		h := sha256.New()
		h.Write([]byte(token))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
		usecases.LoginJWT[signature] = token
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
		}
		c.JSON(http.StatusOK, gin.H{"token": signature})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
	}
	u.userUsecaseImp.UpdateLastLogin(json.Username, time.Now())
}

func (u *userHandler) GetCurrenct(username string, user *model.Login){
	u.userUsecaseImp.GetLoginUser(username, user)
}

func (u *userHandler) UsersListHandler(c *gin.Context) {
	listOfUsers := u.userUsecaseImp.GetAllUser()
	c.JSON(http.StatusOK, gin.H{
		"users": listOfUsers,
	})
}

func (u *userHandler) AddUserHandler(c *gin.Context) {
	var addUser struct {
		Username string
		AddUser	string
		Password string
	}

	if err := c.ShouldBindJSON(&addUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminUsername := addUser.Username

	if u.userUsecaseImp.IsAdmin(adminUsername) {

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not admin"})
		return
	}

	err := u.userUsecaseImp.AddUser(addUser.AddUser, addUser.Password)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	}

}