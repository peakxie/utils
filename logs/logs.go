package logs

import (
	"strings"

	logf "github.com/0x00b/logrus-formatter"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

const (
	LogFileNameLen = 40
)

// logs init
// maxSize 单位M
// maxBackups 最大备份数
// maxAge 最大过期时间 天
// logLevel 0-6:panic-trace,设置值越大，日志越详细
func Init(fileName string, maxSize, maxBackups, maxAge int, logLevel uint32) {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
	}
	logrus.SetReportCaller(true)
	logrus.SetOutput(lumberjackLogger)
	logrus.SetLevel(logrus.Level(logLevel))

	f := &logf.TextFormatter{TimestampFormat: "01-02 15:04:05.000"}
	f.SetFormatAndTagSource(logf.TagBL, logf.FieldKeyTime, logf.TagBR, logf.FieldKeyLevel, logf.FieldKeyMsg)
	f.NoQuoteFields = true
	f.FormatFileName = formatFileName

	logrus.SetFormatter(f)
}

func formatFileName(name string) string {
	idx := strings.LastIndex(name, "/")
	if -1 != idx {
		idx = strings.LastIndex(name[:idx], "/")
		if -1 != idx {
			name = name[idx:]
		}
	}
	if len(name) > LogFileNameLen {
		return name[len(name)-LogFileNameLen:]
	}
	return name
}
