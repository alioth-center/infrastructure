package im

import "github.com/alioth-center/infrastructure/thirdparty/lark/event"

// 此处定义了飞书消息接收事件的数据结构和处理函数
// reference: https://open.feishu.cn/document/server-docs/im-v1/message/events/receive

// MessageReceiverUserID 飞书用户ID组，包含union_id、user_id、open_id，三者至少有一个不为空
type MessageReceiverUserID struct {
	UnionID string `json:"union_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
}

type MessageReceiveEventSender struct {
	SenderID   MessageReceiverUserID `json:"sender_id"`
	SenderType string                `json:"sender_type"`
	TenantKey  string                `json:"tenant_key"`
}

type MessageReceiveEventMessageMentionItem struct {
	Key       string                `json:"key"`
	ID        MessageReceiverUserID `json:"id"`
	Name      string                `json:"name"`
	TenantKey string                `json:"tenant_key"`
}

type MessageReceiveEventMessage struct {
	MessageID   string                                  `json:"message_id"`
	RootID      string                                  `json:"root_id"`
	ParentID    string                                  `json:"parent_id"`
	CreateTime  string                                  `json:"create_time"`
	UpdateTime  string                                  `json:"update_time"`
	ChatID      string                                  `json:"chat_id"`
	ChatType    string                                  `json:"chat_type"`
	MessageType string                                  `json:"message_type"`
	Content     string                                  `json:"content"`
	Mentions    []MessageReceiveEventMessageMentionItem `json:"mentions"`
	UserAgent   string                                  `json:"user_agent"`
}

type MessageReceiveEvent struct {
	Sender  MessageReceiveEventSender  `json:"sender"`
	Message MessageReceiveEventMessage `json:"message"`
}

type MessageReceiveEventHandler struct {
	event.BaseEventHandler[MessageReceiveEvent]
}

func SetMessageHandler(handler func(event *event.FullCallbackRequest[MessageReceiveEvent]) (status int)) {
	if handler != nil {
		// 只设置非空的事件处理函数
		event.SetEventHandler(MessageReceiveEventHandler{BaseEventHandler: event.InitializeBaseEventHandler(handler, event.EnumReceiveMessage)})
	}
}
