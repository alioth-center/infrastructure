package lark

type LarkReceiver struct {
	Type     LarkReceiverIdType `json:"type" yaml:"type" xml:"type"`
	Receiver string             `json:"receiver" yaml:"receiver" xml:"receiver"`
}
