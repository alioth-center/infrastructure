package smtp

import (
	"bytes"
	"html/template"
	"strings"
)

type Content interface {
	setSender(sender string)
	getReceiver() string
	ExportToMailText() []byte
}

type StaticContent struct {
	sender      string
	receiver    string
	subject     string
	content     string
	contentType string
}

func (c *StaticContent) setSender(sender string) {
	c.sender = sender
}

func (c *StaticContent) getReceiver() string {
	return c.receiver
}

func (c *StaticContent) ExportToMailText() []byte {
	builder := strings.Builder{}
	builder.WriteString("From: ")
	builder.WriteString(c.sender)
	builder.WriteString("\n")
	builder.WriteString("To: ")
	builder.WriteString(c.receiver)
	builder.WriteString("\n")
	builder.WriteString("Subject: ")
	builder.WriteString(c.subject)
	builder.WriteString("\n")
	builder.WriteString("MIME-version: 1.0;\n")
	builder.WriteString("Content-Type: ")
	builder.WriteString(c.contentType)
	builder.WriteString("; charset=\"UTF-8\";\n\n")
	builder.WriteString(c.content)
	return []byte(builder.String())
}

func NewStaticTextContent(receiver string, subject string, content string) Content {
	return &StaticContent{
		receiver:    receiver,
		subject:     subject,
		content:     content,
		contentType: "text",
	}
}

func NewStaticHtmlContent(receiver string, subject string, content string) Content {
	return &StaticContent{
		receiver:    receiver,
		subject:     subject,
		content:     content,
		contentType: "text/html",
	}
}

type RenderableContent struct {
	sender      string
	receiver    string
	subject     string
	template    template.Template
	arguments   map[string]any
	contentType string
}

func (c *RenderableContent) setSender(sender string) {
	c.sender = sender
	if c.arguments == nil {
		c.arguments = map[string]any{
			"meta": map[string]string{
				"sender": c.sender,
			},
		}
	} else if c.arguments["meta"] == nil {
		c.arguments["meta"] = map[string]string{
			"sender": c.sender,
		}
	} else {
		c.arguments["meta"].(map[string]string)["sender"] = sender
	}
}

func (c *RenderableContent) getReceiver() string {
	return c.receiver
}

func (c *RenderableContent) ExportToMailText() []byte {
	buffer := bytes.NewBufferString("")
	if c.template.Execute(buffer, c) != nil {
		return []byte{}
	} else {
		return buffer.Bytes()
	}
}

func NewRenderableTextContent(receiver string, subject string, template template.Template) Content {
	arguments := map[string]any{
		"meta": map[string]string{
			"receiver":    receiver,
			"subject":     subject,
			"contentType": "text",
		},
	}
	return &RenderableContent{
		receiver:    receiver,
		subject:     subject,
		template:    template,
		arguments:   arguments,
		contentType: "text",
	}
}

func NewRenderableHtmlContent(receiver string, subject string, template template.Template) Content {
	arguments := map[string]any{
		"meta": map[string]string{
			"receiver":    receiver,
			"subject":     subject,
			"contentType": "text/html",
		},
	}
	return &RenderableContent{
		receiver:    receiver,
		subject:     subject,
		template:    template,
		arguments:   arguments,
		contentType: "text/html",
	}
}
