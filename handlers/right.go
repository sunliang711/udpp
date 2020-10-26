package handlers

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sirupsen/logrus"
	"github.com/sunliang711/udpp/models"
	"github.com/sunliang711/udpp/types"
	"github.com/sunliang711/udpp/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	blockdbDatabase   string
	blockdbCollection = "right"
)

type item struct {
	UdppItemHandler string `json:"udppItemHandler"`
	GmtCreated      string `json:"gmtCreated"`
}
type udppHandlers struct {
	Code string `json:"code"`
	Data struct {
		UdppHandlers []item `json:"udppHandlers"`
	} `json:"data"`
}

// accessList TODO
func accessList(w http.ResponseWriter, req *http.Request) {
	var err error
	query := req.URL.Query()
	pid := query.Get("pid")
	uid := query.Get("uid")
	pageNo := query.Get("pageNo")
	pageSize := "10"

	var config types.ConfigRes
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	mongoDBC := models.Mdb.Database(configDatabase).Collection(configCollection)
	result := mongoDBC.FindOne(ctx, bson.M{"pid": pid})
	err = result.Decode(&config)
	if err != nil {
		errMsg := fmt.Sprintf("query bgColor,themeColor from mongo error: %v", err)
		utils.JSONResponse(1, errMsg, nil, w)
		logrus.Error(errMsg)
		return
	}
	var timeline []types.TimeLineItem
	var uHandlers udppHandlers
	accessURL, err := config.GetURL(types.AccessRightID)
	logrus.Infof("access url: %v", accessURL)
	param := struct {
		PageNo   string `json:"pageNo"`
		PageSize string `json:"pageSize"`
		UserID   string `json:"userId"`
	}{pageNo, pageSize, uid}
	bs, _ := json.Marshal(&param)
	body := bytes.NewReader(bs)

	request, _ := http.NewRequest("POST", accessURL, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-anlink-userid", uid)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		msg := fmt.Sprintf("client.Do error: %v", err)
		logrus.Error(msg)
		utils.JSONResponse(200, msg, nil, w)
	}
	err = json.NewDecoder(resp.Body).Decode(&uHandlers)
	if err != nil {
		msg := fmt.Sprintf("Decode timeline error: %v", err)
		logrus.Error(msg)
		utils.JSONResponse(200, msg, nil, w)
	}
	logrus.Infof("uHandlers: %+v", uHandlers)
	if uHandlers.Code == "0000" {
		for _, v := range uHandlers.Data.UdppHandlers {
			logrus.Infof("item: %+v", v)
			timeline = append(timeline, types.TimeLineItem{v.UdppItemHandler, v.GmtCreated, true})
		}
	}
	utils.JSONResponse(200, "OK", timeline, w)
}

