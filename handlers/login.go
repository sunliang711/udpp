package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/udpp/utils"
)

type credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	username = "ann"
	password = "ann123"
)

func healthHandler(w http.ResponseWriter, req *http.Request) {
	utils.JSONResponse(0, "OK", nil, w)
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
	if c.Username != username || c.Password != password {
		utils.JSONResponse(1, "invalid user or password", nil, w)
		logrus.Errorf("invalid user or password")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["iss"] = "udpp"
	claims["info"] = struct {
		Name string
	}{c.Username}
	//TODO
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		utils.JSONResponse(1, "generate token error", nil, w)
		logrus.Errorf("generate token error")
		return
	}
	logrus.Println(tokenString)
	// utils.JSONResponse(0, "OK", tokenString, w)
	utils.JSONResponse(0, "OK", struct {
		PID      string `json:"pid"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}{"001", c.Username, tokenString}, w)

}
