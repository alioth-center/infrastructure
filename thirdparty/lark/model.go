package lark

// Receiver 飞书消息接收方
type Receiver struct {
	Type     LarkReceiverIdType `json:"type" yaml:"type" xml:"type"`             // 接收方类型
	Receiver string             `json:"receiver" yaml:"receiver" xml:"receiver"` // 接收方ID，根据type不同而不同
}
