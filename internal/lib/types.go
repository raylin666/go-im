package lib

import (
	"context"
	"errors"
	"mt/internal/constant/types"
	"mt/internal/middleware/auth"
	"reflect"
)

// HeaderAppID 获取应用权限认证ID信息
func HeaderAppID(ctx context.Context) (appid types.HeaderAppID, err error) {
	ctxAppID := ctx.Value(auth.AppID)
	if reflect.TypeOf(ctxAppID) != reflect.TypeOf(types.HeaderAppID{}) {
		return appid, errors.New("Inconsistent data types.")
	}

	return ctxAppID.(types.HeaderAppID), nil
}
