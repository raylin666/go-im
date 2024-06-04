package response

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// ErrorHandler 处理错误并返回自定义响应的中间件
func ErrorHandler() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			resp, err := handler(ctx, req)
			if err != nil {
				fromError := errors.FromError(err)
				// Handle the error and return a custom response.
				code := fromError.GetCode()
				message := fromError.GetMessage()
				data := map[string]interface{}{
					"code":    code,
					"message": message,
				}

				fmt.Println(data)
				//if tr, ok := transport.FromServerContext(ctx); ok {
				// tr.Write(code)
				//}
				return data, fmt.Errorf("%d: %s", code, message)
			}
			return resp, nil
		}
	}
}
