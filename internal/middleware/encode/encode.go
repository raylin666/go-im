package encode

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/transport/http"
	pgkCodec "mt/pkg/codec"
	"mt/pkg/utils"
)

// ResponseEncoder 请求响应封装, 统一处理请求响应序列化
func ResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// 通过 Request Header 的 Accept 中提取出对应的编码器
	// 如果找不到则忽略报错，并使用默认 json编码器
	codec, _ := CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}

	// 在 Response Header 中写入编码器的 scheme
	w.Header().Set("Content-Type", utils.HttpContentType(codec.Name()))
	w.Write(data)
	return nil
}

// ErrorEncoder 错误响应封装, 统一处理错误响应序列化
func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	// 通过 Request Header 的 Accept 中提取出对应的编码器
	// 如果找不到则忽略报错，并使用默认 json编码器
	codec, _ := CodecForRequest(r, "Accept")
	data, _ := codec.Marshal(err)

	// 在 Response Header 中写入编码器的 scheme
	w.Header().Set("Content-Type", utils.HttpContentType(codec.Name()))
	w.Write(data)
}

// CodecForRequest get encoding.Codec via http.Request
func CodecForRequest(r *http.Request, name string) (encoding.Codec, bool) {
	for _, accept := range r.Header[name] {
		codec := encoding.GetCodec(utils.HttpContentSubtype(accept))
		if codec != nil {
			return codec, true
		}
	}
	return encoding.GetCodec(pgkCodec.JSON), false
}
