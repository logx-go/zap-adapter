package zapadapter

import (
	"fmt"

	"github.com/logx-go/commons/pkg/commons"
	"github.com/logx-go/contract/pkg/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ logx.Logger = (*ZapAdapter)(nil)
var _ logx.Adapter = (*ZapAdapter)(nil)

// New returns a pointer to a new instance of ZapAdapter
func New(logger *zap.Logger) logx.Adapter {
	return &ZapAdapter{
		logger:    logger,
		fields:    make(map[string]any),
		formatter: nil,
		logLevelMap: map[int]zapcore.Level{
			logx.LogLevelDebug:   zapcore.DebugLevel,
			logx.LogLevelInfo:    zapcore.InfoLevel,
			logx.LogLevelNotice:  zapcore.InfoLevel,
			logx.LogLevelWarning: zapcore.WarnLevel,
			logx.LogLevelError:   zapcore.ErrorLevel,
			logx.LogLevelFatal:   zapcore.FatalLevel,
			logx.LogLevelPanic:   zapcore.PanicLevel,
		},
		logLevelDefault: logx.LogLevelInfo,
	}
}

// ZapAdapter implementation to wrap a format Logger
type ZapAdapter struct {
	logger          *zap.Logger
	formatter       logx.Formatter
	fields          map[string]any
	logLevelMap     map[int]zapcore.Level
	logLevelDefault int
}

func (z *ZapAdapter) clone() *ZapAdapter {
	return &ZapAdapter{
		logger:          z.logger,
		fields:          commons.CloneFieldMap(z.fields),
		formatter:       z.formatter,
		logLevelMap:     z.logLevelMap,
		logLevelDefault: z.logLevelDefault,
	}
}

func (z *ZapAdapter) format(v ...any) string {
	if len(v) < 1 {
		if z.formatter == nil {
			return ""
		}

		return z.formatter.Format("", z.fields)
	}

	msg := fmt.Sprintf(`%v`, v[0])
	fields := z.fields

	for i := 1; i < len(v); i += 2 {
		fieldName := ""
		if i < len(v) {
			fieldName = fmt.Sprintf(`%v`, v[i])
		}

		var fieldValue any
		if i+1 < len(v) {
			fieldValue = v[i+1]
		}

		fields[fieldName] = fieldValue
	}

	if z.formatter == nil {
		return msg
	}

	return z.formatter.Format(msg, fields)
}

func (z *ZapAdapter) convertZapFields(fields map[string]any) []zapcore.Field {
	fieldsZ := make([]zapcore.Field, 0)

	for name, value := range fields {
		if name == logx.FieldNameLogLevel {
			continue
		}

		fieldsZ = append(fieldsZ, zap.Reflect(name, value))
	}

	return fieldsZ
}

func (z *ZapAdapter) convertZapLevel(fields map[string]any) zapcore.Level {
	lvl := commons.GetFieldAsIntOrElse(logx.FieldNameLogLevel, fields, z.logLevelDefault)

	if s, ok := z.logLevelMap[lvl]; ok {
		return s
	}

	return zapcore.InfoLevel
}

func (z *ZapAdapter) Fatal(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Fatal(c.format(v...), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Panic(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Panic(c.format(v...), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Print(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Log(c.convertZapLevel(c.fields), c.format(v...), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Fatalf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Fatal(c.format(fmt.Sprintf(format, v...)), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Panicf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Panic(c.format(fmt.Sprintf(format, v...)), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Printf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)

	c.logger.Log(c.convertZapLevel(c.fields), c.format(fmt.Sprintf(format, v...)), c.convertZapFields(c.fields)...)
}

func (z *ZapAdapter) Debug(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelDebug).Print(v...)
}

func (z *ZapAdapter) Info(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelInfo).Print(v...)
}

func (z *ZapAdapter) Notice(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelNotice).Print(v...)
}

func (z *ZapAdapter) Warning(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelWarning).Print(v...)
}

func (z *ZapAdapter) Error(v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelError).Print(v...)
}

func (z *ZapAdapter) Debugf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelDebug).Printf(format, v...)
}

func (z *ZapAdapter) Infof(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelInfo).Printf(format, v...)
}

func (z *ZapAdapter) Noticef(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelNotice).Printf(format, v...)
}

func (z *ZapAdapter) Warningf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelWarning).Printf(format, v...)
}

func (z *ZapAdapter) Errorf(format string, v ...any) {
	c := z.clone()
	c.fields = commons.SetCallerInfo(1, false, c.fields, logx.FieldNameCallerFunc, logx.FieldNameCallerFile, logx.FieldNameCallerLine)
	c.WithField(logx.FieldNameLogLevel, logx.LogLevelError).Printf(format, v...)
}

func (z *ZapAdapter) WithField(name string, value any) logx.Logger {
	c := z.clone()
	c.fields[name] = value

	return c
}

func (z *ZapAdapter) WithFormatter(formatter logx.Formatter) logx.Adapter {
	c := z.clone()
	c.formatter = formatter

	return c
}
