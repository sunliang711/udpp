package handlers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sunliang711/udpp/utils"
)

//根据pid和uid获取所有的权利
func getRights(w http.ResponseWriter, req *http.Request) {
	var (
	//err error
	)
	query := req.URL.Query()
	pid := query.Get("pid")
	uid := query.Get("uid")

	if len(pid) == 0 || len(uid) == 0 {
		errMsg := fmt.Sprint("Missing pid or uid")
		utils.JSONResponse(1, errMsg, nil, w)
		logrus.Error(errMsg)
		return
	}

	logrus.Debugf("pid: %v,uid: %v", pid, uid)

}

//更新用户权利
func updateRights(w http.ResponseWriter, req *http.Request) {

}
