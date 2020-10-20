package Handlers

import (
	"fmt"
	"jwt-todo/auth-server/Config"
	"jwt-todo/auth-server/Models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

func CreateToken(userid string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessId = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshId = uuid.NewV4().String()
	var err error
	keyConfig := Config.BuildKeyConfig()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["accessId"] = td.AccessId
	atClaims["userId"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(keyConfig.AtSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refreshId"] = td.RefreshId
	rtClaims["userId"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(keyConfig.RfSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessId, ok := claims["accessId"].(string)
		if !ok {
			return nil, err
		}
		userId, ok := claims["userId"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessId: accessId,
			UserId:   userId,
		}, nil
	}
	return nil, err
}
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
func CreateAuth(userid string, td *TokenDetails) error {

	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := Models.Client.Set(td.AccessId, userid, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := Models.Client.Set(td.RefreshId, userid, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
func FetchAuth(authD *AccessDetails) (string, error) {
	userid, err := Models.Client.Get(authD.AccessId).Result()
	if err != nil {
		return "", err
	}

	return userid, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := Models.Client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
