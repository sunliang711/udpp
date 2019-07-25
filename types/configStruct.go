package types

const (
	SafeLaw = 1
	Gdpr    = 2
)

type Option struct {
	Label string `json:"label" bson:"label"`
	Value string `json:"value" bson:"value"`
}
type InfoItem struct {
	Name        string `json:"name" bson:"name"`
	Placeholder string `json:"placeholder" bson:"placeholder"`
	Value       string `json:"value" bson:"value"`
}
type SettingItem struct {
	Title    string     `json:"title" bson:"title"`
	InfoList []InfoItem `json:"infoList" bson:"infoList"`
}

type DetailItem struct {
	Type  string `json:"type" bson:"type"`
	Id    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`

	//type 'desc'
	Desc string `json:"desc" bson:"desc"`

	//type 'checkbox'
	Options []Option `json:"options" bson:"options"`
	Checked []string `json:"checked" bson:"checked"`

	//type 'api'
	Protocol    string `json:"protocol" bson:"protocol"`
	Placeholder string `json:"placeholder" bson:"placeholder"`
	Value       string `json:"value" bson:"value"`

	//type 'settings'
	SettingList []SettingItem `json:"settingList" bson:"settingList"`

	//type 'shareTable'
	CheckboxOptions []Option        `json:"checkboxOptions" bson:"checkboxOptions"`
	ShareList       []ShareListItem `json:"shareList" bson:"shareList"`
}

type ShareListItem struct {
	CompanyName string   `json:"companyName" bson:"companyName"`
	Checked     []string `json:"checked" bson:"checked"`
}

//--------------------------------------------------
type ConfigItem struct {
	RightID     int          `json:"right_id" bson:"right_id"`
	RightType   int          `json:"right_type" bson:"right_type"`
	RightName   string       `json:"right_name" bson:"right_name"`
	Description string       `json:"description" bson:"description"`
	Checked     int          `json:"checked" bson:"checked"`
	Details     []DetailItem `json:"details" bson:"details"`
}

type ConfigRes struct {
	PID        string       `json:"pid" bson:"pid"`
	Bgcolor    string       `json:"bgcolor" bson:"bgcolor"`
	Themecolor string       `json:"themecolor" bson:"themecolor"`
	Config     []ConfigItem `json:"config" bson:"config"`
}

func ConfigResTemplate() (r *ConfigRes) {
	r = &ConfigRes{
		PID:        r.PID,
		Bgcolor:    "",
		Themecolor: "",

		Config: []ConfigItem{
			{
				RightID:     KnownRightID,
				RightType:   SafeLaw,
				RightName:   KnownRightName,
				Description: KnownRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type:  "desc",
						Id:    "1",
						Title: "知情权详细描述",
						Desc:  KnownRightDescLong,
					},
					{
						Type:  "checkbox",
						Id:    "2",
						Title: "告知方式",
						Options: []Option{
							{"站内信", "0"},
							{"手机短信", "1"},
							{"邮件", "2"},
						},
						Checked: []string{},
					},
					{
						Type: "settings",
						Id:   "2",
						SettingList: []SettingItem{
							{
								Title: "企业基本信息配置",
								InfoList: []InfoItem{
									{Name: "名称", Placeholder: "请输入", Value: ""},
									{Name: "联系方式", Placeholder: "请输入", Value: ""},
									{Name: "地址", Placeholder: "请输入", Value: ""},
								},
							},
							{
								Title: "数据保护官信息配置",
								InfoList: []InfoItem{
									{Name: "名称", Placeholder: "请输入", Value: ""},
									{Name: "联系方式", Placeholder: "请输入", Value: ""},
									{Name: "邮箱", Placeholder: "请输入", Value: ""},
								},
							},
						},
					},
				},
			},
			{
				RightID:     AccessRightID,
				RightType:   SafeLaw,
				RightName:   AccessRightName,
				Description: AccessRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type: "desc",
						Id:   "",
						Desc: AccessRightDescLong,
					},
					{
						Type:        "api",
						Id:          "",
						Title:       "API地址",
						Protocol:    "",
						Placeholder: "请输入",
						Value:       "",
					},
				},
			},
			{
				RightID:     ForgetRightID,
				RightType:   SafeLaw,
				RightName:   ForgetRightName,
				Description: ForgetRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type: "desc",
						Id:   "",
						Desc: ForgetRightDescLong,
					},
					{
						Type:  "checkbox",
						Id:    "",
						Title: "遗忘方式",
						Options: []Option{
							{Label: "非业务关联数据定期删除", Value: "0"},
							{Label: "注销后删除所有个人信息", Value: "1"},
						},
						Checked: []string{},
					},
				},
			},
			{
				RightID:     PortableRightID,
				RightType:   SafeLaw,
				RightName:   PortableRightName,
				Description: PortableRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{

						Type: "desc",
						Id:   "",
						Desc: PortableRightDescLong,
					},
					{
						Type:        "api",
						Id:          "",
						Title:       "API地址",
						Protocol:    "",
						Placeholder: "请输入",
						Value:       "",
					},
				},
			},

			{
				RightID:     RefuseRightID,
				RightType:   Gdpr,
				RightName:   RefuseRightName,
				Description: RefuseRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type: "desc",
						Id:   "",
						Desc: RefuseRightDescLong,
					},
					{
						Type:  "checkbox",
						Id:    "",
						Title: "拒绝行为",
						Options: []Option{
							{Label: "获取我的用户画像", Value: "0"},
							{Label: "市场营销", Value: "1"},
						},
						Checked: []string{},
					},
				},
			},
			{
				RightID:     ModifyRightID,
				RightType:   Gdpr,
				RightName:   ModifyRightName,
				Description: ModifyRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type: "desc",
						Id:   "",
						Desc: ModifyRightDescLong,
					},
				},
			},
			{
				RightID:     ShareRightID,
				RightType:   Gdpr,
				RightName:   ShareRightName,
				Description: ShareRightDesc,
				Checked:     0,
				Details: []DetailItem{
					{
						Type: "shareTable",
						CheckboxOptions: []Option{
							{"用户名", "0"},
							{"手机", "1"},
							{"真实姓名", "2"},
							{"身份证号", "3"},
							{"身份证照片", "4"},
						},
						ShareList: []ShareListItem{
							{
								CompanyName: "ele",
								Checked:     []string{},
							},
							{
								CompanyName: "tencent",
								Checked:     []string{},
							},
						},
					},
				},
			},
		},
	}
	return
}

////////////////////////////////////////////////////

//type KnownRight struct {
//	Enabled int `json:"enabled"`
//
//	Desc struct {
//		ShortDesc string `json:"short_desc"`
//		LongDesc  string `json:"long_desc"`
//	} `json:"desc"`
//
//	options []
//	InApp   int `json:"in_app"`
//	Sms     int `json:"sms"`
//	Email   int `json:"email"`
//
//	CompanyName    string `json:"company_name"`
//	CompanyPhone   string `json:"company_phone"`
//	CompanyAddress string `json:"company_address"`
//	GuardName      string `json:"guard_name"`
//	GuardPhone     string `json:"guard_phone"`
//	GuardEmail     string `json:"guard_email"`
//}
//
//type AccessRight struct {
//	Enabled    int    `json:"enabled"`
//	ShortDesc  string `json:"short_desc"`
//	LongDesc   string `json:"long_desc"`
//	ApiAddress string `json:"api_address"`
//}
//
//type ForgetRight struct {
//	Enabled      int    `json:"enabled"`
//	ShortDesc    string `json:"short_desc"`
//	LongDesc     string `json:"long_desc"`
//	DeletePeriod int    `json:"delete_period"`
//	DeleteAll    int    `json:"delete_all"`
//}
//
//type PortableRight struct {
//	Enabled    int    `json:"enabled"`
//	ShortDesc  string `json:"short_desc"`
//	LongDesc   string `json:"long_desc"`
//	ApiAddress string `json:"api_address"`
//}
//
//type SafeLaw struct {
//	KnownRight    KnownRight    `json:"known_right"`
//	AccessRight   AccessRight   `json:"access_right"`
//	ForgetRight   ForgetRight   `json:"forget_right"`
//	PortableRight PortableRight `json:"portable_right"`
//}
//
////--------------------------------------------------
//
////GDPR
////--------------------------------------------------
//type RefuseRight struct {
//	Enabled   int    `json:"enabled"`
//	ShortDesc string `json:"short_desc"`
//	LongDesc  string `json:"long_desc"`
//}
//
//type ModifyRight struct {
//	Enabled   int    `json:"enabled"`
//	ShortDesc string `json:"short_desc"`
//	LongDesc  string `json:"long_desc"`
//}
//
//type ShareRight struct {
//	Enabled   int    `json:"enabled"`
//	ShortDesc string `json:"short_desc"`
//	LongDesc  string `json:"long_desc"`
//}
//
//type Gdpr struct {
//	RefuseRight RefuseRight `json:"refuse_right"`
//	ModifyRight ModifyRight `json:"modify_right"`
//	ShareRight  ShareRight  `json:"share_right"`
//}
//
////--------------------------------------------------
//
////ConfigResponse
//type ConfigResponse struct {
//	SafeLaw SafeLaw `json:"safe_law"`
//	Gdpr    Gdpr    `json:"gdpr"`
//}
