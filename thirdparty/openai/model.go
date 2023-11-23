package openai

import (
	"io"
	"strconv"
)

// ModelObject 模型对象
// reference https://platform.openai.com/docs/api-reference/models/object
type ModelObject struct {
	ID      string `json:"id"`              // 模型的唯一标识符
	Created int64  `json:"created"`         // 模型创建的时间
	Object  string `json:"object"`          // 模型的类型，一般为model
	Owner   string `json:"owner,omitempty"` // 模型的所有者，一般为openai
}

// ListModelRequest 列出模型请求
// reference https://platform.openai.com/docs/api-reference/models/list
type ListModelRequest struct {
}

// ListModelResponseBody 列出模型响应
// reference https://platform.openai.com/docs/api-reference/models/list
type ListModelResponseBody struct {
	Object string        `json:"object"` // 一般为list
	Data   []ModelObject `json:"data"`   // 模型列表
}

// RetrieveModelRequest 获取模型请求
// reference https://platform.openai.com/docs/api-reference/models/retrieve
type RetrieveModelRequest struct {
	Model string // 模型的唯一标识符
}

// RetrieveModelResponseBody 获取模型响应
// reference https://platform.openai.com/docs/api-reference/models/retrieve
type RetrieveModelResponseBody struct {
	ModelObject // 模型对象
}

// ImageItem 生成的图片
// reference https://platform.openai.com/docs/api-reference/images/object
type ImageItem struct {
	Url           string `json:"url,omitempty"`            // openai生成的图片url，在response_format为url时返回，默认为url
	Base64        string `json:"b64_json,omitempty"`       // openai生成的图片base64编码，在response_format为b64_json时返回
	RevisedPrompt string `json:"revised_prompt,omitempty"` // openai生成的图片的提示，如果提示被修改了，会返回修改后的
}

// CreateImageRequestBody 生成图片请求体
// reference https://platform.openai.com/docs/api-reference/images/create
type CreateImageRequestBody struct {
	Prompt         string `json:"prompt"`                    // 生成图片的提示
	Model          string `json:"model,omitempty"`           // 生成图片的模型
	N              int    `json:"n"`                         // 生成图片的数量，1~10之间
	Size           string `json:"size"`                      // 生成图片的尺寸，dall-e-2模型只支持256x256,512x512,1024x1024，dall-e-3模型只支持1024x1024,1792x1024,1024x1792
	Quality        string `json:"quality,omitempty"`         // 质量，仅支持dall-e-3模型，标记为hd时生成的图片质量更高
	Style          string `json:"style,omitempty"`           // 风格，仅支持dall-e-3模型，可以选择natural和vivid
	ResponseFormat string `json:"response_format,omitempty"` // 返回格式，url或者b64_json
	User           string `json:"user,omitempty"`            // 用户的唯一标识符，用于openai跟踪
}

// CreateImageRequest 生成图片请求
// reference https://platform.openai.com/docs/api-reference/images/create
type CreateImageRequest struct {
	Body CreateImageRequestBody
}

// ImageResponseBody 生成图片响应
// reference https://platform.openai.com/docs/api-reference/images/create
type ImageResponseBody struct {
	Created int64       `json:"created"` // openai返回的消息回复时间
	Data    []ImageItem `json:"data"`    // openai生成的图片，一般情况下长度和请求的N一致
}

// ChatRoleEnum 聊天角色枚举
// reference https://platform.openai.com/docs/guides/text-generation/chat-completions-api
type ChatRoleEnum string

func (cre ChatRoleEnum) String() string { return string(cre) }

func getChatRoleEnum(enum ChatRoleEnum) string {
	if _, exist := supportedChatRoleEnum[enum.String()]; !exist {
		return ChatRoleEnumUser.String()
	} else {
		return enum.String()
	}
}

// 聊天角色枚举值
const (
	ChatRoleEnumSystem    ChatRoleEnum = "system"    // 系统，用于给模型进行提示的角色，一般用于给出prompt
	ChatRoleEnumAssistant ChatRoleEnum = "assistant" // 助手，模型一般使用该角色
	ChatRoleEnumUser      ChatRoleEnum = "user"      // 用户，用户一般使用该角色
)

var (
	// 支持的聊天角色枚举
	supportedChatRoleEnum = map[string]ChatRoleEnum{
		ChatRoleEnumSystem.String():    ChatRoleEnumSystem,
		ChatRoleEnumAssistant.String(): ChatRoleEnumAssistant,
		ChatRoleEnumUser.String():      ChatRoleEnumUser,
	}
)

// ChatMessageObject 聊天消息对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type ChatMessageObject struct {
	Role    ChatRoleEnum `json:"role"`
	Message string       `json:"message"`
}

