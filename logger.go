package log

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// logger is a wrapper around zap.Logger and adds a good few features to it.
// It provides layered logs which could be used by separate packages, and could be turned off or on
// separately. Separate layers could also have independent log levels.
// Whenever you change log level it propagates through its children.
type logger struct {
	prefix     string
	skipCaller int
	encoder    zapcore.Encoder
	z          *zap.Logger
	sz         *zap.SugaredLogger
	lvl        zap.AtomicLevel
}

func New(opts ...Option) *logger {
	encodeBuilder := EncoderBuilder().
		WithTimeKey("ts").
		WithLevelKey("level").
		WithNameKey("name").
		WithCallerKey("caller").
		WithMessageKey("msg")

	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	l := &logger{
		lvl:        zap.NewAtomicLevelAt(cfg.level),
		skipCaller: cfg.skipCaller,
	}

	switch cfg.encoder {
	case "json":
		l.encoder = encodeBuilder.JsonEncoder()
	case "console":
		l.encoder = encodeBuilder.ConsoleEncoder()
	}

	cores := []Core{
		zapcore.NewCore(l.encoder, zapcore.Lock(os.Stdout), l.lvl),
	}

	cores = append(cores, cfg.cores...)

	l.z = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(cfg.skipCaller),
	)

	l.sz = l.z.Sugar()

	return l
}

func newNOP() *logger {
	l := &logger{}
	l.z = zap.NewNop()
	l.sz = l.z.Sugar()

	return l
}

func (l *logger) Sugared() *sugaredLogger {
	return &sugaredLogger{
		l: l,
	}
}

func (l *logger) Sync() error {
	return l.z.Sync()
}

func (l *logger) SetLevel(lvl Level) {
	l.lvl.SetLevel(lvl)
}

func (l *logger) With(name string) Logger {
	return l.WithSkip(name, l.skipCaller)
}

func (l *logger) WithSkip(name string, skipCaller int) Logger {
	return l.with(l.z.Core(), name, skipCaller)
}

func (l *logger) WithCore(core Core) Logger {
	return l.with(
		zapcore.NewTee(
			l.z.Core(), core,
		),
		"",
		l.skipCaller,
	)
}

func (l *logger) with(core zapcore.Core, name string, skip int) Logger {
	prefix := l.prefix
	if name != "" {
		prefix = fmt.Sprintf("%s[%s]", l.prefix, name)
	}
	childLogger := &logger{
		prefix:     prefix,
		skipCaller: l.skipCaller,
		encoder:    l.encoder.Clone(),
		z: zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skip),
		),
		sz: zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skip)).Sugar(),
		lvl: l.lvl,
	}

	return childLogger
}

func (l *logger) addPrefix(in string) (out string) {
	if l.prefix != "" {
		sb := &strings.Builder{}
		sb.WriteString(l.prefix)
		sb.WriteRune(' ')
		sb.WriteString(in)
		out = sb.String()

		return out
	}

	return in
}

func (l *logger) WarnOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Warn(guideTxt, fields...)
	}
}

func (l *logger) ErrorOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error(guideTxt, fields...)
	}
}

func (l *logger) checkLevel(lvl Level) bool {
	if l == nil {
		return false
	}

	// Check the level first to reduce the cost of disabled log calls.
	// Since Panic and higher may exit, we skip the optimization for those levels.
	if lvl < zapcore.DPanicLevel && !l.z.Core().Enabled(lvl) {
		return false
	}

	return true
}

func (l *logger) Check(lvl Level, msg string) *CheckedEntry {
	if !l.checkLevel(lvl) {
		return nil
	}

	return l.z.Check(lvl, l.addPrefix(msg))
}

func (l *logger) Debug(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(DebugLevel) {
		return
	}
	if ce := l.z.Check(DebugLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) Info(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(InfoLevel) {
		return
	}
	if ce := l.z.Check(InfoLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) Warn(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(WarnLevel) {
		return
	}
	if ce := l.z.Check(WarnLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) Error(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(ErrorLevel) {
		return
	}
	if ce := l.z.Check(ErrorLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) Fatal(msg string, fields ...Field) {
	if l == nil {
		return
	}
	l.z.Fatal(l.addPrefix(msg), fields...)
}

func (l *logger) RecoverPanic(funcName string, extraInfo interface{}, compensationFunc func()) {
	if r := recover(); r != nil {
		l.Error("Panic Recovered",
			zap.String("Func", funcName),
			zap.Any("Info", extraInfo),
			zap.Any("Recover", r),
			zap.ByteString("StackTrace", debug.Stack()),
		)
		if compensationFunc != nil {
			go compensationFunc()
		}
	}
}

type sugaredLogger struct {
	l *logger
}

var _ SugaredLogger = (*sugaredLogger)(nil)

func (l sugaredLogger) Debugf(template string, args ...interface{}) {
	l.l.sz.Debugf(l.l.addPrefix(template), args...)
}

func (l sugaredLogger) Infof(template string, args ...interface{}) {
	l.l.sz.Infof(l.l.addPrefix(template), args...)
}

func (l sugaredLogger) Printf(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

func (l sugaredLogger) Warnf(template string, args ...interface{}) {
	l.l.sz.Warnf(l.l.addPrefix(template), args...)
}

func (l sugaredLogger) Errorf(template string, args ...interface{}) {
	l.l.sz.Errorf(l.l.addPrefix(template), args...)
}

func (l sugaredLogger) Fatalf(template string, args ...interface{}) {
	l.l.sz.Fatalf(l.l.addPrefix(template), args...)
}

func (l sugaredLogger) Debug(args ...interface{}) {
	l.l.sz.Debug(args...)
}

func (l sugaredLogger) Info(args ...interface{}) {
	l.l.sz.Info(args...)
}

func (l sugaredLogger) Warn(args ...interface{}) {
	l.l.sz.Warn(args...)
}

func (l sugaredLogger) Error(args ...interface{}) {
	l.l.sz.Error(args...)
}

func (l sugaredLogger) Fatal(args ...interface{}) {
	l.l.sz.Fatal(args...)
}

func (l sugaredLogger) Panic(args ...interface{}) {
	l.l.sz.Panic(args...)
}
