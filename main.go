package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	dslogin "eyes/internal/web/ds_login"

	"eyes/internal/domain"
	"eyes/internal/repository"
	"eyes/internal/repository/cache"
	"eyes/internal/repository/dao"
	"eyes/internal/repository/es"
	serviceAzure "eyes/internal/service/aad"
	serviceAlert "eyes/internal/service/alert"
	serviceArticle "eyes/internal/service/article"
	serviceDsLogin "eyes/internal/service/ds_login"
	serviceLogin "eyes/internal/service/login"
	"eyes/internal/web/aad"
	"eyes/internal/web/alert"
	"eyes/internal/web/article"
	"eyes/internal/web/common"
	"eyes/internal/web/login"
	"eyes/internal/web/middleware"
	"eyes/setting"
	"eyes/utility"

	"github.com/elastic/go-elasticsearch/v8"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

func (ai *AppInitializer) Init() *AppInitializer {
	ai.initLogger().
		initElasticsearch().
		initSnowflakeNode().
		initRedisClient().
		initSessionStore().
		initValidatorTranslator("zh").
		initDatabases().
		setupDatabaseCallbacks(ai.cDB)
	return ai
}

func (ai *AppInitializer) initDatabases() *AppInitializer {
	ai.bDB = ai.initDatabase("xiaohongshu_b", ai.logger)
	ai.cDB = ai.initDatabase("xiaohongshu_c", ai.logger)
	return ai
}

type AppInitializer struct {
	conf        *setting.AppConfig
	logger      *zap.Logger
	esClient    *elasticsearch.Client
	node        utility.ISFNode
	redisClient *redis.Client
	store       sessionRedis.Store
	bDB         *gorm.DB
	cDB         *gorm.DB
}

func NewAppInitializer(conf *setting.AppConfig) *AppInitializer {
	return &AppInitializer{conf: conf}
}

//
//func (ai *AppInitializer) Init() *AppInitializer {
//	ai.logger = ai.initLogger()
//	ai.esClient = ai.initElasticsearch()
//	ai.node = ai.initSnowflakeNode()
//	ai.redisClient = ai.initRedisClient()
//	ai.store = ai.initSessionStore()
//	ai.initValidatorTranslator("zh")
//
//	ai.bDB = ai.initDatabase("xiaohongshu_b", ai.logger)
//	ai.cDB = ai.initDatabase("xiaohongshu_c", ai.logger)
//	ai.setupDatabaseCallbacks(ai.cDB)
//	return ai
//}

func (ai *AppInitializer) initLogger() *AppInitializer {
	logger, err := middleware.NewLogger(ai.conf.LogConfig, ai.conf.Mode)
	if err != nil {
		fmt.Printf("初始化日志失败，错误:%v\n", err)
		panic(err)
	}
	ai.logger = logger
	return ai
}

func (ai *AppInitializer) initElasticsearch() *AppInitializer {
	esClient, err := utility.EsInit()
	if err != nil {
		fmt.Printf("初始化Elasticsearch失败，错误: %v\n", err)
		panic(err)
	}
	ai.esClient = esClient
	return ai
}

func (ai *AppInitializer) initSnowflakeNode() *AppInitializer {
	node, err := utility.NewSFNode(ai.conf.StartTime, ai.conf.MachineID)
	if err != nil {
		fmt.Printf("初始化雪花算法节点失败，错误: %+v\n", err)
		panic(err)
	}
	ai.node = node
	return ai
}

func (ai *AppInitializer) initRedisClient() *AppInitializer {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", ai.conf.RedisConfig.Host, ai.conf.RedisConfig.Port),
		Password:     ai.conf.RedisConfig.Password,
		DB:           ai.conf.RedisConfig.DB,
		PoolSize:     ai.conf.RedisConfig.PoolSize,
		MinIdleConns: ai.conf.RedisConfig.MinIdleConns,
	})

	ai.redisClient = redisClient
	return ai
}

