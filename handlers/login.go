package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/udpp/utils"
	"net/http"
	"time"
)

type credential struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	var c credential
	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		utils.JSONResponse(1, "bad login request,valid fields: user and password", nil, w)
		logrus.Errorf("bad login request,valid fields: user and password")
		return
	}

	//TODO
	if c.User != "admin" || c.Password != "admin" {
		utils.JSONResponse(1, "invalid user or password", nil, w)
		logrus.Errorf("invalid user or password")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["iss"] = "udpp"
	claims["info"] = struct {
		Name string
	}{c.User}
	//TODO
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		utils.JSONResponse(1, "generate token error", nil, w)
		logrus.Errorf("generate token error")
		return
	}
	utils.JSONResponse(0, "OK", tokenString, w)

}
