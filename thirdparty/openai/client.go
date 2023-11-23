package openai

import (
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/network/http"
	"strings"
)

const (
	userAgent = "alioth-center/http-client v1.2.1"
)

type Client interface {
	ListModels(req ListModelRequest) (resp ListModelResponseBody, err error)
	RetrieveModel(req RetrieveModelRequest) (resp RetrieveModelResponseBody, err error)
	GenerateImage(req CreateImageRequest) (resp ImageResponseBody, err error)
	CompleteChat(req CompleteChatRequest) (resp CompleteChatResponseBody, err error)
	CreateSpeech(req CreateSpeechRequest) (resp CreateSpeechResponseBody, err error)
	CreateTranscription(req CreateTranscriptionRequest) (resp CreateTranscriptionResponseBody, err error)
	CompleteModeration(req CompleteModerationRequest) (resp CompleteModerationResponseBody, err error)
}

type client struct {
	executor http.Client
	options  Config
}

func (c client) ListModels(_ ListModelRequest) (resp ListModelResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumListModel).
		SetUserAgent(userAgent).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return ListModelResponseBody{}, fmt.Errorf("execute list models request error: %w", executeErr)
	} else if status != http.StatusOK {
		return ListModelResponseBody{}, fmt.Errorf("execute list models request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return ListModelResponseBody{}, fmt.Errorf("parse list models response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) RetrieveModel(req RetrieveModelRequest) (resp RetrieveModelResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumRetrieveModel,
			func(original string) (result string) { return strings.ReplaceAll(original, "{:model:}", req.Model) },
		).
		SetUserAgent(userAgent).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return RetrieveModelResponseBody{}, fmt.Errorf("execute retrieve model request error: %w", executeErr)
	} else if status != http.StatusOK {
		return RetrieveModelResponseBody{}, fmt.Errorf("execute retrieve model request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return RetrieveModelResponseBody{}, fmt.Errorf("parse retrieve model response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) GenerateImage(req CreateImageRequest) (resp ImageResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateImage).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetJsonBody(&req.Body)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return ImageResponseBody{}, fmt.Errorf("execute generate image request error: %w", executeErr)
	} else if status != http.StatusOK {
		return ImageResponseBody{}, fmt.Errorf("execute generate image request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return ImageResponseBody{}, fmt.Errorf("parse generate image response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) CompleteChat(req CompleteChatRequest) (resp CompleteChatResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteChat).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetJsonBody(&req.Body)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CompleteChatResponseBody{}, fmt.Errorf("execute complete chat request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CompleteChatResponseBody{}, fmt.Errorf("execute complete chat request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return CompleteChatResponseBody{}, fmt.Errorf("parse complete chat response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) CreateSpeech(req CreateSpeechRequest) (resp CreateSpeechResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateSpeech).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetJsonBody(&req.Body)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CreateSpeechResponseBody{}, fmt.Errorf("execute create speech request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CreateSpeechResponseBody{}, fmt.Errorf("execute create speech request error: %d: %s", status, string(body))
	} else {
		return body, nil
	}
}

func (c client) CreateTranscription(req CreateTranscriptionRequest) (resp CreateTranscriptionResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateTranscription).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetMultiPartBody("file", req.FormBody.FileName, req.FormBody.File, req.FormBody.ToMultiPartBody())
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("execute create transcription request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("execute create transcription request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("parse create transcription response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) CompleteModeration(req CompleteModerationRequest) (resp CompleteModerationResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteModeration).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetJsonBody(&req.Body)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CompleteModerationResponseBody{}, fmt.Errorf("execute complete moderation request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CompleteModerationResponseBody{}, fmt.Errorf("execute complete moderation request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return CompleteModerationResponseBody{}, fmt.Errorf("parse complete moderation response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func NewClient(options Config) Client {
	return client{
		executor: http.NewClient(),
		options:  options,
	}
}
