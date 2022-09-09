// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type IResponse interface {
	SuccessJson(ctx context.Context, code int, message string, data interface{})
	SuccessJsonDefault(ctx context.Context)
	SuccessCodeDefault() int
	SuccessMessageDefault() string
	View(ctx context.Context, template string, data g.Map) (err error)
}

var localResponse IResponse

func Response() IResponse {
	if localResponse == nil {
		panic("implement not found for interface IResponse, forgot register?")
	}
	return localResponse
}

func RegisterResponse(i IResponse) {
	localResponse = i
}
