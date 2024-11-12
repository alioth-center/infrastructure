package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/pandodao/tokenizer-go"
)

type Client interface {
	CalculateToken(inputs ...string) (tokens int)
	ListModels(ctx context.Context, req ListModelRequest) (resp ListModelResponseBody, err error)
	RetrieveModel(ctx context.Context, req RetrieveModelRequest) (resp RetrieveModelResponseBody, err error)
	GenerateImage(ctx context.Context, req CreateImageRequest) (resp ImageResponseBody, err error)
	CompleteChat(ctx context.Context, req CompleteChatRequest) (resp CompleteChatResponseBody, err error)
	CompleteStreamingChat(ctx context.Context, req CompleteChatRequest) (events <-chan StreamingReplyObject, err error)
	CreateSpeech(ctx context.Context, req CreateSpeechRequest) (resp CreateSpeechResponseBody, err error)
	CreateTranscription(ctx context.Context, req CreateTranscriptionRequest) (resp CreateTranscriptionResponseBody, err error)
	CompleteModeration(ctx context.Context, req CompleteModerationRequest) (resp CompleteModerationResponseBody, err error)
	Embedding(ctx context.Context, req EmbeddingRequest) (resp EmbeddingResponseBody, err error)

	CreateFineTuningJob(ctx context.Context, req CreateFineTuningJobRequest) (resp CreateFineTuningJobResponseBody, err error)
	RetrieveFineTuningJob(ctx context.Context, req RetrieveFineTuningJobRequest) (resp RetrieveFineTuningJobResponseBody, err error)
	ListFineTuningJobs(ctx context.Context, req ListFineTuningJobsRequest) (resp ListFineTuningJobsResponseBody, err error)
	CancelFineTuningJob(ctx context.Context, req CancelFineTuningJobRequest) (resp CancelFineTuningJobResponseBody, err error)
	UploadFile(ctx context.Context, req UploadFileRequest) (resp UploadFileResponseBody, err error)
	ListFiles(ctx context.Context, req ListFilesRequest) (resp ListFilesResponseBody, err error)
	DeleteFile(ctx context.Context, req DeleteFileRequest) (resp DeleteFileResponseBody, err error)
	RetrieveFile(ctx context.Context, req RetrieveFileRequest) (resp RetrieveFileResponseBody, err error)
}

type client struct {
	executor http.Client
	options  Config
	logger   logger.Logger
}

func (c client) CalculateToken(inputs ...string) (tokens int) {
	for _, input := range inputs {
		tokens += tokenizer.MustCalToken(input)
	}

	return tokens
}

func (c client) ListModels(ctx context.Context, _ ListModelRequest) (resp ListModelResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumListModel).
		WithContext(ctx).
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

func (c client) RetrieveModel(ctx context.Context, req RetrieveModelRequest) (resp RetrieveModelResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumRetrieveModel, "model", req.Model).
		WithContext(ctx).
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

func (c client) GenerateImage(ctx context.Context, req CreateImageRequest) (resp ImageResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateImage).
		WithContext(ctx).
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

func (c client) CompleteChat(ctx context.Context, req CompleteChatRequest) (resp CompleteChatResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteChat).
		WithContext(ctx).
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

func (c client) CompleteStreamingChat(ctx context.Context, req CompleteChatRequest) (events <-chan StreamingReplyObject, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteChat).
		WithContext(ctx).
		WithMethod(http.POST).
		WithJsonBody(&req.Body)
	rawRequest, buildErr := request.Build()
	if buildErr != nil {
		return nil, fmt.Errorf("build complete chat request error: %w", buildErr)
	}

	response, executeErr := c.executor.ExecuteRawRequest(rawRequest)
	if executeErr != nil {
		return nil, fmt.Errorf("execute complete chat request error: %w", executeErr)
	}

	if response == nil || response.StatusCode != http.StatusOK {
		return nil, errors.New("complete chat response status code is not 200")
	}

	result := make(chan StreamingReplyObject, 256)
	go func(events chan StreamingReplyObject, body io.ReadCloser) {
		defer close(events)

		for event := range http.ParseServerSentEventFromBody(body, 4096, 256) {
			reply := StreamingReplyObject{}
			payload := event.Data

			// end of the conversation
			if strings.TrimSpace(string(payload)) == "[DONE]" {
				break
			}

			if unmarshalErr := json.Unmarshal(payload, &reply); unmarshalErr != nil {
				c.logger.Error(logger.NewFields(ctx).WithMessage("unmarshal complete chat response error").WithData(map[string]any{"error": unmarshalErr, "event": event}))
				continue
			}

			events <- reply
		}

		_ = body.Close()
	}(result, response.Body)

	return result, nil
}

