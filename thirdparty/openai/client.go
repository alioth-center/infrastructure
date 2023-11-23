package openai

import (
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/network/http"
	"net/url"
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

	CreateFineTuningJob(req CreateFineTuningJobRequest) (resp CreateFineTuningJobResponseBody, err error)
	RetrieveFineTuningJob(req RetrieveFineTuningJobRequest) (resp RetrieveFineTuningJobResponseBody, err error)
	ListFineTuningJobs(req ListFineTuningJobsRequest) (resp ListFineTuningJobsResponseBody, err error)
	CancelFineTuningJob(req CancelFineTuningJobRequest) (resp CancelFineTuningJobResponseBody, err error)
	UploadFile(req UploadFileRequest) (resp UploadFileResponseBody, err error)
	ListFiles(req ListFilesRequest) (resp ListFilesResponseBody, err error)
	DeleteFile(req DeleteFileRequest) (resp DeleteFileResponseBody, err error)
	RetrieveFile(req RetrieveFileRequest) (resp RetrieveFileResponseBody, err error)
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

func (c client) CreateFineTuningJob(req CreateFineTuningJobRequest) (resp CreateFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumCreateFineTuningJob).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetJsonBody(&req.Body)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CreateFineTuningJobResponseBody{}, fmt.Errorf("execute create fine tuning job request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CreateFineTuningJobResponseBody{}, fmt.Errorf("execute create fine tuning job request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return CreateFineTuningJobResponseBody{}, fmt.Errorf("parse create fine tuning job response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) RetrieveFineTuningJob(req RetrieveFineTuningJobRequest) (resp RetrieveFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumRetrieveFineTuningJob,
			func(original string) (result string) { return strings.ReplaceAll(original, "{:id:}", req.Body.ID) },
		).
		SetUserAgent(userAgent).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return RetrieveFineTuningJobResponseBody{}, fmt.Errorf("execute retrieve fine tuning job request error: %w", executeErr)
	} else if status != http.StatusOK {
		return RetrieveFineTuningJobResponseBody{}, fmt.Errorf("execute retrieve fine tuning job request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return RetrieveFineTuningJobResponseBody{}, fmt.Errorf("parse retrieve fine tuning job response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) ListFineTuningJobs(req ListFineTuningJobsRequest) (resp ListFineTuningJobsResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumListFineTuningJobs,
			func(original string) (result string) {
				// build query string from request {after,limit}
				queryParams := url.Values{}
				if req.Body.After != "" {
					queryParams.Set("after", req.Body.After)
				}
				if req.Body.Limit != 0 {
					queryParams.Set("limit", fmt.Sprintf("%d", req.Body.Limit))
				}
				if len(queryParams) == 0 {
					return original
				} else {
					return fmt.Sprintf("%s?%s", original, queryParams.Encode())
				}
			},
		).
		SetUserAgent(userAgent).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return ListFineTuningJobsResponseBody{}, fmt.Errorf("execute list fine tuning jobs request error: %w", executeErr)
	} else if status != http.StatusOK {
		return ListFineTuningJobsResponseBody{}, fmt.Errorf("execute list fine tuning jobs request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return ListFineTuningJobsResponseBody{}, fmt.Errorf("parse list fine tuning jobs response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) CancelFineTuningJob(req CancelFineTuningJobRequest) (resp CancelFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumCancelFineTuningJob,
			func(original string) (result string) { return strings.ReplaceAll(original, "{:id:}", req.Body.ID) },
		).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return CancelFineTuningJobResponseBody{}, fmt.Errorf("execute cancel fine tuning job request error: %w", executeErr)
	} else if status != http.StatusOK {
		return CancelFineTuningJobResponseBody{}, fmt.Errorf("execute cancel fine tuning job request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return CancelFineTuningJobResponseBody{}, fmt.Errorf("parse cancel fine tuning job response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) UploadFile(req UploadFileRequest) (resp UploadFileResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumUploadFile).
		SetUserAgent(userAgent).
		SetMethod(http.POST).
		SetAccept(http.ContentTypeJson).
		SetMultiPartBody("file", req.FormBody.FileName, req.FormBody.File, req.FormBody.ToMultiPartBody())
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("execute upload file request error: %w", executeErr)
	} else if status != http.StatusOK {
		return UploadFileResponseBody{}, fmt.Errorf("execute upload file request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("parse upload file response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) ListFiles(req ListFilesRequest) (resp ListFilesResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumListFiles).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson).
		SetUserAgent(userAgent)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return ListFilesResponseBody{}, fmt.Errorf("execute list files request error: %w", executeErr)
	} else if status != http.StatusOK {
		return ListFilesResponseBody{}, fmt.Errorf("execute list files request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return ListFilesResponseBody{}, fmt.Errorf("parse list files response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) DeleteFile(req DeleteFileRequest) (resp DeleteFileResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumDeleteFile,
			func(original string) (result string) { return strings.ReplaceAll(original, "{:id:}", req.Body.ID) },
		).
		SetMethod(http.DELETE).
		SetAccept(http.ContentTypeJson).
		SetUserAgent(userAgent)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return DeleteFileResponseBody{}, fmt.Errorf("execute delete file request error: %w", executeErr)
	} else if status != http.StatusOK {
		return DeleteFileResponseBody{}, fmt.Errorf("execute delete file request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return DeleteFileResponseBody{}, fmt.Errorf("parse delete file response error: %w", parseErr)
	} else {
		return resp, nil
	}
}

func (c client) RetrieveFile(req RetrieveFileRequest) (resp RetrieveFileResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumRetrieveFile,
			func(original string) (result string) { return strings.ReplaceAll(original, "{:id:}", req.Body.ID) },
		).
		SetMethod(http.GET).
		SetAccept(http.ContentTypeJson).
		SetUserAgent(userAgent)
	status, body, executeErr := c.executor.ExecuteRequest(request).Result()
	if executeErr != nil {
		return RetrieveFileResponseBody{}, fmt.Errorf("execute retrieve file request error: %w", executeErr)
	} else if status != http.StatusOK {
		return RetrieveFileResponseBody{}, fmt.Errorf("execute retrieve file request error: %d: %s", status, string(body))
	} else if parseErr := json.Unmarshal(body, &resp); parseErr != nil {
		return RetrieveFileResponseBody{}, fmt.Errorf("parse retrieve file response error: %w", parseErr)
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
