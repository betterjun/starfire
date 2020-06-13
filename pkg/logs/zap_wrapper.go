package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/exec"
)

// 方便配置，用来初始化日志级别
var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// 包级别的日志
var packageLogger *zap.Logger

// 下面定义包级别的函数，方便用包名调用。

// 初始化包级别的日志，必须在程序启动时调用
func Init(filename, level string, maxSizeMB, maxAgeDay, maxBackups int, dev bool) *zap.Logger {
	if dev {
		packageLogger = NewDevLogger(filename, level, maxSizeMB, maxAgeDay, maxBackups)
	} else {
		packageLogger = NewProdLogger(filename, level, maxSizeMB, maxAgeDay, maxBackups)
	}
	return packageLogger
}

// 获取命名日志，s可为空，表示包级别日志;不为空，则为子级别日志
func Named(s string) *zap.Logger {
	return packageLogger.Named(s)
}

// 封装下面的接口，是为了方便，直接用package名来调用
func Sugar() *zap.SugaredLogger {
	return packageLogger.Sugar()
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return packageLogger.WithOptions(opts...)
}

func With(fields ...zap.Field) *zap.Logger {
	return packageLogger.With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	packageLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	packageLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	packageLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	packageLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	packageLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	packageLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	packageLogger.Fatal(msg, fields...)
}

func Sync() error {
	return packageLogger.Sync()
}

// 获取日志级别
func getLogLevel(level string) zapcore.Level {
	return levelMap[level]
}

// 获取日志文件名
func getFileName(file string) string {
	// 有指定文件名，则直接使用
	if len(file) > 0 {
		return file
	}

	// 使用程序名作为日志文件名
	fileAbsolutePath, _ := exec.LookPath(os.Args[0])
	return fileAbsolutePath + ".log"
}

// 下面生成自定义的日志接口

// 初始化开发模式的logger
// filename 可为空，为空则使用程序名为日志文件名
// level  为"debug", "info",  "warn",  "error", "dpanic", "panic", "fatal"
// maxSizeMB 单个日志文件最大，单位为MB
// maxAgeDay 单个文件最多保留的天数
// maxBackups 最多保留之前的几个日志文件
func NewDevLogger(filename, level string, maxSizeMB, maxAgeDay, maxBackups int) *zap.Logger {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	hook := &lumberjack.Logger{
		Filename:   getFileName(filename), // 日志文件路径
		MaxSize:    maxSizeMB,             // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     maxAgeDay,             // 文件最多保存多少天
		MaxBackups: maxBackups,            // 日志文件最多保存多少个备份
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLogLevel(level))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),                      // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 显示调用者，文件及行号
	caller := zap.AddCaller()
	// DPanic会让程序产生panic
	development := zap.Development()
	// 开启堆栈跟踪
	stackTrace := zap.AddStacktrace(zap.WarnLevel)

	return zap.New(core, caller, development, stackTrace)
}

// 初始化生产模式的logger
func NewProdLogger(filename, level string, maxSizeMB, maxAgeDay, maxBackups int) *zap.Logger {
	hook := &lumberjack.Logger{
		Filename:   getFileName(filename), // 日志文件路径
		MaxSize:    maxSizeMB,             // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     maxAgeDay,             // 文件最多保存多少天
		MaxBackups: maxBackups,            // 日志文件最多保存多少个备份
	}

	// 格式化时间
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLogLevel(level))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                          // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 显示调用者，文件及行号
	caller := zap.AddCaller()
	// 开启堆栈跟踪
	stackTrace := zap.AddStacktrace(zap.WarnLevel)

	return zap.New(core, caller, stackTrace)
}
