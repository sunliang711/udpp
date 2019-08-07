package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TableInfoItem struct {
	Title string `json:"title" bson:"title"`
	Label string `json:"label" bson:"label"`
	Txt   string `json:"txt" bson:"txt"`
}

type TimeLineItem struct {
	Txt   string `json:"txt" bson:"txt"`
	Time  string `json:"time" bson:"time"`
	Seled bool   `json:"seled" bson:"seled"`
}
type Detail struct {
	Type  string `json:"type" bson:"type"`
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`

	// type "checkbox"
	Options  []Option `json:"options" bson:"options"`
	Checked  []string `json:"checked" bson:"checked"`
	Disabled []string `json:"disabled" bson:"disabled"`

	// type "collapseTitle"
	Desc string `json:"desc" bson:"desc"`

	//type "tableInfo"
	TableData [][]TableInfoItem `json:"tableData" bson:"tableData"`

	//type "simpleItems"
	Items []string `json:"items" bson:"items"`

	//type "timeline"
	Data []TimeLineItem `json:"data" bson:"data"`

	//type "download"
	URL string `json:"url" bson:"url"`

	//type "shareTable"
	CheckboxOptions []Option        `json:"checkboxOptions" bson:"checkboxOptions"`
	ShareList       []ShareListItem `json:"shareList" bson:"shareList"`
}

type Permission struct {
	//rightID
	ID         int      `json:"id" bson:"id"`
	Checked    bool     `json:"checked" bson:"checked"`
	Title      string   `json:"title" bson:"title"`
	Desc       string   `json:"desc" bson:"desc"`
	HasDetails bool     `json:"hasDetails" bson:"hasDetails"`
	Details    []Detail `json:"details" bson:"details"`
	Disabled   bool     `json:"disabled" bson:"disabled"`
}

type RightRes struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	PID            string             `json:"pid" bson:"pid"`
	UID            string             `json:"uid" bson:"uid"`
	BgColor        string             `json:"bgcolor" bson:"bgcolor"`
	ThemeColor     string             `json:"themecolor" bson:"themecolor"`
	PermissionList []Permission       `json:"permissionList" bson:"permissionList"`
	Link           string             `json:"link" bson:"link"`
}

func RightResTemplate(pid, uid string) (r *RightRes) {
	r = &RightRes{
		PID:        pid,
		UID:        uid,
		BgColor:    "",
		ThemeColor: "",
		PermissionList: []Permission{
			{
				ID:         KnownRightID,
				Checked:    true,
				Title:      KnownRightName,
				Desc:       KnownRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "checkbox",
						ID:    "100", //要和configStruct中的一一对应，并且相同数组里的ID要唯一
						Title: "告知方式",
						Options: []Option{
							{
								"站内信",
								"0",
							},
							{
								"手机短信",
								"1",
							},
							{
								"邮件",
								"2",
							},
						},
						Checked:  []string{},
						Disabled: []string{},
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "知情权详情",
						Desc:  KnownRightDescLong,
					},
					{
						Type:      "tableInfo",
						ID:        "14",
						TableData: [][]TableInfoItem{},
					},
					{
						Type:  "simpleItems",
						ID:    "15",
						Title: "数据处理涉及使用目的",
						Items: []string{"市场营销", "用户画像", "数据挖掘"},
					},
				},
			},
			{
				ID:         AccessRightID,
				Checked:    true,
				Title:      AccessRightName,
				Desc:       AccessRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "访问权详细",
						Desc:  AccessRightDescLong,
					},
					{
						Type: "timeline",
						ID:   "16",
						Data: []TimeLineItem{},
					},
				},
			},
			{
				ID:         ForgetRightID,
				Checked:    true,
				Title:      ForgetRightName,
				Desc:       ForgetRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "110",
					},
					{
						Type:  "checkbox",
						ID:    "110",
						Title: "遗忘方式",
						Options: []Option{
							{"非业务关联数据定期删除", "0"},
							{"注销后删除所有个人信息", "1"},
						},
						Checked:  []string{},
						Disabled: []string{},
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "遗忘权详细",
						Desc:  ForgetRightDescLong,
					},
				},
			},
			{
				ID:         PortableRightID,
				Checked:    true,
				Title:      PortableRightName,
				Desc:       PortableRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "可携带权",
						Desc:  PortableRightDescLong,
					},
					{
						Type: "download",
						ID:   "22",
						URL:  "",
					},
				},
			},
			{
				ID:         RefuseRightID,
				Checked:    true,
				Title:      RefuseRightName,
				Desc:       RefuseRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "checkbox",
						ID:    "120",
						Title: "拒绝行为",
						Options: []Option{
							{"获取我的用户画像", "0"},
							{"市场营销", "1"},
						},
						Checked:  []string{},
						Disabled: []string{},
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "拒绝权详细",
						Desc:  RefuseRightDescLong,
					},
				},
			},
			{
				ID:         ModifyRightID,
				Checked:    true,
				Title:      ModifyRightName,
				Desc:       ModifyRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "collapseTitle",
						ID:    "13",
						Title: "修正权详细",
						Desc:  ModifyRightDescLong,
					},
				},
			},
			{
				ID:         ShareRightID,
				Checked:    true,
				Title:      ShareRightName,
				Desc:       ShareRightDesc,
				HasDetails: true,
				Details: []Detail{
					{
						Type: "title",
						ID:   "10",
					},
					{
						Type:  "collapseTitle",
						ID:    "14",
						Title: "共享权详细",
						Desc:  ShareRightDescLong,
					},
					{
						Type:            "shareTable",
						ID:              "50",
						CheckboxOptions: []Option{},
						ShareList:       []ShareListItem{},
						Checked:         []string{},
					},
				},
			},
		},
	}
	return
}
