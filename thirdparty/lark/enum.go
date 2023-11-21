package lark

type LarkReceiverIdType string

const (
	LarkReceiverIdTypeOpenID  LarkReceiverIdType = "open_id"
	LarkReceiverIdTypeUserID  LarkReceiverIdType = "user_id"
	LarkReceiverIdTypeUnionID LarkReceiverIdType = "union_id"
	LarkReceiverIdTypeEmail   LarkReceiverIdType = "email"
	LarkReceiverIdTypeChatID  LarkReceiverIdType = "chat_id"
)

var (
	supportedLarkReceiverIdType = map[string]LarkReceiverIdType{
		LarkReceiverIdTypeOpenID.String():  LarkReceiverIdTypeOpenID,
		LarkReceiverIdTypeUserID.String():  LarkReceiverIdTypeUserID,
		LarkReceiverIdTypeUnionID.String(): LarkReceiverIdTypeUnionID,
		LarkReceiverIdTypeEmail.String():   LarkReceiverIdTypeEmail,
		LarkReceiverIdTypeChatID.String():  LarkReceiverIdTypeChatID,
	}
)

func (t LarkReceiverIdType) String() string { return string(t) }

func getLarkReceiverIdType(idType LarkReceiverIdType) string {
	_, exist := supportedLarkReceiverIdType[idType.String()]
	if !exist {
		return LarkReceiverIdTypeOpenID.String()
	} else {
		return idType.String()
	}
}

type LarkMarkdownMessageTheme string

const (
	LarkMarkdownMessageThemeBlue      LarkMarkdownMessageTheme = "blue"
	LarkMarkdownMessageThemeWathet    LarkMarkdownMessageTheme = "wathet"
	LarkMarkdownMessageThemeTurquoise LarkMarkdownMessageTheme = "turquoise"
	LarkMarkdownMessageThemeGreen     LarkMarkdownMessageTheme = "green"
	LarkMarkdownMessageThemeYellow    LarkMarkdownMessageTheme = "yellow"
	LarkMarkdownMessageThemeOrange    LarkMarkdownMessageTheme = "orange"
	LarkMarkdownMessageThemeRed       LarkMarkdownMessageTheme = "red"
	LarkMarkdownMessageThemeCarmine   LarkMarkdownMessageTheme = "carmine"
	LarkMarkdownMessageThemeViolet    LarkMarkdownMessageTheme = "violet"
	LarkMarkdownMessageThemePurple    LarkMarkdownMessageTheme = "purple"
	LarkMarkdownMessageThemeIndigo    LarkMarkdownMessageTheme = "indigo"
	LarkMarkdownMessageThemeGrey      LarkMarkdownMessageTheme = "grey"
)

var (
	supportedLarkMarkdownMessageTheme = map[string]LarkMarkdownMessageTheme{
		LarkMarkdownMessageThemeBlue.String():      LarkMarkdownMessageThemeBlue,
		LarkMarkdownMessageThemeWathet.String():    LarkMarkdownMessageThemeWathet,
		LarkMarkdownMessageThemeTurquoise.String(): LarkMarkdownMessageThemeTurquoise,
		LarkMarkdownMessageThemeGreen.String():     LarkMarkdownMessageThemeGreen,
		LarkMarkdownMessageThemeYellow.String():    LarkMarkdownMessageThemeYellow,
		LarkMarkdownMessageThemeOrange.String():    LarkMarkdownMessageThemeOrange,
		LarkMarkdownMessageThemeRed.String():       LarkMarkdownMessageThemeRed,
		LarkMarkdownMessageThemeCarmine.String():   LarkMarkdownMessageThemeCarmine,
		LarkMarkdownMessageThemeViolet.String():    LarkMarkdownMessageThemeViolet,
		LarkMarkdownMessageThemePurple.String():    LarkMarkdownMessageThemePurple,
		LarkMarkdownMessageThemeIndigo.String():    LarkMarkdownMessageThemeIndigo,
		LarkMarkdownMessageThemeGrey.String():      LarkMarkdownMessageThemeGrey,
	}
)

func (t LarkMarkdownMessageTheme) String() string { return string(t) }

func getLarkMarkdownMessageTheme(theme LarkMarkdownMessageTheme) string {
	if _, exist := supportedLarkMarkdownMessageTheme[theme.String()]; !exist {
		return LarkMarkdownMessageThemeBlue.String()
	} else {
		return theme.String()
	}
}

type LarkImageType string

const (
	LarkImageTypeMessage LarkImageType = "message"
	LarkImageTypeAvatar  LarkImageType = "avatar"
)

var (
	supportedLarkImageType = map[string]LarkImageType{
		LarkImageTypeMessage.String(): LarkImageTypeMessage,
		LarkImageTypeAvatar.String():  LarkImageTypeAvatar,
	}
)

func (t LarkImageType) String() string { return string(t) }

func getLarkImageType(imageType LarkImageType) string {
	if _, exist := supportedLarkImageType[imageType.String()]; !exist {
		return LarkImageTypeMessage.String()
	} else {
		return imageType.String()
	}
}

type LarkFileType string

const (
	LarkFileTypeOpus   LarkFileType = "opus"
	LarkFileTypeMp4    LarkFileType = "mp4"
	LarkFileTypePdf    LarkFileType = "pdf"
	LarkFileTypeDoc    LarkFileType = "doc"
	LarkFileTypeXls    LarkFileType = "xls"
	LarkFileTypePpt    LarkFileType = "ppt"
	LarkFileTypeStream LarkFileType = "stream"
)

var (
	supportedLarkFileType = map[string]LarkFileType{
		LarkFileTypeOpus.String():   LarkFileTypeOpus,
		LarkFileTypeMp4.String():    LarkFileTypeMp4,
		LarkFileTypePdf.String():    LarkFileTypePdf,
		LarkFileTypeDoc.String():    LarkFileTypeDoc,
		LarkFileTypeXls.String():    LarkFileTypeXls,
		LarkFileTypePpt.String():    LarkFileTypePpt,
		LarkFileTypeStream.String(): LarkFileTypeStream,
	}
)

func (t LarkFileType) String() string { return string(t) }

func getLarkFileType(fileType LarkFileType) string {
	if _, exist := supportedLarkFileType[fileType.String()]; !exist {
		return LarkFileTypeStream.String()
	} else {
		return fileType.String()
	}
}
