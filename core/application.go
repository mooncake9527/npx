package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"time"

	"github.com/mooncake9527/npx/core/ebus"

	"os"
	"sync"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mooncake9527/npx/common/utils/ips"
	"github.com/mooncake9527/npx/common/utils/text"
	"github.com/mooncake9527/npx/config"
	"github.com/mooncake9527/npx/core/cache"
	"github.com/mooncake9527/npx/core/locker"
	"github.com/natefinch/lumberjack"
	"gorm.io/gorm"
)

var (
	Cfg       config.AppCfg
	iLog      *slog.Logger //*zap.Logger
	Cache     cache.ICache
	lock      sync.RWMutex
	engine    http.Handler
	dbs       = make(map[string]*gorm.DB, 0)
	RedisLock *locker.Redis
	Started   = make(chan byte, 1)
	ToClose   = make(chan byte, 1)
)

func GetEngine() http.Handler {
	return engine
}

func SetEngine(aEngine http.Handler) {
	engine = aEngine
}

func GetGinEngine() *gin.Engine {
	if Cfg.Server.Mode == ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	var r *gin.Engine
	lock.RLock()
	defer lock.RUnlock()
	if engine == nil {
		engine = gin.New()
	}
	switch engine.(type) {
	case *gin.Engine:
		r = engine.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
	}
	return r
}

func Init() {
	logWrite := logInit()
	Cache = cache.New(Cfg.Cache)
	if Cache.Type() == "redis" {
		r := Cache.(*cache.RedisCache)
		RedisLock = locker.NewRedis(r.GetClient())
	}
	dbInit(logWrite)
	ebus.EventBus.Publish(ebus.EventCoreInit)
}

func Run() {
	addr := fmt.Sprintf("%s:%d", Cfg.Server.GetHost(), Cfg.Server.GetPort())

	//服务启动参数
	srv := &http.Server{
		Addr:           addr,
		Handler:        GetEngine(),
		ReadTimeout:    time.Duration(Cfg.Server.GetReadTimeout()) * time.Second,
		WriteTimeout:   time.Duration(Cfg.Server.GetWriteTimeout()) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//fmt.Println(LOGO)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen: ", err)
		}
	}()

	//fmt.Println(text.Green(`orange github:`) + text.Blue(`https://github.com/mooncake/orange`))
	fmt.Println(text.Green("server started ,listen on: ") + text.Red("[ "+addr+" ]"))

	if Cfg.Server.Mode != ModeProd.String() {
		fmt.Println(text.Blue(fmt.Sprintf("swagger: http://localhost:%d/swagger/index.html", Cfg.Server.Port)))
		ip := ips.GetLocalHost()
		if ip != "" {
			fmt.Println(text.Blue(fmt.Sprintf("swagger: https://%s:%d/swagger/index.html", ip, Cfg.Server.Port)))
		}
	}
	ebus.EventBus.Publish(ebus.EventApplicationStarted)
	Started <- 1
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ebus.EventBus.Publish(ebus.EventApplicationQuit)

	ToClose <- 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slog.Info("server shutdown ...", "time", time.Now())

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server Shutdown:", err)
	}

	slog.Info("server exiting")
	time.Sleep(time.Second * time.Duration(Cfg.Server.GetCloseWait()))
}

func logInit() io.Writer {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	var logWriter io.Writer
	if Cfg.Logger.LogInConsole {
		logWriter = os.Stdout
	} else {
		logWriter = defaultLumberjack()
	}
	level := strings.ToLower(Cfg.Logger.Level)
	switch level {
	case "error":
		opts.Level = slog.LevelError
	case "warn":
		opts.Level = slog.LevelWarn
	default:
		opts.Level = slog.LevelInfo
	}
	if strings.ToLower(Cfg.Logger.Format) == "json" {
		iLog = slog.New(slog.NewJSONHandler(logWriter, &opts))
	} else {
		iLog = slog.New(slog.NewTextHandler(logWriter, &opts))
	}
	slog.SetDefault(iLog)
	return logWriter
}

func defaultLumberjack() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   Cfg.Logger.Director + "/orange.log",
		LocalTime:  true,
		MaxSize:    Cfg.Logger.GetMaxSize(),
		MaxAge:     Cfg.Logger.GetMaxAge(),
		MaxBackups: Cfg.Logger.GetMaxBackups(),
		Compress:   true,
	}
}
