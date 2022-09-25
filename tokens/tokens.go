package tokens

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var Jwt_Key = []byte("keyXXZZ")

type Claims struct {
	Name  string `json:"username"`
	Email string `json:"useremail"`
	jwt.StandardClaims
}

var (
	expirationTime time.Time
	save_token     string
)

func GenerateToken(nameStr string, emailstr string) (tokenStr string, createdCookies bool) {
	if nameStr == "admin" {
		expirationTime = time.Now().UTC().Add(time.Minute * 60)
	} else {
		expirationTime = time.Now().UTC().Add(time.Minute * 10)
	}

	cliams := &Claims{
		Name:  nameStr,
		Email: emailstr,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}

	my_token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	//tk2, _ := my_token.SigningString()

	token, errIntoken := my_token.SignedString(Jwt_Key)
	if errIntoken != nil {
		fmt.Println("error in signed string ")
		return "", false
	}

	return token, true
}

func ValidateToken(ginCtx *gin.Context) (valueCookie string,
	IsValidCookie bool,
	IsExpiredCookie bool,
	creatorName string,
	creatorEmail string,
	err error) {
	claim := &Claims{}
	readCookieByUser, errCookieByUser := ginCtx.Cookie("user_token")
	if errCookieByUser != nil {
		readCookieByAdm, errCookieByAdm := ginCtx.Cookie("adm_token")
		if errCookieByAdm != nil {
			return "", false, false, "", "", errCookieByUser
		}
		save_token = readCookieByAdm
	}
	if readCookieByUser != "" {
		save_token = readCookieByUser
	}

	tokenKey, errParseTokenKey := jwt.ParseWithClaims(save_token, claim, func(t *jwt.Token) (interface{}, error) {
		return Jwt_Key, nil
	})
	if errParseTokenKey != nil {
		if !tokenKey.Valid {
			return save_token, false, true, "", "", errParseTokenKey
		} else if errParseTokenKey == jwt.ErrInvalidKey {
			return save_token, false, true, "", "", errParseTokenKey
		}
	}
	return save_token, true, false, claim.Name, claim.Email, nil
}
