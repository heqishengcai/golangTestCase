package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	Level int
)

const (
	LevelFatal = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var _fileLogWriter = newFileWriter()

// fileLogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type fileLogWriter struct {
	sync.RWMutex // write log order by order and  atomic incr maxLinesCurLines and maxSizeCurSize
	// The opened file

	logSync     sync.RWMutex
	accessSync  sync.RWMutex
	monitorSync sync.RWMutex

	_logFile *os.File
	//_stdLogFile *os.File
	_accesslogFile  *os.File
	_monitorlogFile *os.File

	_log        *logger
	_accesslog  *logger
	_monitorlog *logger

	LogFile, AccessLogFile, MonitorLogFile string

	// Rotate daily
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int
	dailyOpenTime time.Time

	Rotate bool `json:"rotate"`

	Perm string `json:"perm"`

	RotatePerm string `json:"rotateperm"`

	_logFileNameOnly, _logSuffix string // like "project.log", project is fileNameOnly and .log is suffix
	//_stdLogFileNameOnly, _stdLogSuffix string // like "project.log", project is fileNameOnly and .log is suffix
	_accesslogFileNameOnly, _accesslogSuffix   string // like "project.log", project is fileNameOnly and .log is suffix
	_monitorlogFileNameOnly, _monitorlogSuffix string // like "project.log", project is fileNameOnly and .log is suffix
}

// Logger defines the behavior of a log provider.
type Logger interface {
	Init() error
	Log() *logger
	Accesslog() *logger
	Monitorlog() *logger
}

// newFileWriter create a FileLogWriter returning as LoggerInterface.
func newFileWriter() Logger {
	w := &fileLogWriter{
		Daily:          true,
		MaxDays:        7,
		Rotate:         true,
		RotatePerm:     "0444",
		Perm:           "0664",
		LogFile:        "/data/log/go.log",
		AccessLogFile:  "/data/log/access.log",
		MonitorLogFile: "/data/log/monitor.log",
	}
	return w
}

func (w *fileLogWriter) Monitorlog() *logger {
	w.monitorSync.RLock()
	defer w.monitorSync.RUnlock()
	return w._monitorlog
}

func (w *fileLogWriter) Accesslog() *logger {
	w.accessSync.RLock()
	defer w.accessSync.RUnlock()
	return w._accesslog
}

func (w *fileLogWriter) Log() *logger {
	w.logSync.RLock()
	defer w.logSync.RUnlock()
	return w._log
}

func (w *fileLogWriter) Init() error {
	w._logSuffix = filepath.Ext(w.LogFile)
	w._logFileNameOnly = strings.TrimSuffix(w.LogFile, w._logSuffix)
	if w._logSuffix == "" {
		w._logSuffix = ".log"
	}

	w._accesslogSuffix = filepath.Ext(w.AccessLogFile)
	w._accesslogFileNameOnly = strings.TrimSuffix(w.AccessLogFile, w._accesslogSuffix)
	if w._accesslogSuffix == "" {
		w._accesslogSuffix = ".log"
	}

	w._monitorlogSuffix = filepath.Ext(w.MonitorLogFile)
	w._monitorlogFileNameOnly = strings.TrimSuffix(w.MonitorLogFile, w._monitorlogSuffix)
	if w._monitorlogSuffix == "" {
		w._monitorlogSuffix = ".log"
	}
	err := w.startLogger()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (w *fileLogWriter) startLogger() error {
	err := w.createLogFile()
	if err != nil {
		return err
	}
	err = w.createAccessLogFile()
	if err != nil {
		return err
	}
	err = w.createMonitorLogFile()
	if err != nil {
		return err
	}
	return w.initFd()
}

func (w *fileLogWriter) initFd() error {
	_, err := w._logFile.Stat()
	if err != nil {
		return err
	}

	_, err = w._accesslogFile.Stat()
	if err != nil {
		return err
	}

	_, err = w._monitorlogFile.Stat()
	if err != nil {
		return err
	}

	w.dailyOpenTime = time.Now()
	//w.dailyOpenTime = time.Unix(1534348780, 0)
	w.dailyOpenDate = w.dailyOpenTime.Day()
	if w.Daily {
		go w.dailyRotate(w.dailyOpenTime)
	}

	return nil
}

func (w *fileLogWriter) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotate(time.Now().Day()) {
		if err := w.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter: %s\n", err)
		}
	}
	w.Unlock()
}

