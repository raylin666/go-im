package message

const (
	// TypeText 文本消息类型
	TypeText = "text"
)

// Types 消息类型
var Types = []string{
	TypeText,
}

// HasType 消息类型是否存在
func HasType(msgType string) bool {
	for _, itemType := range Types {
		if itemType == msgType {
			return true
		}
	}

	return false
}

// BuilderType 消息类型构造
type BuilderType interface{}
