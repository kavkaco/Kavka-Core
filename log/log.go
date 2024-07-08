package log

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *logger
	logLevel zerolog.Level
)

type logger struct {
	subs   map[string]*SubLogger
	writer io.Writer
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
}

func InitGlobalLogger(config config.Logger) {
	if instance == nil {
		writers := []io.Writer{}

		if slices.Contains(config.Targets, "file") {
			// File writer
			fw := &lumberjack.Logger{
				Filename:   config.Filename,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				Compress:   config.Compress,
			}
			writers = append(writers, fw)
		}

		if slices.Contains(config.Targets, "console") {
			// Console writer.
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		}

		instance = &logger{
			subs:   make(map[string]*SubLogger),
			writer: io.MultiWriter(writers...),
		}

		level, err := zerolog.ParseLevel(strings.ToLower(config.LogLevel))
		if err != nil {
			level = zerolog.InfoLevel // Default to info level if parsing fails.
		}
		zerolog.SetGlobalLevel(level)

		log.Logger = zerolog.New(instance.writer).With().Timestamp().Logger()
	}
}

// NewLoggerLevel initializes the logger level.
func NewLoggerLevel(level zerolog.Level) {
	logLevel = level
}

func getLoggersInst() *logger {
	if instance == nil {
		instance = &logger{
			subs:   make(map[string]*SubLogger),
			writer: zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		}
		log.Logger = zerolog.New(instance.writer).With().Timestamp().Logger()
	}

	return instance
}

// function to set logger level based on env.
func SetLoggerLevel(level string) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel // Default to info level if parsing fails
	}

	logLevel = parsedLevel
}

func GetCurrentLogLevel() zerolog.Level {
	return logLevel
}

func addFields(event *zerolog.Event, keyvals ...interface{}) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[i+1]
		switch v := value.(type) {
		case fmt.Stringer:
			if isNil(v) {
				event.Any(key, v)
			} else {
				event.Stringer(key, v)
			}
		case error:
			event.AnErr(key, v)
		case []byte:
			event.Str(key, hex.EncodeToString(v))
		default:
			event.Any(key, v)
		}
	}

	return event
}

func NewSubLogger(name string) *SubLogger {
	inst := getLoggersInst()
	sl := &SubLogger{
		logger: zerolog.New(inst.writer).With().Timestamp().Logger(),
		name:   name,
	}

	inst.subs[name] = sl

	return sl
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...interface{}) {
	addFields(event, keyvals...).Msg(msg)
}

func (sl *SubLogger) Trace(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...interface{}) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...interface{}) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...interface{}) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...interface{}) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...interface{}) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...interface{}) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...interface{}) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
