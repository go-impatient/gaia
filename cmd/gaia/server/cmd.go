package server

import (
	"errors"
	"fmt"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/go-impatient/gaia/app/conf"
	"github.com/go-impatient/gaia/app/router"
	"github.com/go-impatient/gaia/internal/repository"
	"github.com/go-impatient/gaia/internal/service"
	"github.com/go-impatient/gaia/pkg/http/ginhttp"
	"github.com/go-impatient/gaia/pkg/logger"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		EnvVars: []string{"GAIA_CONFIG"},
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "设置配置文件",
		Value:   "",
	},
}

var Cmd = &cli.Command{
	Name:     "server",
	Usage:    "Gaia 应用管理",
	HideHelp: false,
	Subcommands: cli.Commands{
		&cli.Command{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "运行Gaia服务",
			Action:  start,
			Flags:   flags,
		},
		&cli.Command{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "部署Gaia服务",
			Action:  deploy,
			Flags:   flags,
		},
	},
}

// 运行服务
func start(c *cli.Context) error {
	// var group errgroup.Group
	fileName := c.String("config")
	if len(fileName) == 0 {
		return errors.New("server s -c ./../../config/config.json 或者 server start -config ./../../config/config.json")
	}
	// 获取文件后缀
	fileSuffix := path.Ext(fileName)
	if len(fileSuffix) > 1 {
		// 去掉后缀中.字符
		fileSuffix = fileSuffix[1:]
	}

	// 1. 初始化配置
	cfg, err := conf.InitConfig(fileName, fileSuffix)
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
	sql := service.NewSQL()
	defer sql.Close()

	// 4. 初始化 repositories and services
	userRepo := repository.NewUserRepository(sql.DB)
	adminRepo := repository.NewAdminRepository(sql.DB)
	userSrv := service.NewUserService(userRepo)
	adminSrv := service.NewAdminService(adminRepo)
	s := service.NewServices(userSrv, adminSrv)

	// 5. 初始化应用服务
	serve := ginhttp.NewServer()
	ginhttp.SetRuntimeMode(cfg.App.Mode)
	// 5.1 初始化路由
	r := serve.Router()
	router.RegisterRoutes(r, s)

	// 5.2 启动服务
	serve.Serve()
	//group.Go(func() error {
	//	return serve.RunHTTPServer()
	//})

	// 5.3 健康检查
	//group.Go(func() error {
	//	return serve.PingServer()
	//})

	//if err := group.Wait(); err != nil {
	//	log.Error().Msg(fmt.Sprintf("接口服务停止了：%v", err))
	//}
	//
	//return group.Wait()

	return nil
}

// 编译和部署
func deploy(c *cli.Context) error {
	var grop errgroup.Group
	port := c.Int("port")
	debug := c.Bool("debug")
	config := c.String("config")

	if debug {
		log.Printf("dev")
		// "dev"
	} else {
		log.Printf("prod")
		// "prod"
	}

	if port > 0 {
		return errors.New("start -p")
	}

	if len(config) == 0 {
		return errors.New("start -c")
	}

	log.Info().Msg(fmt.Sprintf("d: %v, p: %v, c: %v", debug, port, config))

	// 编译
	grop.Go(func() error {

		return nil
	})

	// 部署
	grop.Go(func() error {

		return nil
	})

	if err := grop.Wait(); err != nil {
		log.Error().Msg(fmt.Sprintf("编译或者部署失败: %v", err))
	}

	return grop.Wait()
}