func (ai *AppInitializer) initSessionStore() *AppInitializer {
	store, err := sessionRedis.NewStoreWithDB(10, "tcp",
		fmt.Sprintf("%s:%d", ai.conf.RedisConfig.Host, ai.conf.RedisConfig.Port),
		ai.conf.RedisConfig.Password, strconv.Itoa(ai.conf.RedisConfig.DB), []byte("secret"),
	)
	if err != nil {
		fmt.Println("初始化Redis会话存储失败: ", err)
		panic(err)
	}
	ai.store = store
	return ai
}

func (ai *AppInitializer) initValidatorTranslator(locale string) *AppInitializer {
	if err := common.InitTrans(locale); err != nil {
		fmt.Printf("初始化验证器翻译器失败，错误:%v\n", err)
		panic(err)
	}
	return ai
}

func (ai *AppInitializer) initDatabase(dbName string, logger *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		ai.conf.MySQLConfig.User,
		ai.conf.MySQLConfig.Password,
		ai.conf.MySQLConfig.Host,
		ai.conf.MySQLConfig.Port,
		dbName,
	)

	zapLogger := utility.NewZapGormLogger(logger)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zapLogger,
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&domain.DA{})
	if err != nil {
		panic(err)
	}
	return db
}

func (ai *AppInitializer) setupDatabaseCallbacks(db *gorm.DB) *AppInitializer {
	err := db.Callback().Create().Before("*").Register(
		"slow_query_cache_start", func(db *gorm.DB) {
			db.Set("start_time", time.Now())
		},
	)
	if err != nil {
		panic(err)
	}

	err = db.Callback().Create().After("*").Register("slow_query_create_end", func(db *gorm.DB) {
		startTime, ok := db.Get("start_time")
		if ok {
			duration := time.Since(startTime.(time.Time))
			ai.logger.Info("慢查询",
				zap.String("sql", db.Statement.SQL.String()),
				zap.Duration("duration", duration),
			)
		}
	})
	if err != nil {
		panic(err)
	}
	return ai
}

//
//func Init(conf *setting.AppConfig) *AppInitializer {
//	ai := NewAppInitializer(conf)
//	return ai.Init()
//}

func MustLoadConfig() *setting.AppConfig {
	var env string
	flag.StringVar(&env, "env", "dev", "env must in [dev,release,prod]")
	flag.Parse()

	if !(env == DEBUG || env == RELEASE || env == PROD) {
		fmt.Println("env must be in [dev release prod]")
		panic("env must be in [dev release prod]")
	}

	conf, err := setting.LoadConf(env)
	if err != nil {
		fmt.Printf("load configuration failed, err:%v\n", err)
		panic(err)
	}
	return conf
}

const (
	RELEASE = "release"
	DEBUG   = "debug"
	PROD    = "prod"
)

func registerRoutes(engine *gin.Engine, ai *AppInitializer) {
	// monitor
	{

		morCache := cache.NewRedisCache(ai.redisClient)
		morESDao := es.NewMonitorDAO(ai.esClient)

		morESRepo := repository.NewMonitorRepository(morESDao, morCache, ai.node)
		morSrv := serviceAlert.NewMonitorService(morESRepo)

		ctl := alert.NewMonitorController(morSrv, ai.logger)
		ctl.RegisterRoutes(engine)
	}

	// article
	{
		artCache := cache.NewRedisCache(ai.redisClient)
		cArtDao := dao.NewArticleDAO(ai.cDB, ai.node)
		bArtDao := dao.NewArticleDAO(ai.bDB, ai.node)

		cArtRepo := repository.NewArticleRepository(cArtDao, artCache)
		bArtRepo := repository.NewArticleRepository(bArtDao, artCache)
		artSrv := serviceArticle.NewArticleService(cArtRepo, bArtRepo)

		ctl := article.NewArticleController(artSrv, ai.logger)
		ctl.RegisterRoutes(engine)
	}

	// azure aad登录
	{
		azureDao := dao.NewAzureDAO(ai.cDB)
		azureRepo := repository.NewAzureRepository(azureDao)
		azureSrv := serviceAzure.NewAzureService(azureRepo)

		ctl := aad.NewAzureAzureController(azureSrv, ai.logger)
		ctl.RegisterRoutes(engine)
	}

	// login
	{
		loginDao := dao.NewUserDAO(ai.cDB, ai.node)
		loginCache := cache.NewRedisCache(ai.redisClient)
		loginRepo := repository.NewUserRepository(loginDao, loginCache)
		loginSrv := serviceLogin.NewLoginService(loginRepo)
		loginCtl := login.NewController(loginSrv, ai.logger)
		loginCtl.RegisterRoutes(engine)
	}

	// ds-login

	{
		redisClient1 := redis.NewClient(&redis.Options{
			Addr: "localhost:26379",
		})

		dsDao := dao.NewDsUserDAO(ai.cDB)
		redisCache := cache.NewRedisCache(redisClient1)
		dsRepo := repository.NewDsLoginRepository(dsDao, redisCache)
		dsLoginSrv := serviceDsLogin.NewLoginService(dsRepo)
		dsLoginCtl := dslogin.NewDsLoginController(dsLoginSrv, ai.logger)
		dsLoginCtl.RegisterRoutes(engine)

	}
}

