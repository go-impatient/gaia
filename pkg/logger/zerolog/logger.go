package logger

// Source: github.com/xmlking/micro-starter-kit/shared/logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger interface {
	WithLogger(options ...Option) error
	Options() Options
	String() string
}

type defaultLogger struct {
	opts Options
}

func (l *defaultLogger) WithLogger(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}

	// 初始化默认设置
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = nil
	zerolog.LevelFieldName = "level"
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string { return l.String() }

	var newlogger zerolog.Logger

	switch l.opts.Format {
	case JSON:
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		newlogger = zerolog.New(l.opts.Out).
			Level(zerolog.InfoLevel).
			With().Timestamp().Stack().Logger()
	case PRETTY:
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		output := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			if len(l.opts.TimeFormat) > 0 {
				w.TimeFormat = l.opts.TimeFormat
			}
			w.Out = l.opts.Out
			w.NoColor = false
		})
		newlogger = zerolog.New(output).
			Level(zerolog.DebugLevel).
			With().Timestamp().Stack().Logger()
	default:
		return errors.New("unknown log Format string,  defaulting to JSON")
	}

	// Set log Level if not default
	if l.opts.Level != zerolog.NoLevel {
		zerolog.SetGlobalLevel(l.opts.Level)
		newlogger = newlogger.Level(l.opts.Level)
	}

	// Setting timeFormat
	if len(l.opts.TimeFormat) > 0 {
		zerolog.TimeFieldFormat = l.opts.TimeFormat
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		newlogger = newlogger.With().Fields(l.opts.Fields).Logger()
	}

	// Also set it as zerolog's Default logger
	log.Logger = newlogger

	return nil
}

func (l *defaultLogger) Options() Options {
	return l.opts
}

func (l *defaultLogger) String() string {
	return "default"
}

func New(opts ...Option) Logger {
	// 设置默认参数
	options := Options{
		Level:   zerolog.NoLevel,
		Format:  PRETTY,
		Out:     os.Stderr,
		Context: context.Background(),
	}
	l := &defaultLogger{
		opts: options,
	}

	// 初始化
	_ = l.WithLogger(opts...)

	return l
}

func Init(level, format string) {
	var opts []Option
	if len(level) > 0 {
		if lvl, err := zerolog.ParseLevel(level); err != nil {
			log.Fatal().Err(err).Send()
		} else {
			opts = append(opts, WithLevel(lvl))
		}
	}
	if len(format) > 0 {
		if logFmt, err := ParseFormat(format); err != nil {
			log.Fatal().Err(err).Send()
		} else {
			opts = append(opts, WithFormat(logFmt))
		}
	}
	New(opts...)
}
