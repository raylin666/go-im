package message

var _ BuilderType = (*BuilderTextType)(nil)

// BuilderTextType 构建文本消息类型响应
type BuilderTextType struct {
	Text string `json:"text"`
}
