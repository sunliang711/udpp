package types

const (
	KnownRightID       = 1
	KnownRightName     = "知情权"
	KnownRightDesc     = "收集您信息时，必须给与您充分的知情权"
	KnownRightDescLong = "数据控制者收集个人信息必须给予个人充分的知情权。应当向数据主体提供相应的信息以保障数据主体对控制者身份、个人信息处理目的及方式、权利维护途径等内容的知晓。比如，数据控制者在收集与数据主体相关的个人数据时，应当告知数据主体，包括数据控制者的身份与详细联系方式、数据保护官的详细联系方式、数据处理将涉及的个人数据的使用目的，以及处理个人数据的法律依据等。"
)

const (
	AccessRightID       = 2
	AccessRightName     = "访问权"
	AccessRightDesc     = "您有权访问个人数据并有权获知相关信息"
	AccessRightDescLong = "数据主体有权从数据控制者那里得知关于其个人数据是否正在被处理的真实情形，如果其数据正在被处理的话，数据主体应当有权访问个人数据并有权获知相关信息，如该数据处理的目的；相关个人数据的类型；个人数据已经被或将被披露给数据接收者或接收者的类型，特别是当数据的接收者属于第三国或国际组织；在可能的情形下，个人数据将被储存的预期期限等。"
)

const (
	ForgetRightID       = 3
	ForgetRightName     = "遗忘权"
	ForgetRightDesc     = "有权要求商户及时删除您的个人信息数据"
	ForgetRightDescLong = "数据主体有权要求数据控制者及时删除其个人相关数据的权利。当出现“个人数据对于实现其被收集或处理的相关目的不再必要”等六种情形之一时，数据控制者有责任及时删除其个人数据。"
)
const (
	PortableRightID       = 4
	PortableRightName     = "可携带权"
	PortableRightDesc     = "有权要求商户提供与您个人有关的数据下载"
	PortableRightDescLong = "数据主体有权以结构化、通用和机器可读的格式接收其提供给数据控制者的与其有关个人数据，数据主体有权将这些数据传输给另一个数据控制者，而被要求转移数据的控制者应当配合数据主体的相应要求。用户数据可携权不仅赋予用户取得、重复利用相关数据的权利，还赋予用户传输该等数据的权利。"
)
const (
	RefuseRightID       = 5
	RefuseRightName     = "拒绝权"
	RefuseRightDesc     = "有权拒绝商户对您数据的非公益使用方式"
	RefuseRightDescLong = "有关数据主体的数据处理，包括根据这些条款而进行的用户画像，数据主体有权随时提出反对。此时，数据控制者须立即停止针对这部分个人数据的处理行为，除非数据控制者证明，相比数据主体的利益、权利和自由，具有压倒性的正当理由需要进行处理，或者处理是为了提起、行使或辩护法律性主张。"
)
const (
	ModifyRightID       = 6
	ModifyRightName     = "修正权"
	ModifyRightDesc     = "您有权要求商户及时纠正与您相关的个人信息"
	ModifyRightDescLong = "数据主体有权要求数据控制及时地纠正与其相关的不准确个人数据。考虑到处理的目的，数据主体应当有权使不完整的个人数据完整，包括通过提供补充声明的方式进行完善。"
)
const (
	ShareRightID       = 7
	ShareRightName     = "共享权"
	ShareRightDesc     = "经用户的明示授权后，企业可以将用户相关的数据分享给合作伙伴"
	ShareRightDescLong = "用户具有对自身数据的完全控制权。用户数据存储在企业端，当用户明确授权自身数据的共享时，企业可以合法与明示合作伙伴共享用户数据。用户的数据以密文形式存储在云端，共享过程使用PRE实现。公私钥对和数据密钥又DKMS生成和管理。"
)
