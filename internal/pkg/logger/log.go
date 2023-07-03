package logger

// Field is a custom field that the logger preserves in the log.
type Field string

// Fields is a collection of custom fields that are included
// in the log output for a log entry.
type Fields map[Field]interface{}

// Log is the service logger with support of custom fields.
type Log interface {

	// WithFields creates a new logger that is based on this one,
	// but with a new set of custom fields merged into it.
	// In case if the fields with the same keys are already exist
	// in the base log, they are overwritten with a new values.
	WithFields(fields Fields) Log

	// WithField creates a new logger that is based on this one,
	// but with a new custom field merged into it.
	// In case if the specified field already exists
	// in the base log, it is overwritten with a new value.
	WithField(field Field, value interface{}) Log

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	TraceWithFields(fields Fields, args ...interface{})
	TracefWithFields(fields Fields, format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugWithFields(fields Fields, args ...interface{})
	DebugfWithFields(fields Fields, format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoWithFields(fields Fields, args ...interface{})
	InfofWithFields(fields Fields, format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnWithFields(fields Fields, args ...interface{})
	WarnfWithFields(fields Fields, format string, args ...interface{})

	Error(err error, args ...interface{})
	Errorf(err error, format string, args ...interface{})
	ErrorWithFields(err error, fields Fields, args ...interface{})
	ErrorfWithFields(err error, fields Fields, format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	FatalWithFields(fields Fields, args ...interface{})
	FatalfWithFields(fields Fields, format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	PanicWithFields(fields Fields, args ...interface{})
	PanicfWithFields(fields Fields, format string, args ...interface{})
}
