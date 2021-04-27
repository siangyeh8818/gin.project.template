package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	conf "github.com/pnetwork/marvin.connector/internal/config"
	"github.com/siangyeh8818/gin.project.template/internal/handler"
	//"github.com/siangyeh8818/gin.project.template/internal/prometheus"
)

var (
	appConfig conf.Config
	dbPath    string
)

// DBFile indicate where the database is
const DBFile = "marvin-connector.db"

// ServerRun start the server
func ServerRun() {
	fmt.Println("Server Run")
	engine := gin.Default()

	engine.Use(timeoutMiddleware(time.Second * 3600))

	// TODO: database middleware
	// Add your middleware here
	// go func() {
	// 	engine.Use(ConnectDBMiddleware)
	// }()

	// get db connection
	// DBPath indicate where the database is
	/*
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" { // hard coded for development
			dbPath = "./" + DBFile
		} else {
			dbPath = dbPath + "/" + DBFile
		}
		fmt.Printf("DB_PATH: %s\n", dbPath)
		db, _ := database.GetSession(dbPath)
		defer db.Close()

	*/
	// get config path
	configPath := os.Getenv("PATH_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	fmt.Printf("PATH_CONFIG: %s\n", configPath)

	k8sSecret := conf.BaseConfig{}
	k8sSecret.InitConfig(configPath)

	appConfig = conf.Init(k8sSecret)

	env := &handler.Env{
		DBSession: db,
		AppConfig: &appConfig,
	}

	go env.CheckNewPublished()
	if appConfig.BaseConfig.CHECK_UPGRADE_STATUS_ENABLED {
		go env.CheckUpgradeStatus()
	}

	engine.POST("/marvin-connector/test", env.TestUpgrade)
	engine.NoRoute(env.NoRouteHandler)
	engine.GET("/marvin-connector/software/info", env.GetSoftwareInfo)
	engine.GET("/marvin-connector/version/:version", env.GetVersion)
	engine.POST("/marvin-connector/upgrade", env.PostUpgrade)
	engine.GET("/marvin-connector/upgrade", env.GetUpgrade)
	engine.PATCH("/marvin-connector/upgrade", env.PatchUpgrade)
	engine.GET("/marvin-connector/upgrade/log", env.GetUpgradeLog)
	engine.GET("/health", handler.LivenessCheck)
	engine.GET("/healthz", handler.ReadinessCheck)

	engine.Run()
}

// MetricRun exposes prometheus metrics
func MetricRun() {
	// TODO: another port for prometheus metrics
	metrics := gin.Default()

	metrics.GET("/metrics", prometheus.GetMetrics)
	metrics.Run(":9987")
}

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
