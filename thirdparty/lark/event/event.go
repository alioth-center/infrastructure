package event

import (
	"context"
	"encoding/json"
	"fmt"
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
	HandleEvent(ctx context.Context, event *CallbackRequest) (err error)
}

type Enum string

func (e Enum) TriggerEventType() string { return string(e) }

const (
	EnumReceiveMessage Enum = "im.message.receive_v1"
	EnumBotEvent       Enum = "application.bot.menu_v6"
	EnumMemberAdd      Enum = "im.chat.member.user.added_v1"
	EnumMemberDelete   Enum = "im.chat.member.user.deleted_v1"
)

func HandleEvent(ctx context.Context, req *CallbackRequest) (err error) {
	if handler, exist := eventHandlers[req.Header.EventType]; exist && handler != nil {
		// 根据事件类型，调用对应的事件处理函数
		return handler.HandleEvent(ctx, req)
	} else if handler == nil {
		// 事件处理函数未实现，返回错误
		return NewHandlerUnimplementedError(req.Header.EventType)
	} else {
		// 事件类型不存在，返回错误
		return NewHandlerNotFoundError(req.Header.EventType)
	}
}

func HandleEncryptedEvent(ctx context.Context, req *EncryptedRequest, secretKey string) (err error) {
	if secretKey == "" {
		// 未设置密钥，返回错误
		return NewSecretKeyUnsetError()
	}

	decrypted, decryptErr := decrypt(req.Encrypt, secretKey)
	if decryptErr != nil {
		// 解密失败，返回错误
		return fmt.Errorf("failed to decrypt event: %w", decryptErr)
	}

	var request CallbackRequest
	unmarshalErr := json.Unmarshal(decrypted, &request)
	if unmarshalErr != nil {
		// 解析失败，返回错误
		return fmt.Errorf("failed to unmarshal decrypted event: %w", unmarshalErr)
	}

	return HandleEvent(ctx, &request)
}

func HandleEncryptedEventWithChallenge(ctx context.Context, req *EncryptedRequest, secretKey, verificationToken string, callback func(challenge ChallengeResponse) (err error)) (err error) {
	if verificationToken == "" {
		// 未设置验证 token，返回错误
		return NewVerificationTokenUnsetError()
	} else if secretKey == "" {
		// 未设置密钥，返回错误
		return NewSecretKeyUnsetError()
	}

	decrypted, decryptErr := decrypt(req.Encrypt, secretKey)
	if decryptErr != nil {
		// 解密失败，返回错误
		return fmt.Errorf("failed to decrypt event: %w", decryptErr)
	}

	var request CallbackRequest
	unmarshalErr := json.Unmarshal(decrypted, &request)
	if unmarshalErr != nil {
		// 解析失败，返回错误
		return fmt.Errorf("failed to unmarshal decrypted event: %w", unmarshalErr)
	}

	if request.Header.EventType == "" {
		// 事件类型为空，说明是 challenge 请求
		challenge, parseErr := ParseChallengeWithEncryption(decrypted, verificationToken, secretKey)
		if parseErr != nil {
			// 解析失败，返回错误
			return fmt.Errorf("failed to parse challenge: %w", parseErr)
		} else if challenge.Challenge == "" {
			// challenge 为空，返回错误
			return NewChallengeEmptyError()
		}

		// 返回 challenge 响应
		return callback(challenge)
	} else {
		// 返回事件处理结果
		return HandleEvent(ctx, &request)
	}
}

type BaseEventHandler[realEvent any] struct {
	handler   func(ctx context.Context, event *FullCallbackRequest[realEvent]) (err error)
	eventType string
}

func (b BaseEventHandler[realEvent]) TriggerEventType() string {
	return b.eventType
}

func (b BaseEventHandler[realEvent]) HandleEvent(ctx context.Context, event *CallbackRequest) (err error) {
	if event.Header.EventType != b.TriggerEventType() {
		// 事件类型不匹配，返回错误
		return NewHandlerTypeNotMatchError(event.Header.EventType, b.TriggerEventType())
	} else if b.handler == nil {
		// 事件处理函数未实现，返回错误
		return NewHandlerUnimplementedError(event.Header.EventType)
	} else {
		defer func() {
			if finalErr := concurrency.RecoverErr(recover()); finalErr != nil {
				// 事件处理函数内部发生错误，重写错误信息
				err = fmt.Errorf("failed to handle event %s: %w", event.Header.EventType, finalErr)
			}
		}()

		if fullData, getFullDataErr := GetCallbackRequestEventData[realEvent](event, b); getFullDataErr != nil {
			// 事件内容不匹配，返回错误
			return NewTypeNotMatchError(event.Header.EventType, b.TriggerEventType())
		} else {
			// 调用事件处理函数
			return b.handler(ctx, &fullData)
		}
	}
}

func InitializeBaseEventHandler[realEvent any](fn func(ctx context.Context, event *FullCallbackRequest[realEvent]) (err error), eventType Enum) BaseEventHandler[realEvent] {
	return BaseEventHandler[realEvent]{
		handler:   fn,
		eventType: eventType.TriggerEventType(),
	}
}

type HandlerUnimplementedError struct {
	Event string `json:"event"`
}

func (h HandlerUnimplementedError) Error() string {
	return "handler for event " + h.Event + " is not implemented"
}

func NewHandlerUnimplementedError(event string) HandlerUnimplementedError {
	return HandlerUnimplementedError{
		Event: event,
	}
}

type HandlerNotFoundError struct {
	Event string `json:"event"`
}

func (h HandlerNotFoundError) Error() string {
	return "handler for event " + h.Event + " not found"
}

func NewHandlerNotFoundError(event string) HandlerNotFoundError {
	return HandlerNotFoundError{
		Event: event,
	}
}

type HandlerTypeNotMatchError struct {
	Event       string `json:"event"`
	HandlerType string `json:"handler_type"`
}

func (h HandlerTypeNotMatchError) Error() string {
	return "handler for event " + h.Event + " is not match with " + h.HandlerType
}

func NewHandlerTypeNotMatchError(event, handlerType string) HandlerTypeNotMatchError {
	return HandlerTypeNotMatchError{
		Event:       event,
		HandlerType: handlerType,
	}
}

type SecretKeyUnsetError struct{}

func (s SecretKeyUnsetError) Error() string {
	return "secret key is not set"
}

func NewSecretKeyUnsetError() SecretKeyUnsetError {
	return SecretKeyUnsetError{}
}

type VerificationTokenUnsetError struct{}

func (v VerificationTokenUnsetError) Error() string {
	return "verification token is not set"
}

func NewVerificationTokenUnsetError() VerificationTokenUnsetError {
	return VerificationTokenUnsetError{}
}

type ChallengeEmptyError struct{}

func (c ChallengeEmptyError) Error() string {
	return "challenge is empty"
}

func NewChallengeEmptyError() ChallengeEmptyError {
	return ChallengeEmptyError{}
}
