package Util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"sync/atomic"
)

// 全局日志对象（推荐单例模式，避免重复初始化）
var (
	Logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
	logLevel    atomic.Int32 // 原子变量，支持动态调整日志级别（0=Debug, 1=Info, 2=Warn, 3=Error, 4=DPanic, 5=Panic, 6=Fatal）
)

// 初始化日志配置（入口函数）
// serviceName: 服务名（固定字段）
// env: 环境（dev/test/prod）
// logPath: 日志文件路径（prod 环境必填）
func Pre() (err error) {
	// 初始化日志级别（默认 Info）
	logLevel.Store(int32(zapcore.InfoLevel))

	// 根据环境选择编码器和输出
	var core zapcore.Core
	//switch env {
	//case "dev", "test":
	//	core, err = newDevCore()
	//case "prod":
	//	//core, err = newProdCore(logPath)
	//default:
	//	return fmt.Errorf("invalid env: %s (support dev/test/prod)", env)
	//}
	//if err != nil {
	//	return err
	//}

	core, err = newDevCore()

	// 添加固定字段（全局生效，所有日志都会包含）
	baseFields := zap.Fields(
		//zap.String("service", serviceName),
		//zap.String("env", env),
		zap.Int("pid", os.Getpid()),
		zap.String("go_version", runtime.Version()),
	)

	// 构建 Logger（添加调用者信息、堆栈追踪）
	Logger = zap.New(
		core,
		zap.WithCaller(true),             // 输出调用文件和行号（生产环境建议保留，便于问题定位）
		zap.AddCallerSkip(1),             // 跳过当前函数（InitLogger），显示真实调用位置
		zap.AddStacktrace(zap.InfoLevel), // 仅 Error 及以上级别输出堆栈
		baseFields,                       // 固定字段
		zap.Hooks(metricHook),            // 自定义钩子（示例：日志 metrics 统计）
	)

	// 构建 SugaredLogger（更易用的 API，支持 fmt 风格格式化）
	sugarLogger = Logger.Sugar()

	// 注册退出时同步日志（确保缓冲区日志写入文件）
	return nil
}

func Post() {
	defer Logger.Sync()
}

// newDevCore：开发环境核心配置（控制台彩色输出、人类可读格式）
func newDevCore() (zapcore.Core, error) {
	// 控制台输出器（彩色）
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey, // 不输出函数名（简化输出）
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色级别（DEBUG/INFO/WARN/ERROR）
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间格式：2006-01-02T15:04:05.000Z0700
		EncodeDuration: zapcore.StringDurationEncoder,    // 耗时格式：1.234ms
		EncodeCaller:   zapcore.ShortCallerEncoder,       // 调用者格式：pkg/file.go:123
	})

	// 动态级别检查器（支持运行时调整级别）
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(logLevel.Load())
	})

	// 输出到控制台
	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), levelEnabler), nil
}

// newProdCore：生产环境核心配置（JSON 格式、日志轮转、采样）
//func newProdCore(logPath string) (zapcore.Core, error) {
//	// 确保日志目录存在
//	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
//		return nil, fmt.Errorf("create log dir failed: %w", err)
//	}
//
//	// 日志轮转配置（lumberjack）
//	rotator := &lumberjack.Logger{
//		Filename:   logPath, // 日志文件路径（如：./logs/app.log）
//		MaxSize:    100,     // 单个文件最大大小（MB）
//		MaxBackups: 7,       // 保留历史文件数（7天）
//		MaxAge:     7,       // 保留历史文件天数
//		Compress:   true,    // 压缩历史文件（gzip）
//		LocalTime:  true,    // 使用本地时间命名备份文件
//	}
//
//	// JSON 编码器（结构化输出，便于日志收集）
//	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
//		TimeKey:        "timestamp",
//		LevelKey:       "level",
//		NameKey:        "Logger",
//		CallerKey:      "caller",
//		FunctionKey:    zapcore.OmitKey,
//		MessageKey:     "message",
//		StacktraceKey:  "stacktrace",
//		LineEnding:     zapcore.DefaultLineEnding,
//		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 级别小写（debug/info/warn/error）
//		EncodeTime:     zapcore.RFC3339NanoEncoder,     // 时间格式（带纳秒，便于排序）：2006-01-02T15:04:05.999999999Z07:00
//		EncodeDuration: zapcore.SecondsDurationEncoder, // 耗时格式（秒，便于统计）：1.234
//		EncodeCaller:   zapcore.FullCallerEncoder,      // 调用者格式（完整路径）：github.com/xxx/pkg/file.go:123
//	})
//
//	// 采样策略（避免日志风暴：1秒内最多记录100条相同级别日志，超出的采样）
//	sampler := zapcore.NewSamplerWithOptions(
//		zapcore.LevelEnablerFunc(func(lvl zapcore.Level) bool {
//			return lvl >= zapcore.Level(logLevel.Load())
//		}),
//		time.Second, // 采样窗口
//		100,         // 窗口内最大日志数
//		10,          // 超出后每N条记录1条
//	)
//
//	// 输出到文件（结合轮转和采样）
//	return zapcore.NewCore(jsonEncoder, zapcore.AddSync(rotator), sampler), nil
//}

// 自定义钩子：日志 metrics 统计（示例）
func metricHook(entry zapcore.Entry) error {
	// 可扩展：记录不同级别日志的数量（如 Prometheus 指标）
	// 示例：prometheus.CounterVec{Name: "log_total", Labels: {"level", "service"}}.WithLabelValues(entry.Level.String(), serviceName).Inc()
	return nil
}

// 动态调整日志级别（支持运行时调用，如通过HTTP接口）
func SetLogLevel(lvl string) error {
	level, err := zapcore.ParseLevel(lvl)
	if err != nil {
		return err
	}
	logLevel.Store(int32(level))
	sugarLogger.Infof("log level updated to %s", lvl)
	return nil
}

// 获取当前日志级别
func GetLogLevel() string {
	return zapcore.Level(logLevel.Load()).String()
}
