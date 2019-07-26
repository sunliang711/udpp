package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sunliang711/udpp/models"
	"github.com/sunliang711/udpp/types"
	"github.com/sunliang711/udpp/utils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	configDatabase   = "udpp"
	configCollection = "config"
)

type getConfigReq struct {
	PID string `json:"pid"`
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	pid := query.Get("pid")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := models.Mdb.Database(configDatabase).Collection(configCollection)
	n, err := c.CountDocuments(ctx, bson.M{"pid": pid})
	if err != nil {
		errMsg := fmt.Sprintf("Query count error with pid: %v", pid)
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	var res types.ConfigRes
	if n == 0 {
		logrus.Infof("No current config data with pid: %v,set default", pid)
		//set default value
		res = *types.ConfigResTemplate(pid)

	} else {
		singleResult := c.FindOne(ctx, bson.M{"pid": pid})
		singleResult.Decode(&res)
	}

	utils.JSONResponse(0, "OK", res, w)

}

//更新对应pid的勾选状态，以及bgcolor，themecolor
func updateConfig(w http.ResponseWriter, req *http.Request) {
	var (
		c   types.ConfigRes
		err error
	)

	err = json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		errMsg := "Bad request data format"
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	logrus.Debugf("Request data: %+v\n", c)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	C := models.Mdb.Database(configDatabase).Collection(configCollection)
	n, err := C.CountDocuments(ctx, bson.M{"pid": c.PID})
	if err != nil {
		errMsg := fmt.Sprintf("count document error with pid: %v,%v", c.PID, err)
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	if n == 0 {
		//insert
		logrus.Infof("Insert...")
		_, err = C.InsertOne(ctx, &c)
		if err != nil {
			utils.JSONResponse(1, "insert failed", nil, w)
			logrus.Infof("insert failed: %v", err)
			return
		}
		logrus.Info("OK")
	} else if n == 1 {
		logrus.Infof("Update...")
		_, err = C.UpdateOne(ctx, bson.M{"pid": c.PID}, bson.M{"$set": &c})
		if err != nil {
			errMsg := fmt.Sprintf("Update failed: %v", err)
			utils.JSONResponse(1, errMsg, nil, w)
			logrus.Infof(errMsg, err)
			return
		}
		logrus.Info("OK")
	} else {
		errMsg := fmt.Sprintf("More than 1 record in db")
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}

	utils.JSONResponse(0, "OK", nil, w)
}
