package im

import (
	"context"

	"github.com/alioth-center/infrastructure/thirdparty/lark/event"
)

// 此处定义了飞书群聊用户被删除事件的数据结构和处理函数
// reference: https://open.feishu.cn/document/server-docs/group/chat-member/event/deleted-2

type MemberDeleteUserID struct {
	UnionID string `json:"union_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
}

type MemberDeleteUser struct {
	Name      string             `json:"name"`
	TenantKey string             `json:"tenant_key"`
	UserID    MemberDeleteUserID `json:"user_id"`
}

type MemberDeleteGroupNameTranslation struct {
	ZhCn string `json:"zh_cn,omitempty"`
	EnUs string `json:"en_us,omitempty"`
	JaJp string `json:"ja_jp,omitempty"`
}

type MemberDeleteEvent struct {
	ChatID            string                           `json:"chat_id"`
	OperatorID        MemberDeleteUserID               `json:"operator_id"`
	External          bool                             `json:"external"`
	OperatorTenantKey string                           `json:"operator_tenant_key"`
	Users             []MemberDeleteUser               `json:"users"`
	Name              string                           `json:"name"`
	I18NNames         MemberDeleteGroupNameTranslation `json:"i18n_names"`
}

type MemberDeleteEventHandler struct {
	event.BaseEventHandler[MemberDeleteEvent]
}

func SetMemberDeleteHandler(handler func(ctx context.Context, event *event.FullCallbackRequest[MemberDeleteEvent]) (err error)) {
	if handler != nil {
		// 只设置非空的事件处理函数
		event.SetEventHandler(MemberDeleteEventHandler{BaseEventHandler: event.InitializeBaseEventHandler(handler, event.EnumMemberDelete)})
	}
}
