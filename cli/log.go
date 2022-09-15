package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	logTimestampFormatText   = "2006-01-02T15:04:05.000Z07:00" // like time.RFC3339Nano, but shorter and fixed length
	logTimestampFormatPretty = "01-02|15:04:05.000"
)

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer      io.Writer
	Formatter   logrus.Formatter
	LogLevelMin logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	b, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = hook.Writer.Write(b)

	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.LogLevelMin+1]
}

func setLoggingBasic() error {
	// default logging level
	level, err := logrus.ParseLevel(venom.GetString(fmt.Sprintf("log%vlevel", keyDelimiter)))
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	logrus.SetReportCaller(false)

	return nil
}

var ErrUnknownLogTarget = errors.New("unknown log target defined")
var ErrLogFilePathRequired = errors.New("log file path required")

const logFileMode = 0660

func setLoggingHook() error { // nolint: cyclop
	formatter, err := getFormatter(venom.GetString(fmt.Sprintf("log%vformat", keyDelimiter)))
	if err != nil {
		return err
	}

	var writer io.Writer

	switch venom.GetString(fmt.Sprintf("log%vtarget", keyDelimiter)) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "file":
		filePath := venom.GetString(fmt.Sprintf("log%vfile%vpath", keyDelimiter, keyDelimiter))
		if filePath == "" {
			return fmt.Errorf("%w, but was: %v", ErrLogFilePathRequired, filePath)
		}

		writer, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, logFileMode)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w", ErrUnknownLogTarget)
	}

	switch venom.GetString(fmt.Sprintf("log%vtarget", keyDelimiter)) {
	case "stdout":
		fallthrough
	case "stderr":
		fallthrough
	case "file":
		logrus.AddHook(&WriterHook{
			Writer:      writer,
			Formatter:   formatter,
			LogLevelMin: logrus.GetLevel(),
		})
	default:
		return fmt.Errorf("%w", ErrUnknownLogTarget)
	}

	// Send direct logs to nowhere, everything is handled by the hook now
	logrus.SetOutput(io.Discard)

	return nil
}

func initLog() {
	pflags := rootCmd.PersistentFlags()
	pflags.String(fmt.Sprintf("log%vlevel", keyDelimiter), logrus.WarnLevel.String(), "Define minimum log level, one of: "+strings.Join(logLevels(), ", ")+".")
	pflags.String(fmt.Sprintf("log%vtarget", keyDelimiter), "stderr", "Define log output target, one of: stdout, stderr, file.")
	pflags.String(fmt.Sprintf("log%vformat", keyDelimiter), "pretty", "Define default log format, one of: json, keyvalue, pretty.")
	pflags.String(fmt.Sprintf("log%vfile%vpath", keyDelimiter, keyDelimiter), "", fmt.Sprintf("Path to log file. Only required if --log%starget is set to file.", keyDelimiter))

	// Setup basic logging before any hook is setup, so error gets correctly output if something fails before that
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: logTimestampFormatText,
	})
	logrus.SetOutput(os.Stderr)
	logrus.SetReportCaller(true)
}

func logLevels() []string {
	levels := make([]string, len(logrus.AllLevels))
	for k, v := range logrus.AllLevels {
		levels[k] = v.String()
	}

	return levels
}

var ErrUnknownLogFormat = errors.New("unknown log format")

func getFormatter(format string) (logrus.Formatter, error) {
	switch format {
	case "pretty":
		return &logrus.TextFormatter{
			ForceColors: true,
			// QuoteEmptyFields: true,
			FullTimestamp:   true,
			TimestampFormat: logTimestampFormatPretty,
			// DisableSorting:  true,
			PadLevelText: true,
		}, nil

	case "text":
		fallthrough
	case "keyvalue":
		return &logrus.TextFormatter{
			FullTimestamp: true,
			// QuoteEmptyFields: true,
			TimestampFormat: logTimestampFormatText,
		}, nil

	case "json":
		return &logrus.JSONFormatter{
			PrettyPrint:       false,
			DisableHTMLEscape: true,
			TimestampFormat:   logTimestampFormatText,
		}, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownLogFormat, format)
	}
}