//
//func main() {
//	conf := MustLoadConfig()
//	ai := Init(conf)
//
//	// logger, esClient, node, redisClient, store, bDB, cDB := Init(conf)
//	logger, esClient, node, redisClient, bDB, cDB := ai.logger, ai.esClient, ai.node, ai.redisClient, ai.bDB, ai.cDB
//
//	logger.Info("env-mode:", zap.String("env-mode", conf.Mode))
//	logger.Info("sentry-dsn: ", zap.String("sentry-dsn", viper.GetString("sentry.dsn")))
//
//	err := sentry.Init(sentry.ClientOptions{
//		Dsn:                viper.GetString("sentry.dsn"),
//		EnableTracing:      true,
//		TracesSampleRate:   0.2,
//		ProfilesSampleRate: 0.2,
//	})
//	if err != nil {
//		panic(err)
//	}
//	if conf.Mode == RELEASE || conf.Mode == PROD {
//		gin.SetMode(gin.ReleaseMode)
//	} else {
//		gin.SetMode(DEBUG)
//	}
//
//	engine := gin.New()
//	engine.Use(sentrygin.New(sentrygin.Options{}))
//	engine.Use(requestid.New())
//	engine.Use(sessions.Sessions("session-id", ai.store))
//	engine.Use(middleware.GinLogger(logger), middleware.GinRecovery(true, logger))
//
//	//err := engine.SetTrustedProxies([]string{"*"})
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	engine.LoadHTMLGlob("templates/*")
//
//	// profile-performance
//	{
//		if gin.Mode() == gin.DebugMode {
//			pprof.Register(engine)
//		}
//	}
//
//	// general-router
//	{
//		engine.NoRoute(func(c *gin.Context) {
//			c.HTML(http.StatusNotFound, "404.html", nil)
//		})
//
//		engine.GET("/health", func(c *gin.Context) {
//			c.JSON(http.StatusOK, gin.H{"code": 200, "success": "ok"})
//		})
//	}
//
//	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
//	fmt.Printf("%s service starts running.", addr)
//	err = engine.Run(addr)
//	if err != nil {
//		logger.Error("服务异常", zap.Error(err))
//	}
//}

func main() {
	conf := MustLoadConfig()
	// conf := setting.LoadConfig()
	appInitializer := NewAppInitializer(conf).Init()

	// 创建 Gin 引擎
	r := gin.Default()
	if conf.Mode == RELEASE || conf.Mode == PROD {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(DEBUG)
	}

	registerRoutes(r, appInitializer)

	// 定义服务器地址
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动服务器
	go func() {
		fmt.Printf("%s service starts running.", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appInitializer.logger.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	// 监听信号进行优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appInitializer.logger.Info("正在关闭服务器...")

	// 创建上下文对象，设置超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		appInitializer.logger.Fatal("服务器关闭失败", zap.Error(err))
	}

	appInitializer.logger.Info("服务器已关闭")
}
