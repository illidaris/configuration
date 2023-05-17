package configuration

import (
	"context"
	"fmt"
)

var _ = ILogger(&DefaultLogger{})
var defaultLogger ILogger

type ILogger interface {
	DebugCtx(ctx context.Context, msg string)
	InfoCtx(ctx context.Context, msg string)
	WarnCtx(ctx context.Context, msg string)
	ErrorCtx(ctx context.Context, msg string)
}

type DefaultLogger struct {
}

func (l *DefaultLogger) DebugCtx(_ context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) InfoCtx(_ context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) WarnCtx(_ context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) ErrorCtx(_ context.Context, msg string) {
	fmt.Println(msg)
}