// ReplyChoiceObject 回复选择对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type ReplyChoiceObject struct {
	Index        int               `json:"index"`         // 回复序号
	Message      ChatMessageObject `json:"message"`       // 回复信息
	FinishReason string            `json:"finish_reason"` // 停止回复原因
}

// UsageObject token使用量对象
// reference https://platform.openai.com/docs/api-reference/chat/object
type UsageObject struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提问token消耗
	CompletionTokens int `json:"completion_tokens"` // 回复token消耗
	TotalTokens      int `json:"total_tokens"`      // 全部token消耗
}

// CompleteChatRequestBody 聊天请求体
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatRequestBody struct {
	Model            string              `json:"model"`                       // 生成图片的模型
	Messages         []ChatMessageObject `json:"messages"`                    // 提问信息
	Temperature      float64             `json:"temperature,omitempty"`       // 温度采样，0~2，越高越随机
	TopP             float64             `json:"top_p,omitempty"`             // 核采样，0~1，和temperature只有一个能生效
	N                int                 `json:"n,omitempty"`                 // 需要多少消息回复
	Stream           bool                `json:"stream,omitempty"`            // 流式传输开关
	MaxTokens        int                 `json:"max_tokens,omitempty"`        // 允许消耗的最大令牌数
	PresencePenalty  float64             `json:"presence_penalty,omitempty"`  // 创新惩罚，-2~+2，越高越可能出现新东西
	FrequencyPenalty float64             `json:"frequency_penalty,omitempty"` // 重复惩罚，-2~+2，越高越不可能重复
	User             string              `json:"user,omitempty"`              // 用户的唯一标识符，用于openai跟踪
}

// CompleteChatRequest 聊天请求
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatRequest struct {
	Body CompleteChatRequestBody
}

// CompleteChatResponseBody 聊天响应
// reference https://platform.openai.com/docs/api-reference/chat/create
type CompleteChatResponseBody struct {
	ID      string              `json:"id"`      // openai提供的回复id
	Object  string              `json:"object"`  // openai标记的返回对象，此处固定为chat.completion
	Created int64               `json:"created"` // openai返回的消息回复时间
	Choices []ReplyChoiceObject `json:"choices"` // openai的回复，一般情况下只有一个元素
	Usage   UsageObject         `json:"usage"`   // openai的token使用量
}

type VoiceEnum string

func (ve VoiceEnum) String() string { return string(ve) }

func getVoiceEnum(enum VoiceEnum) string {
	if _, exist := supportedVoiceEnum[enum.String()]; !exist {
		return VoiceEnumAlloy.String()
	} else {
		return enum.String()
	}
}

// 语音枚举值
const (
	VoiceEnumAlloy   VoiceEnum = "alloy"   // alloy音色，成年女性，中性声音
	VoiceEnumEcho    VoiceEnum = "echo"    // echo音色，成年男性，中性声音
	VoiceEnumFable   VoiceEnum = "fable"   // fable音色，成年女性，中性声音
	VoiceEnumOnyx    VoiceEnum = "onyx"    // onyx音色，成年男性，低沉声音
	VoiceEnumNova    VoiceEnum = "nova"    // nova音色，成年女性，年轻声音
	VoiceEnumShimmer VoiceEnum = "shimmer" // shimmer音色，成年女性，中性声音
)

var (
	// 支持的语音枚举
	supportedVoiceEnum = map[string]VoiceEnum{
		VoiceEnumAlloy.String():   VoiceEnumAlloy,
		VoiceEnumEcho.String():    VoiceEnumEcho,
		VoiceEnumFable.String():   VoiceEnumFable,
		VoiceEnumOnyx.String():    VoiceEnumOnyx,
		VoiceEnumNova.String():    VoiceEnumNova,
		VoiceEnumShimmer.String(): VoiceEnumShimmer,
	}
)

// CreateSpeechRequestBody 生成语音请求体
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechRequestBody struct {
	Model          string    `json:"model"`                     // 生成语音的模型，目前支持tts-1和tts-1-hd
	Input          string    `json:"input"`                     // 需要生成语音的文本
	Voice          VoiceEnum `json:"voice"`                     // 生成语音的声音
	ResponseFormat string    `json:"response_format,omitempty"` // 返回格式，支持mp3,aac,flac,opus
	Speed          float64   `json:"speed,omitempty"`           // 语速，支持0.25~4.0，1.0为正常语速
}

// CreateSpeechRequest 生成语音请求
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechRequest struct {
	Body CreateSpeechRequestBody
}

