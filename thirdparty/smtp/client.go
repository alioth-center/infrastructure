package smtp

import (
	"crypto/tls"
	"net/smtp"
	"strconv"
)

type Client interface {
	initClient() error
	closeClient() error
	SendContent(content Content) error
}

type tlsClient struct {
	secret    string
	server    string
	port      int
	sender    string
	signature string
	enableTLS bool

	client *smtp.Client
}

func (c *tlsClient) initClient() error {
	if c.client != nil {
		return nil
	}

	serverAddr := c.server + ":" + strconv.Itoa(c.port)
	if conn, dialErr := tls.Dial("tcp", serverAddr, &tls.Config{
		ServerName:         c.server,
		InsecureSkipVerify: c.enableTLS, // nolint:gosec
		MinVersion:         tls.VersionTLS12,
	}); dialErr != nil {
		return NewDialSmtpServerError(dialErr)
	} else if client, dialClientErr := smtp.NewClient(conn, c.server); dialClientErr != nil {
		return NewDialSmtpServerError(dialClientErr)
	} else if authErr := client.Auth(smtp.PlainAuth("", c.sender, c.secret, c.server)); authErr != nil {
		return NewDialSmtpServerError(authErr)
	} else {
		c.client = client
		return nil
	}
}

func (c *tlsClient) closeClient() error {
	if c.client == nil {
		return nil
	}

	if closeErr := c.client.Quit(); closeErr != nil {
		return NewDialSmtpServerError(closeErr)
	} else {
		c.client = nil
		return nil
	}
}

func (c *tlsClient) SendContent(content Content) error {
	// 检查是否已经连接
	if c.client == nil {
		if initClientErr := c.initClient(); initClientErr != nil {
			return initClientErr
		}
	}

	// 发送邮件
	content.setSender(c.signature)
	payload := content.ExportToMailText()
	if mailToErr := c.client.Mail(c.sender); mailToErr != nil {
		return NewInitMailContentError(payload, mailToErr)
	} else if rcptErr := c.client.Rcpt(content.getReceiver()); rcptErr != nil {
		return NewInitMailContentError(payload, rcptErr)
	} else if writer, initWriterErr := c.client.Data(); initWriterErr != nil {
		return NewInitMailContentError(payload, initWriterErr)
	} else if _, writeContentErr := writer.Write(payload); writeContentErr != nil {
		return NewInitMailContentError(payload, writeContentErr)
	} else if closeErr := writer.Close(); closeErr != nil {
		return NewDialSmtpServerError(closeErr)
	} else {
		return nil
	}
}

func NewTLSClient(cfg Config) (client Client, err error) {
	client = &tlsClient{
		secret:    cfg.Secret,
		server:    cfg.Server,
		port:      cfg.Port,
		sender:    cfg.Sender,
		signature: cfg.Signature,
		enableTLS: cfg.EnableTLS,
	}
	if initClientErr := client.initClient(); initClientErr != nil {
		return nil, initClientErr
	} else {
		return client, nil
	}
}
