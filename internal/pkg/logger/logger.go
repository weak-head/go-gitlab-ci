package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

const (

	// FieldNode is a unique node name this service is running on
	FieldNode = "node"

	// FieldService is a unique name of this service
	FieldService = "service"

	// FieldPackage is a unique name of the package
	FieldPackage = "package"

	// FieldFunction is a unique name of the function in a package
	FieldFunction = "function"

	// FieldError is a unique error id
	FieldError = "error"

	// FieldCorrelation is a unique correlation id
	FieldCorrelation = "correlationId"
)

// Config is a logger configuration
type Config struct {

	// Supported log levels:
	// - trace
	// - debug
	// - info
	// - warning
	// - error
	// - fatal
	// - panic
	Level string

	// Supported log formatters:
	// - text
	// - json
	Formatter string
}

// log implements service logger with tracing and custom field support.
type log struct {
	config Config
	logger *logrus.Logger
	fields logrus.Fields
}

// New creates a new service logger.
func New(config Config) (Log, error) {

	// Logrus is used as the underlying logger
	l := &log{
		logger: logrus.New(),
		fields: logrus.Fields{},
	}

	// Parse and apply the configuration
	if err := l.applyConfig(config); err != nil {
		return nil, err
	}

	return l, nil
}

// NewNullLogger creates a new discarding null logger for unit test purposes.
func NewNullLogger() (Log, *logtest.Hook) {
	logger, hook := logtest.NewNullLogger()
	return &log{
		logger: logger,
	}, hook
}

// WithFields combines this logger with a new custom fields.
func (l *log) WithFields(fields Fields) Log {
	return &log{
		logger: l.logger,
		fields: l.combineFields(fields),
	}
}

// WithField combines this logger with a new custom fields.
func (l *log) WithField(field Field, value interface{}) Log {
	return l.WithFields(Fields{field: value})
}

func (l *log) Trace(args ...interface{}) {
	l.logger.WithFields(l.fields).Trace(args...)
}

func (l *log) Tracef(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Tracef(format, args...)
}

func (l *log) TraceWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Trace(args...)
}

func (l *log) TracefWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Tracef(format, args...)
}

func (l *log) Debug(args ...interface{}) {
	l.logger.WithFields(l.fields).Debug(args...)
}

func (l *log) Debugf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Debugf(format, args...)
}

func (l *log) DebugWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Debug(args...)
}

func (l *log) DebugfWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Debugf(format, args...)
}

func (l *log) Info(args ...interface{}) {
	l.logger.WithFields(l.fields).Info(args...)
}

func (l *log) Infof(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Infof(format, args...)
}

func (l *log) InfoWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Info(args...)
}

func (l *log) InfofWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Infof(format, args...)
}

func (l *log) Warn(args ...interface{}) {
	l.logger.WithFields(l.fields).Warn(args...)
}

func (l *log) Warnf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Warnf(format, args...)
}

func (l *log) WarnWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Warn(args...)
}

func (l *log) WarnfWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Warnf(format, args...)
}

func (l *log) Error(err error, args ...interface{}) {
	combined := l.combineFields(Fields{FieldError: err})
	l.logger.WithFields(combined).Error(args...)
}

func (l *log) Errorf(err error, format string, args ...interface{}) {
	combined := l.combineFields(Fields{FieldError: err})
	l.logger.WithFields(combined).Errorf(format, args...)
}

func (l *log) ErrorWithFields(err error, fields Fields, args ...interface{}) {
	combined := l.combineFields(fields)
	combined[FieldError] = err
	l.logger.WithFields(combined).Error(args...)
}

func (l *log) ErrorfWithFields(err error, fields Fields, format string, args ...interface{}) {
	combined := l.combineFields(fields)
	combined[FieldError] = err
	l.logger.WithFields(combined).Errorf(format, args...)
}

func (l *log) Fatal(args ...interface{}) {
	l.logger.WithFields(l.fields).Fatal(args...)
}

func (l *log) Fatalf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Fatalf(format, args...)
}

func (l *log) FatalWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Fatal(args...)
}

func (l *log) FatalfWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Fatalf(format, args...)
}

func (l *log) Panic(args ...interface{}) {
	l.logger.WithFields(l.fields).Panic(args...)
}

func (l *log) Panicf(format string, args ...interface{}) {
	l.logger.WithFields(l.fields).Panicf(format, args...)
}

func (l *log) PanicWithFields(fields Fields, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Panic(args...)
}

func (l *log) PanicfWithFields(fields Fields, format string, args ...interface{}) {
	l.logger.WithFields(l.combineFields(fields)).Panicf(format, args...)
}

// applyConfig
func (l *log) applyConfig(config Config) error {
	l.config = config

	// log formatter
	formatter, err := getFormatter(config.Formatter)
	if err != nil {
		return err
	}
	l.logger.SetFormatter(formatter)

	// log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return err
	}
	l.logger.SetLevel(level)

	// log output
	l.logger.SetOutput(os.Stdout)

	return nil
}

// combineFields combines existing logger fields with a new ones.
func (l *log) combineFields(fields Fields) logrus.Fields {
	combined := logrus.Fields{}

	// preserve original fields
	for k, v := range l.fields {
		combined[k] = v
	}

	// add / overwrite with the new fields
	for k, v := range fields {
		combined[string(k)] = v
	}

	return combined
}

// getFormatter parses and creates a new log formatter.
func getFormatter(formatter string) (logrus.Formatter, error) {
	switch formatter {
	case "text":
		return &logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		}, nil

	case "json":
		return &logrus.JSONFormatter{}, nil

	default:
		return &logrus.JSONFormatter{}, nil
	}
}
