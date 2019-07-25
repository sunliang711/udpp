package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/sunliang711/udpp/utils"
)

var (
	signKey   []byte
	verifyKey []byte
)

func init() {
	var err error
	signKey, err = ioutil.ReadFile("keys/udpp.rsa")
	if err != nil {
		panic("Read udpp.rsa error")
	}

	verifyKey, err = ioutil.ReadFile("keys/udpp.rsa.pub")
	if err != nil {
		panic("Read udpp.rsa.pub error")
	}
}

func auth(next http.Handler) http.Handler {
	if !viper.GetBool("auth") {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		})

		if err != nil {
			switch err.(type) {
			case *jwt.ValidationError:
				vErr := err.(*jwt.ValidationError)
				switch vErr.Errors {
				case jwt.ValidationErrorExpired:
					utils.JSONResponse(1, "token expired", nil, w)
					logrus.Info("token expired")
					return
				default:
					utils.JSONResponse(1, "parsing token error", nil, w)
					logrus.Info("Parsing token error")
					return
				}
			default:
				utils.JSONResponse(2, "parsing token error", nil, w)
				logrus.Info("Parsing token error")
				return
			}
		}

		if token.Valid {
			next.ServeHTTP(w, req)
		} else {
			utils.JSONResponse(1, "invalid token", nil, w)
			logrus.Info("invalid token")
		}
	})
}
