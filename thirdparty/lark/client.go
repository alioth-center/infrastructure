package lark

import (
	"context"
	"fmt"
	"github.com/alioth-center/infrastructure/trace"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"io"
)

type Client interface {
	initClient(cfg Config)
	UploadFile(ctx context.Context, fileName string, fileType LarkFileType, fileContent io.Reader) (fileKey string, err error)
	UploadMediaFile(ctx context.Context, fileName string, fileType LarkFileType, fileDuration int, fileContent io.Reader) (fileKey string, err error)
	UploadMessageImage(ctx context.Context, imageContent io.Reader) (imageKey string, err error)
	UploadAvatarImage(ctx context.Context, imageContent io.Reader) (imageKey string, err error)
	SendTextMessage(ctx context.Context, receiver LarkReceiver, text string) (err error)
	SendMarkdownMessage(ctx context.Context, receiver LarkReceiver, markdownHeader, markdownContent string, theme LarkMarkdownMessageTheme) (err error)
}

type client struct {
	larkCore *lark.Client
}

func (c *client) initClient(cfg Config) {
	if c.larkCore != nil {
		return
	}

	c.larkCore = lark.NewClient(cfg.AppID, cfg.AppSecret)
}

func (c *client) uploadImage(ctx context.Context, imageType LarkImageType, imageContent io.Reader) (imageKey string, err error) {
	request := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			Image(imageContent).
			ImageType(getLarkImageType(imageType)).
			Build()).
		Build()

	uploadResult, uploadErr := c.larkCore.Im.Image.Create(ctx, request)
	if uploadErr != nil {
		return "", fmt.Errorf("failed to upload image: %w", uploadErr)
	} else if uploadResult.Code != 0 {
		return "", fmt.Errorf("failed to upload image: %s", uploadResult.Msg)
	}

	return *uploadResult.Data.ImageKey, nil
}

func (c *client) buildTextMessage(ctx context.Context, receiver LarkReceiver, text string) (context.Context, *larkim.CreateMessageReq) {
	traceId, newCtx := trace.GetTraceID(ctx)
	message := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(getLarkReceiverIdType(receiver.Type)).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(receiver.Receiver).
			Uuid(traceId).Content(text).
			Build(),
		).
		Build()
	return newCtx, message
}

func (c *client) buildMarkdownMessage(ctx context.Context, receiver LarkReceiver, markdownHeader, markdownContent string, theme LarkMarkdownMessageTheme) (context.Context, *larkim.CreateMessageReq, error) {
	traceId, newCtx := trace.GetTraceID(ctx)

	header := larkcard.NewMessageCardHeader().
		Template(getLarkMarkdownMessageTheme(theme)).
		Title(larkcard.NewMessageCardPlainText().Content(markdownHeader)).
		Build()

	content := larkcard.NewMessageCardMarkdown().
		Content(markdownContent).
		Build()

	card, buildErr := larkcard.NewMessageCard().
		Config(larkcard.NewMessageCardConfig()).
		Header(header).
		Elements([]larkcard.MessageCardElement{content}).
		Build().String()
	if buildErr != nil {
		return newCtx, nil, fmt.Errorf("failed to build markdown message card: %w", buildErr)
	}

	message := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(getLarkReceiverIdType(receiver.Type)).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(receiver.Receiver).
			Uuid(traceId).
			Content(card).
			Build(),
		).Build()

	return newCtx, message, nil
}

func (c *client) UploadFile(ctx context.Context, fileName string, fileType LarkFileType, fileContent io.Reader) (fileKey string, err error) {
	request := larkim.NewCreateFileReqBuilder().
		Body(larkim.NewCreateFileReqBodyBuilder().
			FileType(getLarkFileType(fileType)).
			FileName(fileName).
			File(fileContent).
			Build(),
		).
		Build()

	createResult, createErr := c.larkCore.Im.File.Create(ctx, request)
	if createErr != nil {
		return "", fmt.Errorf("failed to create file: %w", createErr)
	} else if createResult.Code != 0 {
		return "", fmt.Errorf("failed to create file: %s", createResult.Msg)
	}

	return *createResult.Data.FileKey, nil
}

func (c *client) UploadMediaFile(ctx context.Context, fileName string, fileType LarkFileType, fileDuration int, fileContent io.Reader) (fileKey string, err error) {
	request := larkim.NewCreateFileReqBuilder().
		Body(larkim.NewCreateFileReqBodyBuilder().
			FileType(getLarkFileType(fileType)).
			FileName(fileName).
			Duration(fileDuration).
			File(fileContent).
			Build(),
		).
		Build()

	createResult, createErr := c.larkCore.Im.File.Create(ctx, request)
	if createErr != nil {
		return "", fmt.Errorf("failed to create file: %w", createErr)
	} else if createResult.Code != 0 {
		return "", fmt.Errorf("failed to create file: %s", createResult.Msg)
	}

	return *createResult.Data.FileKey, nil
}

func (c *client) UploadMessageImage(ctx context.Context, imageContent io.Reader) (imageKey string, err error) {
	return c.uploadImage(ctx, LarkImageTypeMessage, imageContent)
}

func (c *client) UploadAvatarImage(ctx context.Context, imageContent io.Reader) (imageKey string, err error) {
	return c.uploadImage(ctx, LarkImageTypeAvatar, imageContent)
}

func (c *client) SendTextMessage(ctx context.Context, receiver LarkReceiver, text string) (err error) {
	sendResult, sendErr := c.larkCore.Im.Message.Create(c.buildTextMessage(ctx, receiver, text))
	if sendErr != nil {
		return fmt.Errorf("failed to send text message: %w", sendErr)
	} else if sendResult.Code != 0 {
		return fmt.Errorf("failed to send text message: %s", sendResult.Msg)
	}

	return nil
}

func (c *client) SendMarkdownMessage(ctx context.Context, receiver LarkReceiver, markdownHeader, markdownContent string, theme LarkMarkdownMessageTheme) (err error) {
	sendCtx, sendMsg, buildErr := c.buildMarkdownMessage(ctx, receiver, markdownHeader, markdownContent, theme)
	if buildErr != nil {
		return fmt.Errorf("failed to send markdown message: %w", buildErr)
	}
	sendResult, sendErr := c.larkCore.Im.Message.Create(sendCtx, sendMsg)
	if sendErr != nil {
		return fmt.Errorf("failed to send markdown message: %w", sendErr)
	} else if sendResult.Code != 0 {
		return fmt.Errorf("failed to send markdown message: %s", sendResult.Msg)
	}

	return nil
}

func NewClient(cfg Config) Client {
	client := &client{}
	client.initClient(cfg)
	return client
}
