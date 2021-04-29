package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siangyeh8818/gin.project.template/internal/routers"

	//"github.com/EDDYCJY/go-gin-example/models"
	"github.com/siangyeh8818/gin.project.template/pkg/gredis"
	"github.com/siangyeh8818/gin.project.template/pkg/logging"
	"github.com/siangyeh8818/gin.project.template/pkg/prometheus"
	"github.com/siangyeh8818/gin.project.template/pkg/setting"
	"github.com/siangyeh8818/gin.project.template/pkg/util"
)

func init() {
	setting.Setup()
	//models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
}

func main() {
	//server.ServerRun()
	//go server.MetricRun()
	ServerRun()
}

func ServerRun() {
	/*gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
	*/

	engine := gin.Default()

	// 導入路由
	engine = routers.InitRouter()

	engine.Run()
}

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
