package im

import (
	"context"

	"github.com/alioth-center/infrastructure/thirdparty/lark/event"
)

// 此处定义了飞书机器人事件的数据结构和处理函数
// reference: https://open.feishu.cn/document/client-docs/bot-v3/events/menu

type BotEventOperatorID struct {
	UnionID string `json:"union_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"open_id,omitempty"`
}

type BotEventOperator struct {
	OperatorName string             `json:"operator_name"`
	OperatorID   BotEventOperatorID `json:"operator_id"`
}

type BotEvent struct {
	Operator  BotEventOperator `json:"operator"`
	EventKey  string           `json:"event_key"`
	Timestamp int64            `json:"timestamp"`
}

type BotEventHandler struct {
	event.BaseEventHandler[BotEvent]
}

func SetBotEventHandler(handler func(ctx context.Context, event *event.FullCallbackRequest[BotEvent]) (err error)) {
	if handler != nil {
		// 只设置非空的事件处理函数
		event.SetEventHandler(BotEventHandler{BaseEventHandler: event.InitializeBaseEventHandler(handler, event.EnumBotEvent)})
	}
}
