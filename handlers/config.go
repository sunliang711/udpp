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
	database         = "udpp"
	configCollection = "config"
)

type getConfigReq struct {
	PID string `json:"pid"`
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	var (
		r   getConfigReq
		err error
	)
	err = json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		errMsg := fmt.Sprintf("getConfig: invalid request body: %v", err)
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	c := models.Mdb.Database(database).Collection(configCollection)
	n, err := c.CountDocuments(ctx, bson.M{"pid": r.PID})
	if err != nil {
		errMsg := "query count error"
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	var res types.ConfigRes
	if n == 0 {
		logrus.Infof("No current data,set default")
		//set default value
		res = *types.ConfigResTemplate()

	} else {
		singleResult := c.FindOne(ctx, bson.M{"pid": r.PID})
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
		errMsg := "bad request data format"
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	logrus.Debugf("request data: %+v\n", c)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	C := models.Mdb.Database(database).Collection(configCollection)
	n, err := C.CountDocuments(ctx, bson.M{"pid": c.PID})
	if err != nil {
		errMsg := fmt.Sprintf("count document error with pid: %v,%v", c.PID, err)
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}
	if n == 0 {
		//insert
		logrus.Infof("insert...")
		_, err = C.InsertOne(ctx, &c)
		if err != nil {
			utils.JSONResponse(1, "insert failed", nil, w)
			logrus.Infof("insert failed: %v", err)
			return
		}
		logrus.Info("OK")
	} else if n == 1 {
		logrus.Infof("update...")
		_, err = C.UpdateOne(ctx, bson.M{"pid": c.PID}, bson.M{"$set": &c})
		if err != nil {
			utils.JSONResponse(1, "update failed", nil, w)
			logrus.Infof("update failed: %v", err)
			return
		}
		logrus.Info("OK")
	} else {
		errMsg := fmt.Sprintf("more than 1 record in db")
		logrus.Errorf(errMsg)
		utils.JSONResponse(1, errMsg, nil, w)
		return
	}

	utils.JSONResponse(0, "OK", nil, w)

}

//获取全量的right list，并根据pid，把打勾的标出来,格式为 `商户配置页` 定义的格式
//func getConfig(w http.ResponseWriter, req *http.Request) {
//	var (
//		r   getConfigReq
//		err error
//		sql string
//		id  int
//	)
//	err = json.NewDecoder(req.Body).Decode(&r)
//	if err != nil {
//		errMsg := "getConfig: invalid request body"
//		logrus.Errorf(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//
//	pid := r.PID
//	logrus.Debugf("getConfig: get pid: %d", pid)
//
//	//知情权
//	//types.KnownRight
//
//	//1. get shortDesc,longDesc
//	//    select shortDesc,longDesc from `right` where rightID = types.KnownRight;
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err := models.Db.Query(sql, types.KnownRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	var (
//		shortDesc string
//		longDesc  string
//	)
//	if rows.Next() {
//		rows.Scan(&shortDesc, &longDesc)
//		logrus.Debugf("shortDesc: %v;longDesc: %v", shortDesc, longDesc)
//	} else {
//		errMsg := fmt.Sprintf("no shortDesc,longDesc data in right table")
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//
//	//2. 根据pid查找对应的知情权配置项，没有则给出默认值
//	sql = "select enabled,inApp,sms,email,companyName,companyPhone,guardName,guardPhone,guardEmail from `knownRight` where pid=?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query data from knownRight error,with pid: %v", pid)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	var (
//		enabled int
//		inApp   int
//		sms     int
//		email   int
//
//		companyName    string
//		companyPhone   string
//		companyAddress string
//
//		guardName  string
//		guardPhone string
//		guardEmail string
//	)
//
//	if rows.Next() {
//		rows.Scan(&enabled, &inApp, &sms, &email, &companyName, &companyPhone, &companyAddress, &guardName, &guardPhone, &guardEmail)
//		logrus.Debugf("enabled: %v,inApp: %v,sms: %v,email: %v", enabled, inApp, sms, email)
//		logrus.Debugf("companyName: %v,companyPhone: %v,companyAddress: %v", companyName, companyPhone, companyAddress)
//		logrus.Debugf("guardName: %v,guardPhone: %v,guardEmail: %v", guardName, guardPhone, guardEmail)
//	} else {
//		logrus.Infof("No data in knownRight,use default value(zero value).")
//	}
//
//	//3. 1和2构成对象ConfigItem(是SafeLaw的一部分）
//	id += 1
//	descDetail := types.DetailItem{
//		Type:  "desc",
//		Id:    fmt.Sprintf("%d", id),
//		Title: "知情权详细描述",
//		Desc:  longDesc,
//	}
//	id += 1
//	checkboxDetail := types.DetailItem{
//		Type:  "checkbox",
//		Id:    fmt.Sprintf("%d", id),
//		Title: "告知方式",
//		Options: []types.Option{
//			{Label: "站内信", Value: "0"},
//			{Label: "手机短信", Value: "1"},
//			{Label: "邮件", Value: "2"},
//		},
//		Checked: utils.OptionCheckd(inApp, sms, email),
//	}
//	id += 1
//	settingsDetail := types.DetailItem{
//		Type: "settings",
//		Id:   fmt.Sprintf("%d", id),
//		SettingList: []types.SettingItem{
//			{
//				Title: "企业基本信息配置",
//				InfoList: []types.InfoItem{
//					{Name: "名称", Placeholder: "请输入", Value: companyName},
//					{Name: "联系方式", Placeholder: "请输入", Value: companyPhone},
//					{Name: "地址", Placeholder: "请输入", Value: companyAddress},
//				},
//			},
//			{
//				Title: "数据保护官信息配置",
//				InfoList: []types.InfoItem{
//					{Name: "名称", Placeholder: "请输入", Value: guardName},
//					{Name: "联系方式", Placeholder: "请输入", Value: guardPhone},
//					{Name: "邮箱", Placeholder: "请输入", Value: guardEmail},
//				},
//			},
//		},
//	}
//	knownItem := types.ConfigItem{
//		RightID:     types.KnownRightID,
//		RightType:   types.SafeLaw,
//		RightName:   "知情权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//			checkboxDetail,
//			settingsDetail,
//		},
//	}
//	logrus.Debugf("rightItem: %+v", knownItem)
//
//	//访问权
//	//types.AccessRight
//	//1. get shortDesc,longDesc
//	//    select shortDesc,longDesc from `right` where rightID = types.accessRight;
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err = models.Db.Query(sql, types.AccessRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	if rows.Next() {
//		rows.Scan(&shortDesc, &longDesc)
//		logrus.Debugf("shortDesc: %v;longDesc: %v", shortDesc, longDesc)
//	} else {
//		errMsg := fmt.Sprintf("no shortDesc,longDesc data in right table")
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	//2. 从accessRight表查询
//	sql = "select enabled,apiAddress from `accessRight` where pid = ?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query enabled,apiAddress from accessRight error: %v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//
//	var apiAddress string
//	if rows.Next() {
//		rows.Scan(&enabled, &apiAddress)
//	}
//	id += 1
//	descDetail = types.DetailItem{
//		Type: "desc",
//		Id:   fmt.Sprintf("%d", id),
//		Desc: longDesc,
//	}
//	id += 1
//	apiDetail := types.DetailItem{
//		Type:        "api",
//		Id:          fmt.Sprintf("%d", id),
//		Title:       "API地址",
//		Protocol:    "",
//		Placeholder: "请输入",
//		Value:       apiAddress,
//	}
//	accessItem := types.ConfigItem{
//		RightID:     types.AccessRightID,
//		RightType:   types.SafeLaw,
//		RightName:   "访问权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//			apiDetail,
//		},
//	}
//	logrus.Debugf("accessItem: %+v", accessItem)
//
//	//遗忘权
//	//types.ForgetRight
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err = models.Db.Query(sql, types.ForgetRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//
//	sql = "select enabled,deletePeriod,deleteAll from `forgetRight` where pid = ?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query enabled,deletePeriod,deleteAll from forgetRight error: %v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	var (
//		deletePeriod int
//		deleteAll    int
//	)
//	if rows.Next() {
//		rows.Scan(&enabled, &deletePeriod, &deleteAll)
//	}
//
//	id += 1
//	descDetail = types.DetailItem{
//		Type: "desc",
//		Id:   fmt.Sprintf("%d", id),
//		Desc: longDesc,
//	}
//	id += 1
//	checkboxDetail = types.DetailItem{
//		Type:  "checkbox",
//		Id:    fmt.Sprintf("%d", id),
//		Title: "遗忘方式",
//		Options: []types.Option{
//			{Label: "非业务关联数据定期删除", Value: "0"},
//			{Label: "注销后删除所有个人信息", Value: "1"},
//		},
//		Checked: utils.OptionCheckd(deletePeriod, deleteAll),
//	}
//	forgetItem := types.ConfigItem{
//		RightID:     types.ForgetRightID,
//		RightType:   types.SafeLaw,
//		RightName:   "遗忘权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//			checkboxDetail,
//		},
//	}
//	logrus.Debugf("forget item: %+v", forgetItem)
//
//	//可携带权
//	//types.PortableRight
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err = models.Db.Query(sql, types.PortableRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	sql = "select enabled,apiAddress from `portableRight` where pid = ?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query enabled,apiAddress from portableRight error: %v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	if rows.Next() {
//		rows.Scan(&enabled, &apiAddress)
//	}
//
//	id += 1
//	descDetail = types.DetailItem{
//		Type: "desc",
//		Id:   fmt.Sprintf("%d", id),
//		Desc: longDesc,
//	}
//
//	id += 1
//	apiDetail = types.DetailItem{
//		Type:        "api",
//		Id:          fmt.Sprintf("%d", id),
//		Title:       "API地址",
//		Protocol:    "",
//		Placeholder: "请输入",
//		Value:       apiAddress,
//	}
//
//	portableItem := types.ConfigItem{
//		RightID:     types.PortableRightID,
//		RightType:   types.SafeLaw,
//		RightName:   "可携带权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//			apiDetail,
//		},
//	}
//	logrus.Debugf("portable item: %+v", portableItem)
//
//	//GDPR
//	//拒绝权
//	//types.RefuseRight
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err = models.Db.Query(sql, types.RefuseRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//
//	sql = "select enabled,picture,market from `refuseRight` where pid = ?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query enabled,picture,market from refuseRight error: %v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	var (
//		picture int
//		market  int
//	)
//	if rows.Next() {
//		rows.Scan(&picture, &make())
//	}
//	id += 1
//	descDetail = types.DetailItem{
//		Type: "desc",
//		Id:   fmt.Sprintf("%d", id),
//		Desc: longDesc,
//	}
//
//	id += 1
//	checkboxDetail = types.DetailItem{
//		Type:  "checkbox",
//		Id:    fmt.Sprintf("%d", id),
//		Title: "拒绝行为",
//		Options: []types.Option{
//			{Label: "获取我的用户画像", Value: "0"},
//			{Label: "市场营销", Value: "1"},
//		},
//		Checked: utils.OptionCheckd(picture, market),
//	}
//	refuseItem := types.ConfigItem{
//		RightID:     types.RefuseRightID,
//		RightType:   types.Gdrp,
//		RightName:   "拒绝权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//			checkboxDetail,
//		},
//	}
//	logrus.Debugf("refuse item: %+v", refuseItem)
//
//	//修改权
//	//types.ModifyRight
//	sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	rows, err = models.Db.Query(sql, types.ModifyRightID)
//	if err != nil {
//		errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	sql = "select enabled from `modifyRight` where pid = ?"
//	rows, err = models.Db.Query(sql, pid)
//	if err != nil {
//		errMsg := fmt.Sprintf("query eanbled from modifyRight error: %v", err)
//		logrus.Error(errMsg)
//		utils.JSONResponse(1, errMsg, nil, w)
//		return
//	}
//	id += 1
//	descDetail = types.DetailItem{
//		Type: "desc",
//		Id:   fmt.Sprintf("%d", id),
//		Desc: longDesc,
//	}
//
//	modifyItem := types.ConfigItem{
//		RightID:     types.ModifyRightID,
//		RightType:   types.Gdrp,
//		RightName:   "修正权",
//		Description: shortDesc,
//		Checked:     enabled,
//		Details: []types.DetailItem{
//			descDetail,
//		},
//	}
//	logrus.Debugf("modify item: %+v", modifyItem)
//	//TODO
//	////共享权
//	////types.ShareRight
//	//sql = "select shortDesc,longDesc from `right` where rightID = ?"
//	//rows, err = models.Db.Query(sql, types.ShareRightID)
//	//if err != nil {
//	//	errMsg := fmt.Sprintf("query shortDesc,longDesc error:%v", err)
//	//	logrus.Error(errMsg)
//	//	utils.JSONResponse(1, errMsg, nil, w)
//	//	return
//	//}
//	// sql = "select enabled, TODO from `shareRight` where pid = ?"
//	// rows,err = models.Db.Query(sql,types.ShareRightID)
//	// if err != nil{
//	// 	errMsg := fmt.Sprintf("query enabled,TODO from shareRight error: %v",err)
//	// 	logrus.Error(errMsg)
//	// 	utils.JSONResponse(1,errMsg,nil,w)
//	//	 return
//	// }
//	// var(
//	// 	TODO
//	// )
//	// if rows.Next(){
//	// 	rows.Scan(&enabled,& TODO)
//	// }
//
//	res := types.ConfigRes{
//		//TODO add shareItem
//		Config: []types.ConfigItem{knownItem, accessItem, forgetItem, portableItem, refuseItem, modifyItem},
//	}
//
//	utils.JSONResponse(0, "OK", res, w)
//}