func (c client) CreateSpeech(ctx context.Context, req CreateSpeechRequest) (resp CreateSpeechResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateSpeech).
		WithContext(ctx).
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

	resp = response.RawBody()
	return resp, nil
}

func (c client) CreateTranscription(ctx context.Context, req CreateTranscriptionRequest) (resp CreateTranscriptionResponseBody, err error) {
	multipart := http.NewMultipartBodyBuilder().WithFile("file", req.FormBody.FileName, req.FormBody.File)
	for k, v := range req.FormBody.ToMultiPartBody() {
		multipart = multipart.WithForm(k, v)
	}
	multipartBody, contentType, buildMultipartErr := multipart.Build()
	if buildMultipartErr != nil {
		return CreateTranscriptionResponseBody{}, fmt.Errorf("build create transcription multipart body error: %w", buildMultipartErr)
	}

	request := c.options.buildBaseRequest(EndpointEnumCreateTranscription).
		WithContext(ctx).
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

func (c client) CompleteModeration(ctx context.Context, req CompleteModerationRequest) (resp CompleteModerationResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCompleteModeration).
		WithContext(ctx).
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

func (c client) Embedding(ctx context.Context, req EmbeddingRequest) (resp EmbeddingResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumEmbedding).
		WithContext(ctx).
		WithMethod(http.POST).
		WithAccept(http.ContentTypeJson).
		WithJsonBody(&req.Body)
	response, executeErr := c.executor.ExecuteRequest(request)
	if executeErr != nil {
		return EmbeddingResponseBody{}, fmt.Errorf("execute embedding request error: %w", executeErr)
	}

	code, message := response.Status()
	if code != http.StatusOK {
		return EmbeddingResponseBody{}, &ResponseStatusError{StatusCode: code, Status: message}
	}

	bindErr := response.BindJson(&resp)
	if bindErr != nil {
		return EmbeddingResponseBody{}, fmt.Errorf("parse embedding response error: %w", bindErr)
	}

	return resp, nil
}

func (c client) CreateFineTuningJob(ctx context.Context, req CreateFineTuningJobRequest) (resp CreateFineTuningJobResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCreateFineTuningJob).
		WithContext(ctx).
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

func (c client) RetrieveFineTuningJob(ctx context.Context, req RetrieveFineTuningJobRequest) (resp RetrieveFineTuningJobResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumRetrieveFineTuningJob, "id", req.Body.ID).
		WithContext(ctx).
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

func (c client) ListFineTuningJobs(ctx context.Context, req ListFineTuningJobsRequest) (resp ListFineTuningJobsResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumListFineTuningJobs).
		WithContext(ctx).
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

func (c client) CancelFineTuningJob(ctx context.Context, req CancelFineTuningJobRequest) (resp CancelFineTuningJobResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumCancelFineTuningJob, "id", req.Body.ID).
		WithContext(ctx).
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

func (c client) UploadFile(ctx context.Context, req UploadFileRequest) (resp UploadFileResponseBody, err error) {
	multipart := http.NewMultipartBodyBuilder().WithFile("file", req.FormBody.FileName, req.FormBody.File)
	for k, v := range req.FormBody.ToMultiPartBody() {
		multipart = multipart.WithForm(k, v)
	}
	multipartBody, contentType, buildMultipartErr := multipart.Build()
	if buildMultipartErr != nil {
		return UploadFileResponseBody{}, fmt.Errorf("build upload file multipart body error: %w", buildMultipartErr)
	}

	request := c.options.buildBaseRequest(EndpointEnumUploadFile).
		WithContext(ctx).
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

func (c client) ListFiles(ctx context.Context, req ListFilesRequest) (resp ListFilesResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumListFiles).
		WithContext(ctx).
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

func (c client) DeleteFile(ctx context.Context, req DeleteFileRequest) (resp DeleteFileResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumDeleteFile, "id", req.Body.ID).
		WithContext(ctx).
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

func (c client) RetrieveFile(ctx context.Context, req RetrieveFileRequest) (resp RetrieveFileResponseBody, err error) {
	request := c.options.buildBaseRequest(EndpointEnumRetrieveFile, "id", req.Body.ID).
		WithContext(ctx).
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
		logger:   logger,
	}
}

func NewCustomClient(opts Config, cli http.Client, logger logger.Logger) Client {
	return client{
		executor: cli,
		options:  opts,
		logger:   logger,
	}
}