// CreateSpeechResponseBody 生成语音响应
// reference https://platform.openai.com/docs/api-reference/audio/createSpeech
type CreateSpeechResponseBody []byte // openai生成的语音，根据response_format返回不同的格式

// CreateTranscriptionRequestBody 生成转录请求体
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionRequestBody struct {
	File           io.Reader // 需要转录的音频文件
	FileName       string    // 音频文件的名称
	Model          string    // 进行音频转写的模型
	Language       string    // 音频文件的语言，ISO-639-1标准
	Prompt         string    // 音频文件的提示
	ResponseFormat string    // 返回格式，支持json和txt
	Temperature    float64   // 温度采样，0~1，越高越随机
}

func (ctr CreateTranscriptionRequestBody) ToMultiPartBody() map[string]string {
	result := map[string]string{
		"model": ctr.Model,
	}
	if ctr.Language != "" {
		result["language"] = ctr.Language
	}
	if ctr.Prompt != "" {
		result["prompt"] = ctr.Prompt
	}
	if ctr.ResponseFormat != "" {
		result["response_format"] = ctr.ResponseFormat
	}
	if ctr.Temperature != 0 {
		result["temperature"] = strconv.FormatFloat(ctr.Temperature, 'f', -1, 64)
	}
	return result
}

// CreateTranscriptionRequest 生成转录请求
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionRequest struct {
	FormBody CreateTranscriptionRequestBody
}

// CreateTranscriptionResponseBody 生成转录响应
// reference https://platform.openai.com/docs/api-reference/audio/createTranscription
type CreateTranscriptionResponseBody struct {
	Text string `json:"text"` // openai生成的文本
}

// ModerationCategoryObject 内容分类对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationCategoryObject struct {
	Sexual                bool `json:"sexual"`                 // 性内容
	Hate                  bool `json:"hate"`                   // 仇恨内容
	Harassment            bool `json:"harassment"`             // 骚扰内容
	SelfHarm              bool `json:"self-harm"`              // 自残内容
	SexualMinors          bool `json:"sexual/minors"`          // 未成年人性内容
	HateThreatening       bool `json:"hate/threatening"`       // 仇恨威胁内容
	ViolenceGraphic       bool `json:"violence/graphic"`       // 暴力内容
	SelfHarmIntent        bool `json:"self-harm/intent"`       // 自残意图内容
	SelfHarmInstr         bool `json:"self-harm/instructions"` // 自残指导内容
	HarassmentThreatening bool `json:"harassment/threatening"` // 骚扰威胁内容
	Violence              bool `json:"violence"`               // 暴力内容
}

// ModerationCategoryScoreObject 内容分类得分对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationCategoryScoreObject struct {
	Sexual                float64 `json:"sexual"`                 // 性内容
	Hate                  float64 `json:"hate"`                   // 仇恨内容
	Harassment            float64 `json:"harassment"`             // 骚扰内容
	SelfHarm              float64 `json:"self-harm"`              // 自残内容
	SexualMinors          float64 `json:"sexual/minors"`          // 未成年人性内容
	HateThreatening       float64 `json:"hate/threatening"`       // 仇恨威胁内容
	ViolenceGraphic       float64 `json:"violence/graphic"`       // 暴力内容
	SelfHarmIntent        float64 `json:"self-harm/intent"`       // 自残意图内容
	SelfHarmInstr         float64 `json:"self-harm/instructions"` // 自残指导内容
	HarassmentThreatening float64 `json:"harassment/threatening"` // 骚扰威胁内容
	Violence              float64 `json:"violence"`               // 暴力内容
}

// ModerationResultObject 内容审核结果对象
// reference: https://platform.openai.com/docs/api-reference/moderations/object
type ModerationResultObject struct {
	Flagged        bool                          `json:"flagged"`         // 是否被标记
	Categories     ModerationCategoryObject      `json:"categories"`      // 内容分类
	CategoryScores ModerationCategoryScoreObject `json:"category_scores"` // 内容分类得分
}

// CompleteModerationRequestBody 内容审核请求体
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationRequestBody struct {
	Model string `json:"model"` // 进行内容审核的模型
	Input string `json:"input"` // 需要审核的文本
}

// CompleteModerationRequest 内容审核请求
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationRequest struct {
	Body CompleteModerationRequestBody
}

// CompleteModerationResponseBody 内容审核响应
// reference https://platform.openai.com/docs/api-reference/moderations/create
type CompleteModerationResponseBody struct {
	ID      string                   `json:"id"`      // openai提供的回复id
	Model   string                   `json:"model"`   // 进行内容审核的模型
	Results []ModerationResultObject `json:"results"` // 内容审核结果
}
