package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
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
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/spf13/viper"

	"eyes/setting"
	"eyes/utility"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/sessions"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

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

func (ai *AppInitializer) Init() *AppInitializer {
	ai.logger = ai.initLogger()
	ai.esClient = ai.initElasticsearch()
	ai.node = ai.initSnowflakeNode()
	ai.redisClient = ai.initRedisClient()
	ai.store = ai.initSessionStore()
	ai.initValidatorTranslator("zh")

	ai.bDB = ai.initDatabase("xiaohongshu_b", ai.logger)
	ai.cDB = ai.initDatabase("xiaohongshu_c", ai.logger)
	ai.setupDatabaseCallbacks(ai.cDB)
	return ai
}

func (ai *AppInitializer) initLogger() *zap.Logger {
	logger, err := middleware.NewLogger(ai.conf.LogConfig, ai.conf.Mode)
	if err != nil {
		fmt.Printf("初始化日志失败，错误:%v\n", err)
		panic(err)
	}
	return logger
}

func (ai *AppInitializer) initElasticsearch() *elasticsearch.Client {
	esClient, err := utility.EsInit()
	if err != nil {
		fmt.Printf("初始化Elasticsearch失败，错误: %v\n", err)
		panic(err)
	}
	return esClient
}

func (ai *AppInitializer) initSnowflakeNode() utility.ISFNode {
	node, err := utility.NewSFNode(ai.conf.StartTime, ai.conf.MachineID)
	if err != nil {
		fmt.Printf("初始化雪花算法节点失败，错误: %+v\n", err)
		panic(err)
	}
	return node
}

func (ai *AppInitializer) initRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", ai.conf.RedisConfig.Host, ai.conf.RedisConfig.Port),
		Password:     ai.conf.RedisConfig.Password,
		DB:           ai.conf.RedisConfig.DB,
		PoolSize:     ai.conf.RedisConfig.PoolSize,
		MinIdleConns: ai.conf.RedisConfig.MinIdleConns,
	})
	return redisClient
}

func (ai *AppInitializer) initSessionStore() sessionRedis.Store {
	store, err := sessionRedis.NewStoreWithDB(10, "tcp",
		fmt.Sprintf("%s:%d", ai.conf.RedisConfig.Host, ai.conf.RedisConfig.Port),
		ai.conf.RedisConfig.Password, strconv.Itoa(ai.conf.RedisConfig.DB), []byte("secret"),
	)
	if err != nil {
		fmt.Println("初始化Redis会话存储失败: ", err)
		panic(err)
	}
	return store
}

func (ai *AppInitializer) initValidatorTranslator(locale string) {
	if err := common.InitTrans(locale); err != nil {
		fmt.Printf("初始化验证器翻译器失败，错误:%v\n", err)
		panic(err)
	}
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

func (ai *AppInitializer) setupDatabaseCallbacks(db *gorm.DB) {
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
}

func Init(conf *setting.AppConfig) *AppInitializer {
	ai := NewAppInitializer(conf)
	return ai.Init()
}

func MustLoadConfig() *setting.AppConfig {
	var env string
	flag.StringVar(&env, "env", "dev", "env must in [dev,release,prod]")
	flag.Parse()

	if !(env == "dev" || env == "release" || env == "prod") {
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

func main() {
	conf := MustLoadConfig()
	ai := Init(conf)

	// logger, esClient, node, redisClient, store, bDB, cDB := Init(conf)
	logger, esClient, node, redisClient, bDB, cDB := ai.logger, ai.esClient, ai.node, ai.redisClient, ai.bDB, ai.cDB

	logger.Info("env-mode:", zap.String("env-mode", conf.Mode))
	logger.Info("sentry-dsn: ", zap.String("sentry-dsn", viper.GetString("sentry.dsn")))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:                viper.GetString("sentry.dsn"),
		EnableTracing:      true,
		TracesSampleRate:   0.2,
		ProfilesSampleRate: 0.2,
	})
	if err != nil {
		panic(err)
	}
	if conf.Mode == RELEASE || conf.Mode == PROD {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(DEBUG)
	}

	engine := gin.New()
	engine.Use(sentrygin.New(sentrygin.Options{}))
	engine.Use(requestid.New())
	engine.Use(sessions.Sessions("session-id", ai.store))
	engine.Use(middleware.GinLogger(logger), middleware.GinRecovery(true, logger))

	//err := engine.SetTrustedProxies([]string{"*"})
	//if err != nil {
	//	panic(err)
	//}

	engine.LoadHTMLGlob("templates/*")

	// profile-performance
	{
		if gin.Mode() == gin.DebugMode {
			pprof.Register(engine)
		}
	}

	// general-router
	{
		engine.NoRoute(func(c *gin.Context) {
			c.HTML(http.StatusNotFound, "404.html", nil)
		})

		engine.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 200, "success": "ok"})
		})
	}

	// monitor
	{

		morCache := cache.NewRedisCache(redisClient)
		morESDao := es.NewMonitorDAO(esClient)

		morESRepo := repository.NewMonitorRepository(morESDao, morCache, node)
		morSrv := serviceAlert.NewMonitorService(morESRepo)

		ctl := alert.NewMonitorController(morSrv, logger)
		ctl.RegisterRoutes(engine)
	}

	// article
	{
		artCache := cache.NewRedisCache(redisClient)
		cArtDao := dao.NewArticleDAO(cDB, node)
		bArtDao := dao.NewArticleDAO(bDB, node)

		cArtRepo := repository.NewArticleRepository(cArtDao, artCache)
		bArtRepo := repository.NewArticleRepository(bArtDao, artCache)
		artSrv := serviceArticle.NewArticleService(cArtRepo, bArtRepo)

		ctl := article.NewArticleController(artSrv, logger)
		ctl.RegisterRoutes(engine)
	}

	// azure aad登录
	{
		azureDao := dao.NewAzureDAO(cDB)
		azureRepo := repository.NewAzureRepository(azureDao)
		azureSrv := serviceAzure.NewAzureService(azureRepo)

		ctl := aad.NewAzureAzureController(azureSrv, logger)
		ctl.RegisterRoutes(engine)
	}

	// login
	{
		loginDao := dao.NewUserDAO(cDB, node)
		loginCache := cache.NewRedisCache(redisClient)
		loginRepo := repository.NewUserRepository(loginDao, loginCache)
		loginSrv := serviceLogin.NewLoginService(loginRepo)
		loginCtl := login.NewController(loginSrv, logger)
		loginCtl.RegisterRoutes(engine)
	}

	// ds-login

	{
		redisClient1 := redis.NewClient(&redis.Options{
			Addr: "localhost:26379",
		})

		dsDao := dao.NewDsUserDAO(cDB)
		redisCache := cache.NewRedisCache(redisClient1)
		dsRepo := repository.NewDsLoginRepository(dsDao, redisCache)
		dsLoginSrv := serviceDsLogin.NewLoginService(dsRepo)
		dsLoginCtl := dslogin.NewDsLoginController(dsLoginSrv, logger)
		dsLoginCtl.RegisterRoutes(engine)

	}

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	fmt.Printf("%s service starts running.", addr)
	err = engine.Run(addr)
	if err != nil {
		logger.Error("服务异常", zap.Error(err))
	}
}
