package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/app/router"
	"github.com/go-impatient/gaia/internal/database"
	"github.com/go-impatient/gaia/internal/model"
	"github.com/go-impatient/gaia/internal/repository"
	"github.com/go-impatient/gaia/internal/service"
	"github.com/go-impatient/gaia/pkg/logger"
	"github.com/go-impatient/gaia/pkg/server"
	signal "github.com/go-impatient/gaia/pkg/server/siganl"
)

func Init(fileName string) {
	// 1. 初始化配置
	cfg, err := conf.InitConfig(fileName)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	//log.Info().Msg(fmt.Sprintf("配置参数: %+v", cfg.App))
	//log.Info().Msg(fmt.Sprintf("Config: %+v",conf.Config.App))
	//jsonBytes, _ := json.Marshal(conf.Config.App)
	//fmt.Println(string(jsonBytes))

	// 2. 初始化日志
	logger.InitLogger(cfg.Log.Level, cfg.Log.Format)

	// 3. 初始化数据库
	// Models represents all models..
	var Models = []interface{}{
		&model.App{},
		&model.User{},
		&model.AppMeta{},
		&model.UserMeta{},
	}
	sql := database.NewSQL()
	sql.Migrate(Models)
	defer sql.Close()

	// 4. 初始化 repositories and services
	repositories := repository.NewRepositories(sql.DB)
	services := service.NewServices(repositories)

	// 5. 初始化应用服务
	app := gin.New()
	serve := server.NewServer(server.App(app))
	switch cfg.App.Mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("unknown mode")
	}
	// 5.1 初始化路由
	router.RegisterRoutes(app, services)

	// 5.2 启动服务
	ctx := signal.WithContext(context.Background())
	serve.Serve(ctx)
}
