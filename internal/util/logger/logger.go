package logger

/*
Copy paste from https://github.com/learning-cloud-native-go/myapp/blob/step-6/util/logger/logger.go
*/

import (
	"context"
	"io"
	stdlog "log"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"go-react-app/internal/config"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(isDebug bool, conf *config.LogConfig) *Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logLevel, err := zerolog.ParseLevel(conf.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
		if isDebug {
			logLevel = zerolog.DebugLevel
		}
	}

	zerolog.SetGlobalLevel(logLevel)

	err = os.MkdirAll(conf.FolderPath, os.ModePerm)
	if err != nil {
		errors.Wrap(err, "cannot create logs folder")
		log.Fatal().Err(err).Stack().Caller().Msg("")
	}

	logFilePath := path.Join(conf.FolderPath, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		errors.Wrap(err, "cannot create log file")
		log.Fatal().Err(err).Stack().Caller().Msg("")
	}

	writer := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, file)

	logger := zerolog.New(writer).With().Timestamp().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	return &Logger{logger: &logger}
}

// Output duplicates the global logger and sets w as its output.
func (l *Logger) Output(w io.Writer) zerolog.Logger {
	if w == nil {
		return *l.logger
	}
	return l.logger.Output(w)
}

// With creates a child logger with the field added to its context.
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

// Sample returns a logger with the s sampler.
func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func (l *Logger) Ctx(ctx context.Context) *Logger {
	return &Logger{logger: zerolog.Ctx(ctx)}
}
