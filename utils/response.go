package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//JSONResponse response web client json data
func JSONResponse(code int, msg string, data interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(resp{code, msg, data})
	if err != nil {
		fmt.Println("JSONResponse(): json.Marshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
