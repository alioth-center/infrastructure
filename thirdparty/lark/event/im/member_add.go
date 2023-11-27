package im

import event "github.com/alioth-center/infrastructure/thirdparty/lark/event"

// 此处定义了飞书群聊用户进群事件的数据结构和处理函数
// reference: https://open.feishu.cn/document/server-docs/group/chat-member/event/added

type MemberAddUserID struct {
	UnionID string `json:"union_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
}

type MemberAddUser struct {
	Name      string          `json:"name"`
	TenantKey string          `json:"tenant_key"`
	UserID    MemberAddUserID `json:"user_id"`
}

type MemberAddGroupNameTranslation struct {
	ZhCn string `json:"zh_cn"`
	EnUs string `json:"en_us"`
	JaJp string `json:"ja_jp"`
}

type MemberAddEvent struct {
	ChatID            string                        `json:"chat_id"`
	OperatorID        MemberAddUserID               `json:"operator_id"`
	External          bool                          `json:"external"`
	OperatorTenantKey string                        `json:"operator_tenant_key"`
	Users             []MemberAddUser               `json:"users"`
	Name              string                        `json:"name"`
	I18NNames         MemberAddGroupNameTranslation `json:"i18n_names"`
}

type MemberAddEventHandler struct {
	event.BaseEventHandler[MemberAddEvent]
}

func SetMemberAddHandler(handler func(event *event.FullCallbackRequest[MemberAddEvent]) (status int)) {
	if handler != nil {
		// 只设置非空的事件处理函数
		event.SetEventHandler(MemberAddEventHandler{BaseEventHandler: event.InitializeBaseEventHandler(handler, event.EnumMemberAdd)})
	}
}
