package event

import (
	"encoding/json"
	"fmt"

	"github.com/alioth-center/infrastructure/utils/values"
)

type CallbackRequestHeader struct {
	EventID    string `json:"event_id"`    // 事件ID
	Token      string `json:"token"`       // 验证token
	CreateTime string `json:"create_time"` // 事件创建时间，毫秒级时间戳
	EventType  string `json:"event_type"`  // 事件类型
	TenantKey  string `json:"tenant_key"`  // 企业ID
	AppID      string `json:"app_id"`      // 应用ID
}

type CallbackRequest struct {
	Schema string                `json:"schema,omitempty"` // 版本，只有2.0的版本有这个字段
	Header CallbackRequestHeader `json:"header"`           // 请求头
	Event  any                   `json:"event,omitempty"`  // 事件内容
}

type FullCallbackRequest[eventData any] struct {
	Schema string                `json:"schema,omitempty"` // 版本，只有2.0的版本有这个字段
	Header CallbackRequestHeader `json:"header"`           // 请求头
	Event  eventData             `json:"event,omitempty"`  // 事件内容
}

func GetCallbackRequestEventData[data any](request *CallbackRequest, event Handler) (fullData FullCallbackRequest[data], err error) {
	if request.Header.EventType != event.TriggerEventType() {
		// 事件类型不匹配，返回错误
		return values.Nil[FullCallbackRequest[data]](), NewTypeNotMatchError(event.TriggerEventType(), request.Header.EventType)
	}

	var resultData data
	if request.Event != nil {
		// 当且仅当事件内容不为空时，才进行序列化和反序列化
		if marshalBytes, marshalErr := json.Marshal(&request.Event); marshalErr != nil {
			return values.Nil[FullCallbackRequest[data]](), fmt.Errorf("failed to marshal event: %w", marshalErr)
		} else if unmarshalErr := json.Unmarshal(marshalBytes, &resultData); unmarshalErr != nil {
			return values.Nil[FullCallbackRequest[data]](), fmt.Errorf("failed to unmarshal event: %w", unmarshalErr)
		} else {
			result := FullCallbackRequest[data]{
				Schema: request.Schema,
				Header: request.Header,
				Event:  resultData,
			}

			return result, nil
		}
	} else {
		// 当事件内容为空时，直接返回
		result := FullCallbackRequest[data]{
			Schema: request.Schema,
			Header: request.Header,
		}

		return result, nil
	}
}

type TypeNotMatchError struct {
	ExpectedType string
	ActualType   string
}

func (e TypeNotMatchError) Error() string {
	return fmt.Sprintf("type not match, expected: %s, actual: %s", e.ExpectedType, e.ActualType)
}

func NewTypeNotMatchError(expectedType string, actualType string) error {
	return TypeNotMatchError{
		ExpectedType: expectedType,
		ActualType:   actualType,
	}
}