//根据pid和uid获取所有的权利
func getRights(w http.ResponseWriter, req *http.Request) {
	var (
		err error
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

	//1. 先从mongodb获取pid 对应的config信息
	var config types.ConfigRes
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	mongoDBC := models.Mdb.Database(configDatabase).Collection(configCollection)
	result := mongoDBC.FindOne(ctx, bson.M{"pid": pid})
	err = result.Decode(&config)
	if err != nil {
		errMsg := fmt.Sprintf("query bgColor,themeColor from mongo error: %v", err)
		utils.JSONResponse(1, errMsg, nil, w)
		logrus.Error(errMsg)
		return
	}
	logrus.Debugf("bgColor: %v themeColor: %v", config.Bgcolor, config.Themecolor)

	//q. 从blockDB中根据pid，uid读取当前状态，不存在时新建模版返回
	blockDBC := models.BlockDb.Database(blockdbDatabase).Collection(blockdbCollection)
	n, err := blockDBC.CountDocuments(ctx, bson.M{"pid": pid, "uid": uid})
	if err != nil {
		errMsg := fmt.Sprintf("query count right from blockDB with pid: %v uid: %v error: %v", pid, uid, err)
		utils.JSONResponse(1, errMsg, nil, w)
		logrus.Error(errMsg)
		return
	}
	var rightRes types.RightRes
	if n == 0 {
		logrus.Info("No data,use right template")
		rightRes = *types.RightResTemplate(pid, uid)
		rightRes.ID = primitive.NewObjectID()
		logrus.Infof("new objectID: %v", rightRes.ID)
	} else if n == 1 {
		logrus.Infof("FindOne with pid: %v,uid: %v in blockDB", pid, uid)
		result := blockDBC.FindOne(ctx, bson.M{"pid": pid, "uid": uid})
		err = result.Decode(&rightRes)
		logrus.Debugf("FindOne result: %+v", rightRes)
		if err != nil {
			errMsg := fmt.Sprintf("Decode RightRes error: %v", err)
			logrus.Error(errMsg)
			utils.JSONResponse(1, errMsg, nil, w)
			return
		}
	} else {
		errMsg := fmt.Sprintf("More than one Right with pid: %v uid: %v", pid, uid)
		utils.JSONResponse(1, errMsg, nil, w)
		logrus.Error(errMsg)
		return
	}

	//TODO always from config
	// rightRes.BgColor
	// rightRes.ThemeColor
	rightRes.BgColor = config.Bgcolor
	rightRes.ThemeColor = config.Themecolor
	rightRes.Link = fmt.Sprintf("%s%s", viper.GetString("blockAddress"), rightRes.ID.Hex())

	//知情权 config -> Details -> type"settings"       == rightRes -> permission -> Details -> tableInfo -> tableData
	var tableInfoItem [][]types.TableInfoItem
	for _, c := range config.Config {
		if c.RightID == types.KnownRightID {
			logrus.Debugf("known right")
			for _, d := range c.Details {
				logrus.Debugf("c.Details: %+v", c.Details)
				if d.Type == "settings" {
					logrus.Debugf("config -> details -> settings")
					for _, s := range d.SettingList {
						var ti []types.TableInfoItem
						logrus.Debugf("s: %v", s)
						var infoItem types.TableInfoItem
						for _, info := range s.InfoList {
							infoItem.Title = s.Title
							infoItem.Label = info.Name
							infoItem.Txt = info.Value
							logrus.Debugf("append %v to ti", infoItem)
							ti = append(ti, infoItem)
						}
						tableInfoItem = append(tableInfoItem, ti)
					}
					break
				}
			}
			break
		}
	}
	logrus.Debugf("tableInfoItem: %v", tableInfoItem)
	for i, r := range rightRes.PermissionList {
		if r.ID == types.KnownRightID {
			for j, d := range r.Details {
				if d.Type == "tableInfo" {
					rightRes.PermissionList[i].Details[j].TableData = tableInfoItem
					break
				}
			}
			break
		}
	}

	//访问权
	var timeline []types.TimeLineItem
	accessURL, err := config.GetURL(types.AccessRightID)
	logrus.Infof("access url: %v", accessURL)
	if err != nil {
		logrus.Warnf("Cannot find url of access right")
	} else {
		accessRes, err := http.Get(accessURL + fmt.Sprintf("?uid=", uid))
		if err != nil {
			logrus.Warnf("Get access timeline error: %v", err)
		} else {
			err = json.NewDecoder(accessRes.Body).Decode(&timeline)
			if err != nil {
				logrus.Warnf("Decode access timeline error: %v", err)
			}
		}
	}
	//for test
	timeline = []types.TimeLineItem{
		{"用户画像", "2019/6/25 12:30", true},
		{"市场营销", "2019/6/26 12:40", true},
		{"大数据分析", "2019/6/27 12:48", true},
		{"短信通知", "2019/6/28 12:49", true},
	}
	for i, p := range rightRes.PermissionList {
		if p.ID == types.AccessRightID {
			for j, d := range p.Details {
				if d.Type == "timeline" {
					rightRes.PermissionList[i].Details[j].Data = timeline
				}
			}
			break
		}
	}

	//可携带权
	downURL, err := config.GetURL(types.PortableRightID)
	if err != nil {
		logrus.Warnf("Cannot find url of portable right")
	}
	logrus.Infof("download url: %v", downURL)

	for i, p := range rightRes.PermissionList {
		if p.ID == types.PortableRightID {
			for j, d := range p.Details {
				if d.Type == "download" {
					rightRes.PermissionList[i].Details[j].URL = downURL
					break
				}
			}
			// rightRes.PermissionList[i].Details = append(rightRes.PermissionList[i].Details, types.Detail{
			// 	Type: "download",
			// 	ID:   "22",
			// 	URL:  downURL,
			// })
			break
		}
	}

	//共享权
	for i, p := range rightRes.PermissionList {
		if p.ID == types.ShareRightID {
			for _, c := range config.Config {
				if c.RightID == types.ShareRightID {
					//TODO replace all "shareTable" with c.Details
					checked := rightRes.PermissionList[i].Details[2].Checked
					rightRes.PermissionList[i].Details = []types.Detail{rightRes.PermissionList[i].Details[0], rightRes.PermissionList[i].Details[1]}
					for _, d := range c.Details {
						if d.Type == "shareTable" {
							rightRes.PermissionList[i].Details = append(rightRes.PermissionList[i].Details, types.Detail{
								Type:            "shareTable",
								CheckboxOptions: d.CheckboxOptions,
								ShareList:       d.ShareList,
								Checked:         checked,
							})

						}
					}
					break
				}
			}
			break
		}
	}

	//合并checkbox
	for _, configItem := range config.Config {
		rightID := configItem.RightID
		for i, permission := range rightRes.PermissionList {
			if permission.ID == rightID {
				logrus.Debugf("rightID: %v", rightID)
				rightRes.PermissionList[i].Disabled = (configItem.Checked == 0)

				id2DisableList := make(map[string][]string)
				for _, detailItem := range configItem.Details {
					var disabledList []string
					if detailItem.Options != nil {
						for _, o := range detailItem.Options {
							logrus.Debugf("detailItem.Checked: %v", detailItem.Checked)
							if !utils.IsIn(o.Value, detailItem.Checked) {
								logrus.Debugf("append o.Value: %v to disabledList", o.Value)
								disabledList = append(disabledList, o.Value)
							}
						}
						id2DisableList[detailItem.Id] = disabledList
						logrus.Debugf("id2DisableList[%v]: %v", detailItem.Id, id2DisableList[detailItem.Id])
					}
				}

				for j, d := range rightRes.PermissionList[i].Details {
					if d.Options != nil {
						logrus.Debugf("set disabled: %v", id2DisableList[d.ID])
						logrus.Debugf("d.ID: %v", d.ID)
						rightRes.PermissionList[i].Details[j].Disabled = id2DisableList[d.ID]
					}
				}

				break
			}
		}

	}

	logrus.Debugf("rightRes: %+v", rightRes)
	utils.JSONResponse(0, "OK", rightRes, w)

}

//更新用户权利
func updateRights(w http.ResponseWriter, req *http.Request) {
	var (
		rightRes types.RightRes
		err      error
	)

	err = json.NewDecoder(req.Body).Decode(&rightRes)
	if err != nil {
		errMsg := fmt.Sprint("bad request data format")
		logrus.Error(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	logrus.Debugf("updateRights: request body: %+v", rightRes)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	C := models.BlockDb.Database(blockdbDatabase).Collection(blockdbCollection)
	n, err := C.CountDocuments(ctx, bson.M{"pid": rightRes.PID, "uid": rightRes.UID})
	if err != nil {
		errMsg := fmt.Sprintf("count document error with pid: %v,uid: %v,%v", rightRes.PID, rightRes.UID, err)
		logrus.Error(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}

	if n == 0 {
		//insert
		logrus.Infof("insert...")
		_, err := C.InsertOne(ctx, &rightRes)
		if err != nil {
			errMsg := fmt.Sprintf("insert failed: %v", err)
			utils.JSONResponse(1, errMsg, nil, w)
			logrus.Error(errMsg, err)
			return
		}
		logrus.Info("OK")
	} else if n == 1 {
		logrus.Infof("update...")
		_, err = C.UpdateOne(ctx, bson.M{"pid": rightRes.PID, "uid": rightRes.UID}, bson.M{"$set": &rightRes})
		if err != nil {
			errMsg := fmt.Sprintf("update failed: %v", err)
			utils.JSONResponse(1, errMsg, nil, w)
			logrus.Error(errMsg, err)
			return
		}
	} else {
		errMsg := fmt.Sprintf("more than 1 record in db")
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}

	userID, _ := strconv.Atoi(rightRes.UID)
	bs, err := json.Marshal(rightRes)
	//var b bytes.Buffer
	//binary.Write(&b,binary.BigEndian,&rightRes)
	if err != nil {
		logrus.Errorf("Encode rightRes to json error: %v", err)
	} else {
		h := md5.New()
		h.Write(bs)
		md5sum := h.Sum(nil)

		rightResBs, _ := json.Marshal(rightRes)
		//toShop
		toShop := struct {
			UserID             int    `json:"userId"`
			UdppPrivacyContent string `json:"udppPrivacyContent"`
			Hash               string `json:"hash"`
			BlockAddress       string `json:"blockaddress"`
		}{
			userID,
			string(rightResBs),
			fmt.Sprintf("%x", md5sum),
			fmt.Sprintf("%s%s", viper.GetString("blockAddress"), rightRes.ID.Hex()),
		}
		toShopBs, _ := json.Marshal(toShop)
		logrus.Debugf("toShopBs: %s", string(toShopBs))
		var b io.Reader
		b = bytes.NewBuffer(toShopBs)
		link := "http://17391-zis-stargate-tech-anlink-web-gateway.test.za-tech.net/tech-anlink-user/app/v2/anlink/sign/insertUdppPrivacy"
		res, err := http.Post(link, "application/json", b)
		if err != nil {
			logrus.Errorf("Post error: %v", err)
		} else {
			io.Copy(os.Stdout, res.Body)
		}
	}

	utils.JSONResponse(0, "OK", nil, w)
}
