package openai

import (
	"github.com/alioth-center/infrastructure/network/http"
	"strings"
)

type Config struct {
	ApiKey          string            `json:"api_key" yaml:"api_key" xml:"api_key"`
	BaseUrl         string            `json:"base_url,omitempty" yaml:"base_url,omitempty" xml:"base_url,omitempty"`
	BetaFeatures    string            `json:"beta_features,omitempty" yaml:"beta_features,omitempty" xml:"beta_features,omitempty"`
	CustomEndpoints map[string]string `json:"custom_endpoints,omitempty" yaml:"custom_endpoints,omitempty" xml:"custom_endpoints,omitempty"`
	CustomUserAgent string            `json:"custom_user_agent,omitempty" yaml:"custom_user_agent,omitempty" xml:"custom_user_agent,omitempty"`
}

type EndpointEnum string

func (e EndpointEnum) String() string { return string(e) }

const (
	EndpointEnumListModel             EndpointEnum = "list_models"          // 列出模型
	EndpointEnumRetrieveModel         EndpointEnum = "retrieve_model"       // 检索模型
	EndpointEnumCreateImage           EndpointEnum = "create_image"         // 创建图片
	EndpointEnumCompleteChat          EndpointEnum = "complete_chat"        // 完成聊天
	EndpointEnumCreateSpeech          EndpointEnum = "create_speech"        // 创建语音
	EndpointEnumCreateTranscription   EndpointEnum = "create_transcription" // 创建转录
	EndpointEnumCompleteModeration    EndpointEnum = "complete_moderation"  // 完成审核
	EndpointEnumCreateFineTuningJob   EndpointEnum = "create_fine_tuning"   // 创建微调
	EndpointEnumRetrieveFineTuningJob EndpointEnum = "retrieve_fine_tuning" // 检索微调
	EndpointEnumListFineTuningJobs    EndpointEnum = "list_fine_tuning"     // 列出微调
	EndpointEnumCancelFineTuningJob   EndpointEnum = "cancel_fine_tuning"   // 取消微调
	EndpointEnumUploadFile            EndpointEnum = "upload_file"          // 上传文件
	EndpointEnumListFiles             EndpointEnum = "list_files"           // 列出文件
	EndpointEnumRetrieveFile          EndpointEnum = "retrieve_file"        // 检索文件
	EndpointEnumDeleteFile            EndpointEnum = "delete_file"          // 删除文件
)

const (
	// defaultBaseUrl 默认的请求地址
	defaultBaseUrl = "https://api.openai.com/v1"

	// defaultUserAgent 默认的user-agent
	defaultUserAgent = "alioth-center/http-client v1.2.1"
)

var (

	// defaultEndpoints 默认的endpoint
	defaultEndpoints = map[EndpointEnum]string{
		EndpointEnumListModel:             "models",
		EndpointEnumRetrieveModel:         "models/{:model:}",
		EndpointEnumCreateImage:           "images/generations",
		EndpointEnumCompleteChat:          "chat/completions",
		EndpointEnumCreateSpeech:          "audio/speech",
		EndpointEnumCreateTranscription:   "audio/transcriptions",
		EndpointEnumCompleteModeration:    "moderations",
		EndpointEnumCreateFineTuningJob:   "fine_tuning/jobs",
		EndpointEnumRetrieveFineTuningJob: "fine_tuning/jobs/{:id:}",
		EndpointEnumListFineTuningJobs:    "fine_tuning/jobs",
		EndpointEnumCancelFineTuningJob:   "fine_tuning/jobs/{:id:}/cancel",
		EndpointEnumUploadFile:            "files",
		EndpointEnumListFiles:             "files",
		EndpointEnumRetrieveFile:          "files/{:id:}",
		EndpointEnumDeleteFile:            "files/{:id:}",
	}
)

func (c Config) getRequestUrl(endpoint EndpointEnum) string {
	result := strings.Builder{}

	// 如果config中没有设置base_url，则使用defaultBaseUrl
	if c.BaseUrl != "" {
		result.WriteString(c.BaseUrl)
	} else {
		result.WriteString(defaultBaseUrl)
	}
	result.WriteString("/")

	// 如果config中没有设置custom_endpoints，或者custom_endpoints中没有设置endpoint，则使用defaultEndpoints
	if c.CustomEndpoints != nil {
		if customEndpoint, existEndpoint := c.CustomEndpoints[endpoint.String()]; existEndpoint {
			result.WriteString(customEndpoint)
			return result.String()
		}
	}

	result.WriteString(defaultEndpoints[endpoint])
	return result.String()
}

func (c Config) getUserAgent() string {
	if c.CustomUserAgent != "" {
		return c.CustomUserAgent
	} else {
		return defaultUserAgent
	}
}

func (c Config) buildBaseRequest(endpoint EndpointEnum, injectUrlFunction ...func(original string) (result string)) http.Request {
	baseRequest := http.NewRequest()

	if c.BetaFeatures != "" {
		// 如果设置了beta_features，则注入beta_features
		baseRequest.SetHeader("OpenAI-Beta", c.BetaFeatures)
	}

	if injectUrlFunction == nil || len(injectUrlFunction) == 0 {
		// 如果没有注入函数，则直接使用默认的url
		return baseRequest.
			SetUrl(c.getRequestUrl(endpoint)).
			SetHeader("Authorization", "Bearer "+c.ApiKey)
	} else if injectUrlFunction[0] == nil {
		// 如果注入函数为nil，则直接使用默认的url
		return baseRequest.
			SetUrl(c.getRequestUrl(endpoint)).
			SetHeader("Authorization", "Bearer "+c.ApiKey)
	} else {
		// 如果注入函数不为nil，则使用注入函数注入url
		return baseRequest.
			SetUrl(injectUrlFunction[0](c.getRequestUrl(endpoint))).
			SetHeader("Authorization", "Bearer "+c.ApiKey)
	}
}
