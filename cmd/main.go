package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"github.com/raylin666/go-utils/auth"
	utLogger "github.com/raylin666/go-utils/logger"
	"github.com/raylin666/go-utils/server/system"
	"mt/config"
	"mt/internal/app"
	pkgLogger "mt/pkg/logger"
	"mt/pkg/utils"

	kratosConfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// ID service id
	ID string
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../.env.yaml", "config path, eg: -conf .env.yaml")
}

func newApp(tools *app.Tools, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(ID),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(tools.Logger()),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	c := kratosConfig.New(
		kratosConfig.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc config.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	ID = bc.App.Id
	Name = bc.App.Name
	Version = bc.App.Version

	// 打印启动信息
	app.NewLogo(&bc)

	// 初始化 Environment
	var env = system.NewEnvironment(bc.Environment)
	bc.Environment = env.Value()

	// 初始化 Datetime
	var datetime = system.NewDatetime(
		system.WithLocation(bc.Datetime.Location),
		system.WithCSTLayout(bc.Datetime.CstLayout))

	// 初始化 JWT 鉴权认证
	var jwt = auth.NewJWT(bc.Jwt.App, bc.Jwt.Key, bc.Jwt.Secret)

	// 初始化 Logger
	logger, err := pkgLogger.NewJSONLogger(
		// ut_logger.WithDisableConsole(),
		utLogger.WithField(utLogger.AppKey, bc.App.Name),
		utLogger.WithField(utLogger.EnvironmentKey, bc.Environment),
		utLogger.WithTimeLayout(bc.Datetime.CstLayout),
		//	项目访问日志存放文件
		utLogger.WithPathFileRotation(fmt.Sprintf("%s/runtime/logs/%s.log", utils.ProjectPath(1), bc.App.Name), utLogger.PathFileRotationOption{
			MaxSize:    int(bc.Log.MaxSize),
			MaxAge:     int(bc.Log.MaxAge),
			MaxBackups: int(bc.Log.MaxBackups),
			LocalTime:  bc.Log.LocalTime,
			Compress:   bc.Log.Compress,
		}))
	if err != nil {
		panic(err)
	}

	// 创建公共工具实例
	var tools = app.NewTools(logger, datetime, jwt)

	appMT, cleanup, err := wireApp(bc.Server, bc.Data, bc.App, bc.Websocket, tools)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := appMT.Run(); err != nil {
		panic(err)
	}
}