func (w *fileLogWriter) cuttingLog(fileName, fileNameOnly, suffix string, rotatePerm int64, file *os.File) error {
	//fmt.Printf("fileName: %s, fileNameOnly: %s, suffix: %s \n", fileName, fileNameOnly, suffix)
	num := 1
	logfName := ""

	_, err := os.Lstat(fileName)
	if err != nil {
		//even if the file is not exist or other ,we should RESTART the logger
		return err
	}

	logfName = fmt.Sprintf("%s%s.%s", fileNameOnly, suffix, w.dailyOpenTime.Format("2006-01-02"))
	_, err = os.Lstat(logfName)
	for ; err == nil && num <= 10000; num++ {
		logfName = fileNameOnly +
			fmt.Sprintf(".%03d%s.%s", num, suffix, w.dailyOpenTime.Format("2006-01-02"))
		_, err = os.Lstat(logfName)
	}

	// return error if the last file checked still existed
	if err == nil {
		return err
	}

	// close fileWriter before rename
	//file.Close()

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	err = os.Rename(fileName, logfName)
	if err != nil {
		return err
	}

	err = os.Chmod(logfName, os.FileMode(rotatePerm))
	if err != nil {
		return err
	}

	return nil
}

// DoRotate means it need to write file in new file.
// new file name like xx.2013-01-01.log (daily) or xx.001.log (by line or size)
func (w *fileLogWriter) doRotate(logTime time.Time) error {
	// file exists
	// Find the next available number

	rotatePerm, err := strconv.ParseInt(w.RotatePerm, 8, 64)
	if err != nil {
		return err
	}

	//logFile
	err1 := w.cuttingLog(w.LogFile, w._logFileNameOnly, w._logSuffix, rotatePerm, w._logFile)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "doRotate err: %s", err1.Error())
		goto RESTART_LOGGER
	}

	//accessLogFile
	err1 = w.cuttingLog(w.AccessLogFile, w._accesslogFileNameOnly, w._accesslogSuffix, rotatePerm, w._accesslogFile)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "doRotate err: %s", err1.Error())
		goto RESTART_LOGGER
	}

	//monitorLogFile
	err1 = w.cuttingLog(w.MonitorLogFile, w._monitorlogFileNameOnly, w._monitorlogSuffix, rotatePerm, w._monitorlogFile)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "doRotate err: %s", err1.Error())
		goto RESTART_LOGGER
	}

RESTART_LOGGER:

	startLoggerErr := w.startLogger()

	if startLoggerErr != nil {
		return startLoggerErr
	}

	return nil
}

func (w *fileLogWriter) needRotate(day int) bool {
	return (w.Daily && day != w.dailyOpenDate)
}

