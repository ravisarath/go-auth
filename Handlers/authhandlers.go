package Handlers

import (
	"fmt"
	"jwt-todo/auth-server/Config"
	"jwt-todo/auth-server/Models"
	"math/rand"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessId     string
	RefreshId    string
	AtExpires    int64
	RtExpires    int64
}
type Todo struct {
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

type AccessDetails struct {
	AccessId string
	UserId   string
}
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateAccount(c *gin.Context) {
	var todo Models.Account
	c.BindJSON(&todo)

	repr := Models.CreateAccount(&todo)

	if !repr {
		c.JSON(http.StatusUnauthorized, "Username or Email already excist")
		return
	} else {
		c.JSON(http.StatusOK, todo)
		return
	}

}
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	resp, account := Models.Login(u.Username, u.Password)
	userId := account.ID.Hex()
	if !resp {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := CreateToken(userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := CreateAuth(userId, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	userdetails := map[string]string{
		"id":       userId,
		"username": account.Username,
		"group":    account.Group,
		"company":  account.Company,
	}
	err = Models.UpdateToken(userId, ts.AccessToken, ts.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Something went wrong")
	}
	err = Models.UpdateLoginTime(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Something went wrong")
	}

	tokenConfig := Config.BuildTokenConfig()
	c.SetCookie("acces_token", ts.AccessToken, tokenConfig.MaxAge,
		tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)
	c.SetCookie("refresh_token", ts.RefreshToken, tokenConfig.MaxAge,
		tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)

	c.JSON(
		http.StatusOK, gin.H{"data": userdetails})

}

func CreateTodo(c *gin.Context) {
	var td *Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	userId, err := FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	td.UserID = userId

	c.JSON(http.StatusCreated, td)
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	keyConfig := Config.BuildKeyConfig()
	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(keyConfig.RfSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		refreshId, ok := claims["refreshId"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, ok := claims["userId"].(string)

		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}

		deleted, delErr := DeleteAuth(refreshId)
		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		ts, createErr := CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}

		err1 := Models.RefreshUpdateToken(userId, ts.AccessToken, ts.RefreshToken)
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, "Something went wrong")
		}

		tokenConfig := Config.BuildTokenConfig()
		c.SetCookie("acces_token", ts.AccessToken, tokenConfig.MaxAge,
			tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)
		c.SetCookie("refresh_token", ts.RefreshToken, tokenConfig.MaxAge,
			tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)
		c.JSON(http.StatusCreated, "Token Refreshed")
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
func Logout(c *gin.Context) {
	au, err := ExtractTokenMetadata(c.Request)
	tokenConfig := Config.BuildTokenConfig()
	c.SetCookie("acces_token", "", tokenConfig.MaxAge,
		tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)
	c.SetCookie("refresh_token", "", tokenConfig.MaxAge,
		tokenConfig.Path, tokenConfig.Domain, tokenConfig.Secure, tokenConfig.HttpOnly)
	c.Set("is_logged_in", false)
	if err != nil {

		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	deleted, delErr := DeleteAuth(au.AccessId)
	if delErr != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	err = Models.RemoveToken(au.UserId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unable to update")
	}
	err = Models.UpdateLogoutTime(au.UserId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Something went wrong")
	}

	c.JSON(http.StatusOK, "Successfully logged out")
}
