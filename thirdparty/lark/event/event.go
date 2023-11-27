package event

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

var (
	eventHandlers = map[string]Handler{
		EnumReceiveMessage.TriggerEventType(): nil,
		EnumBotEvent.TriggerEventType():       nil,
		EnumMemberAdd.TriggerEventType():      nil,
		EnumMemberDelete.TriggerEventType():   nil,
	}
)

func SetEventHandler(handler Handler) {
	if handler != nil {
		if evt, exist := eventHandlers[handler.TriggerEventType()]; exist && evt == nil {
			// 只有当事件处理函数未实现时，并且传入的事件处理函数不为空时，才进行设置
			eventHandlers[handler.TriggerEventType()] = handler
		}
	}
}

type Handler interface {
	TriggerEventType() string
	HandleEvent(event *CallbackRequest) (status int)
}

type Enum string

func (e Enum) TriggerEventType() string { return string(e) }

const (
	EnumReceiveMessage Enum = "im.message.receive_v1"
	EnumBotEvent       Enum = "application.bot.menu_v6"
	EnumMemberAdd      Enum = "im.chat.member.user.added_v1"
	EnumMemberDelete   Enum = "im.chat.member.user.deleted_v1"
)

func HandleEvent(req *CallbackRequest) (status int) {
	if handler, ok := eventHandlers[req.Header.EventType]; ok && handler != nil {
		// 根据事件类型，调用对应的事件处理函数
		return handler.HandleEvent(req)
	} else if handler == nil {
		// 事件处理函数未实现，返回错误
		return http.StatusNotImplemented
	} else {
		// 事件类型不存在，返回错误
		return http.StatusNotFound
	}
}

type BaseEventHandler[realEvent any] struct {
	handler   func(event *FullCallbackRequest[realEvent]) (status int)
	eventType string
}

func (b BaseEventHandler[realEvent]) TriggerEventType() string {
	return b.eventType
}

func (b BaseEventHandler[realEvent]) HandleEvent(event *CallbackRequest) (status int) {
	if event.Header.EventType != b.TriggerEventType() {
		// 事件类型不匹配，返回错误
		return http.StatusBadRequest
	} else if b.handler == nil {
		// 事件处理函数未实现，返回错误
		return http.StatusNotImplemented
	} else {
		defer func() {
			if finalErr := concurrency.RecoverErr(recover()); finalErr != nil {
				// 事件处理函数内部发生错误，返回错误
				status = http.StatusServiceUnavailable
			}
		}()

		if fullData, getFullDataErr := GetCallbackRequestEventData[realEvent](event, b); getFullDataErr != nil {
			// 事件内容不匹配，返回错误
			return http.StatusBadRequest
		} else {
			// 调用事件处理函数
			return b.handler(&fullData)
		}
	}
}

func InitializeBaseEventHandler[realEvent any](fn func(event *FullCallbackRequest[realEvent]) (status int), eventType Enum) BaseEventHandler[realEvent] {
	return BaseEventHandler[realEvent]{
		handler:   fn,
		eventType: eventType.TriggerEventType(),
	}
}