func (w *fileLogWriter) createLogFile() error {
	// Open the log file
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return err
	}

	//正常日志
	logFileFd, err := os.OpenFile(w.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return err
	}

	// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
	os.Chmod(w.LogFile, os.FileMode(perm))

	w.logSync.Lock()
	defer w.logSync.Unlock()

	if w._logFile != nil {
		w._logFile.Close()
	}
	w._logFile = logFileFd

	if w._log != nil {
		w._log._log.SetOutput(logFileFd)
	} else {
		w._log = &logger{_log: log.New(logFileFd, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
	}

	return nil
}

func (w *fileLogWriter) createAccessLogFile() error {
	// Open the log file
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return err
	}

	w.accessSync.Lock()
	defer w.accessSync.Unlock()

	//accessLogFile
	accessLogFileFd, err := os.OpenFile(w.AccessLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return err
	}

	os.Chmod(w.AccessLogFile, os.FileMode(perm))

	if w._accesslogFile != nil {
		w._accesslogFile.Close()
	}
	w._accesslogFile = accessLogFileFd

	if w._accesslog != nil {
		w._accesslog._log.SetOutput(accessLogFileFd)
	} else {
		w._accesslog = &logger{_log: log.New(accessLogFileFd, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
	}

	return nil
}

func (w *fileLogWriter) createMonitorLogFile() error {
	// Open the log file
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return err
	}

	w.monitorSync.Lock()
	defer w.monitorSync.Unlock()

	//monitorLogFile
	monitorLogFileFd, err := os.OpenFile(w.MonitorLogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err != nil {
		return err
	}

	os.Chmod(w.MonitorLogFile, os.FileMode(perm))

	if w._monitorlogFile != nil {
		w._monitorlogFile.Close()
	}
	w._monitorlogFile = monitorLogFileFd

	if w._monitorlog != nil {
		w._monitorlog._log.SetOutput(monitorLogFileFd)
	} else {
		w._monitorlog = &logger{_log: log.New(monitorLogFileFd, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
	}

	return nil
}

func init() {
	err := _fileLogWriter.Init()
	fmt.Printf("_fileLogWriter err: %v\n", err)
}

//
//******************monitor.log文件
//日志格式举例：
//[error][20180529 13:50:55] [122] [log/a.php:27][read_ihead:56 read_body:123 total:400][host(12.33.33.33:3600) word (aa cc)] Read mysql failed
//整体日志列分割采用[]进行分割，如果字段为空在[]不输出内容
//括号中内容为与falcon对应关系，在做指标监控时
//1、日志级别：可以分debug info warn error fatal，如果作为指标监控时，请采用info字段，对应falcon的json字段在下方括号内对应关系标注
//2、时间：必须按例子中规定格式输出内容 （timestamp）
//3、指标或错误：主要用于运维根据错误号来配置报警等级、接收人等，必须按要求输出且需要在wiki记录，错误与指标只能用小写英语单词加下划线或者纯数字，该字段占位最长10个字符（metric）
//4、程序文件名与行号：可为空
//5、耗时：超时日志必须打印 总耗时：毫秒。其他两个字段用于自定义扩展，为各阶段毫秒，可为空（value）
//6、参数：可为空 （tags）
//7、描述：可为空
//http://wiki.corp.ttyongche.com:8360/confluence/pages/viewpage.action?pageId=20120648
//
//New 实例化
const (
	//日志级别
	MONITOR_FATAL = "fatal"
	MONITOR_ERROR = "error"
	MONITOR_WARN  = "warn"
	MONITOR_INFO  = "info"
	MONITOR_DEBUG = "debug"

	//错误code
	MONITOR_ERROR_INTERFACE_TIMEOUT        = 101 //调用接口超时
	MONITOR_ERROR_INTERFACE_EXCEPTION      = 102 //接口返回数据异常
	MONITOR_ERROR_INTERFACE_EXCEED_MAXTIME = 103 //接口超过最大时间
)

//MONITOR_ERROR_INTERFACE_TIMEOUT   = 101    //调用接口超时
func MonitorErrorfInterTimeout(consumingTime uint64, format string, v ...interface{}) {
	MonitorErrorf(MONITOR_ERROR_INTERFACE_TIMEOUT, consumingTime, format, v...)
}

//MONITOR_ERROR_INTERFACE_EXCEPTION = 102    //接口返回数据异常
func MonitorErrorfInterException(consumingTime uint64, format string, v ...interface{}) {
	MonitorErrorf(MONITOR_ERROR_INTERFACE_EXCEPTION, consumingTime, format, v...)
}

//MONITOR_ERROR_INTERFACE_EXCEPTION = 102    //接口超过最大时间
func MonitorErrorfInterExceedMaxTime(consumingTime uint64, format string, v ...interface{}) {
	MonitorErrorf(MONITOR_ERROR_INTERFACE_EXCEED_MAXTIME, consumingTime, format, v...)
}

//[error][20180529 13:50:55] [122] [log/a.php:27][read_ihead:56 read_body:123 total:400][host(12.33.33.33:3600) word (aa cc)] Read mysql failed
func MonitorErrorf(reasonCode uint32, consumingTime uint64, format string, v ...interface{}) {
	monitorLog(MONITOR_ERROR, reasonCode, consumingTime, fmt.Sprintf(format, v...))
}

func MonitorInfo(reasonCode uint32, consumingTime uint64, format string, v ...interface{}) {
	monitorLog(MONITOR_INFO, reasonCode, consumingTime, fmt.Sprintf(format, v...))
}

func MonitorFatal(reasonCode uint32, consumingTime uint64, format string, v ...interface{}) {
	monitorLog(MONITOR_FATAL, reasonCode, consumingTime, fmt.Sprintf(format, v...))
}

func MonitorWarn(reasonCode uint32, consumingTime uint64, format string, v ...interface{}) {
	monitorLog(MONITOR_WARN, reasonCode, consumingTime, fmt.Sprintf(format, v...))
}

func MonitorDebug(reasonCode uint32, consumingTime uint64, format string, v ...interface{}) {
	monitorLog(MONITOR_DEBUG, reasonCode, consumingTime, fmt.Sprintf(format, v...))
}

//监控日志通用结构
func monitorLog(logLevel string, reasonCode uint32, consumingTime uint64, desc string) {
	//获取代码文件名称与代码行数
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	err := _fileLogWriter.Monitorlog().Output(0, fmt.Sprintf("[%s][%s] [%d] [%s/%d] [read_ihead:0 read_body:0 total:%d] [] %s",
		logLevel,
		time.Now().Format("20060102 15:04:05"),
		reasonCode,
		file,
		line,
		consumingTime,
		desc,
	))

	if err != nil {
		fmt.Fprintf(os.Stderr, "write monitor log to file fail, err: %s", err.Error())
		panic("monitorLog")
	}
}

//
//******************access.log文件
//
func AccessInfof(format string, v ...interface{}) {
	err := _fileLogWriter.Accesslog().Output(LevelInfo, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write access info log to file fail, err: %s", err.Error())
	}
}

func AccessErrorf(format string, v ...interface{}) {
	err := _fileLogWriter.Accesslog().Output(LevelError, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write access error log to file fail, err: %s", err.Error())
	}
}

//
//*******************_log文件日志
//
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Fatal(s string) {
	err := _fileLogWriter.Log().Output(LevelFatal, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Fatal log to file fail, err: %s", err.Error())
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	err := _fileLogWriter.Log().Output(LevelFatal, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Fatalf log to file fail, err: %s", err.Error())
	}
	os.Exit(1)
}

func Error(s string) {
	err := _fileLogWriter.Log().Output(LevelError, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Error log to file fail, err: %s", err.Error())
	}
}

func Errorf(format string, v ...interface{}) {
	err := _fileLogWriter.Log().Output(LevelError, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Errorf log to file fail, err: %s", err.Error())
	}
}

func Warn(s string) {
	err := _fileLogWriter.Log().Output(LevelWarning, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Warn log to file fail, err: %s", err.Error())
	}
}

func Warnf(format string, v ...interface{}) {
	err := _fileLogWriter.Log().Output(LevelWarning, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Warnf log to file fail, err: %s", err.Error())
	}
}

func Info(s string) {
	err := _fileLogWriter.Log().Output(LevelInfo, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Info log to file fail, err: %s", err.Error())
	}
}

func Infof(format string, v ...interface{}) {
	err := _fileLogWriter.Log().Output(LevelInfo, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Infof log to file fail, err: %s", err.Error())
	}
}

func Debug(s string) {
	err := _fileLogWriter.Log().Output(LevelDebug, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Debug log to file fail, err: %s", err.Error())
	}
}

func Debugf(format string, v ...interface{}) {
	err := _fileLogWriter.Log().Output(LevelDebug, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Debugf log to file fail, err: %s", err.Error())
	}
}

func SetLogLevel(level Level) {
	_fileLogWriter.Log().SetLogLevel(level)
}

type logger struct {
	_log *log.Logger
	//小于等于该级别的level才会被记录
	logLevel Level
}

func (l *logger) Output(level Level, s string) error {
	if l.logLevel < level {
		return nil
	}
	formatStr := "[UNKNOWN] %s"
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m %s"
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m %s"
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m %s"
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m %s"
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m %s"
	}
	s = fmt.Sprintf(formatStr, s)
	return l._log.Output(3, s)
}

func (l *logger) Fatal(s string) {
	err := l.Output(LevelFatal, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Fatal log to file fail, err: %s", err.Error())
	}
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	err := l.Output(LevelFatal, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Fatalf log to file fail, err: %s", err.Error())
	}
	os.Exit(1)
}

func (l *logger) Error(s string) {
	err := l.Output(LevelError, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Error log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	err := l.Output(LevelError, fmt.Sprintf(format, v...))
	if err != nil {
		panic("Errorf")
		fmt.Fprintf(os.Stderr, "write Errorf log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Warn(s string) {
	err := l.Output(LevelWarning, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Warn log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	err := l.Output(LevelWarning, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Warnf log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Info(s string) {
	err := l.Output(LevelInfo, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Info log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	err := l.Output(LevelInfo, fmt.Sprintf(format, v...))
	if err != nil {
		panic("Errorf")
		fmt.Fprintf(os.Stderr, "write Infof log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Debug(s string) {
	err := l.Output(LevelDebug, s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Debug log to file fail, err: %s", err.Error())
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	err := l.Output(LevelDebug, fmt.Sprintf(format, v...))
	if err != nil {
		fmt.Fprintf(os.Stderr, "write Debugf log to file fail, err: %s", err.Error())
	}
}

func (l *logger) SetLogLevel(level Level) {
	l.logLevel = level
}
