package openai

import (
	"fmt"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/values"
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
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return ListModelResponseBody{}, fmt.Errorf("execute list models request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return ListModelResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return ListModelResponseBody{}, fmt.Errorf("parse list models response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) RetrieveModel(req RetrieveModelRequest) (resp RetrieveModelResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumRetrieveModel, "model", req.Model).
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return RetrieveModelResponseBody{}, fmt.Errorf("execute retrieve model request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return RetrieveModelResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return RetrieveModelResponseBody{}, fmt.Errorf("parse retrieve model response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) GenerateImage(req CreateImageRequest) (resp ImageResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateImage).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return ImageResponseBody{}, fmt.Errorf("execute generate image request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return ImageResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return ImageResponseBody{}, fmt.Errorf("parse generate image response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CompleteChat(req CompleteChatRequest) (resp CompleteChatResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteChat).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CompleteChatResponseBody{}, fmt.Errorf("execute complete chat request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CompleteChatResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CompleteChatResponseBody{}, fmt.Errorf("parse complete chat response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CreateSpeech(req CreateSpeechRequest) (resp CreateSpeechResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateSpeech).
		WithMethod(http.POST).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CreateSpeechResponseBody{}, fmt.Errorf("execute create speech request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CreateSpeechResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CreateSpeechResponseBody{}, fmt.Errorf("parse create speech response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CreateTranscription(req CreateTranscriptionRequest) (resp CreateTranscriptionResponseBody, err error) {
	multipart := http.NewMultipartBodyBuilder().WithFile("file", req.FormBody.FileName, req.FormBody.File)
	for k, v := range req.FormBody.ToMultiPartBody() {
		multipart = multipart.WithForm(k, v)
	}
	multipartBody, contentType, buildMultipartErr := multipart.Build()
	if buildMultipartErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("build create transcription multipart body error: %w", buildMultipartErr)
	}

	request := c.options.buildBaseRequest(EndpointEnumCreateTranscription).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithBody(multipartBody).
		WithContentType(contentType)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("execute create transcription request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CreateTranscriptionResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("parse create transcription response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CompleteModeration(req CompleteModerationRequest) (resp CompleteModerationResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteModeration).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CompleteModerationResponseBody{}, fmt.Errorf("execute complete moderation request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CompleteModerationResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CompleteModerationResponseBody{}, fmt.Errorf("parse complete moderation response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CreateFineTuningJob(req CreateFineTuningJobRequest) (resp CreateFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumCreateFineTuningJob).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CreateFineTuningJobResponseBody{}, fmt.Errorf("execute create fine tuning job request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CreateFineTuningJobResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CreateFineTuningJobResponseBody{}, fmt.Errorf("parse create fine tuning job response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) RetrieveFineTuningJob(req RetrieveFineTuningJobRequest) (resp RetrieveFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumRetrieveFineTuningJob, "id", req.Body.ID).
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return RetrieveFineTuningJobResponseBody{}, fmt.Errorf("execute retrieve fine tuning job request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return RetrieveFineTuningJobResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return RetrieveFineTuningJobResponseBody{}, fmt.Errorf("parse retrieve fine tuning job response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) ListFineTuningJobs(req ListFineTuningJobsRequest) (resp ListFineTuningJobsResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumListFineTuningJobs).
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson).
		WithQueryIgnoreEmptyValue("after", req.Body.After).
		WithQueryIgnoreEmptyValue("limit", values.IntToString(req.Body.Limit))
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return ListFineTuningJobsResponseBody{}, fmt.Errorf("execute list fine tuning jobs request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return ListFineTuningJobsResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return ListFineTuningJobsResponseBody{}, fmt.Errorf("parse list fine tuning jobs response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CancelFineTuningJob(req CancelFineTuningJobRequest) (resp CancelFineTuningJobResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumCancelFineTuningJob, "id", req.Body.ID).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return CancelFineTuningJobResponseBody{}, fmt.Errorf("execute cancel fine tuning job request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return CancelFineTuningJobResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return CancelFineTuningJobResponseBody{}, fmt.Errorf("parse cancel fine tuning job response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) UploadFile(req UploadFileRequest) (resp UploadFileResponseBody, err error) {
	multipart := http.NewMultipartBodyBuilder().WithFile("file", req.FormBody.FileName, req.FormBody.File)
	for k, v := range req.FormBody.ToMultiPartBody() {
		multipart = multipart.WithForm(k, v)
	}
	multipartBody, contentType, buildMultipartErr := multipart.Build()
	if buildMultipartErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("build upload file multipart body error: %w", buildMultipartErr)
	}

	request := c.options.
		buildBaseRequest(EndpointEnumUploadFile).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithBody(multipartBody).
		WithContentType(contentType)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("execute upload file request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return UploadFileResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("parse upload file response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) ListFiles(req ListFilesRequest) (resp ListFilesResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumListFiles).
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson).
		WithUserAgent(c.options.getUserAgent()).
		WithQueryIgnoreEmptyValue("purpose", req.Body.Purpose)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return ListFilesResponseBody{}, fmt.Errorf("execute list files request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return ListFilesResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return ListFilesResponseBody{}, fmt.Errorf("parse list files response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) DeleteFile(req DeleteFileRequest) (resp DeleteFileResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumDeleteFile, "id", req.Body.ID).
		WithMethod(http.DELETE).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return DeleteFileResponseBody{}, fmt.Errorf("execute delete file request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return DeleteFileResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return DeleteFileResponseBody{}, fmt.Errorf("parse delete file response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) RetrieveFile(req RetrieveFileRequest) (resp RetrieveFileResponseBody, err error) {
	request := c.options.
		buildBaseRequest(EndpointEnumRetrieveFile, "id", req.Body.ID).
		WithMethod(http.GET).
		WithAccept(http.ContentTypeJson)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return RetrieveFileResponseBody{}, fmt.Errorf("execute retrieve file request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return RetrieveFileResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return RetrieveFileResponseBody{}, fmt.Errorf("parse retrieve file response error: %w", bindErr)
	}

	return resp, nil
}

func NewClient(options Config, logger logger.Logger) Client {
	return client{
		executor: http.NewLoggerClient(logger),
		options:  options,
	}
}
