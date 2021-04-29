package prometheus

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO
func GetMetrics(c *gin.Context) {
	fmt.Println("getMetrics")

	// Response
	defer c.JSON(http.StatusOK, `process_virtual_memory_max_bytes -1
	# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
	# TYPE promhttp_metric_handler_requests_in_flight gauge
	promhttp_metric_handler_requests_in_flight 1
	# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
	# TYPE promhttp_metric_handler_requests_total counter
	promhttp_metric_handler_requests_total{code="200"} 418
	promhttp_metric_handler_requests_total{code="500"} 0
	promhttp_metric_handler_requests_total{code="503"} 0`,
	)
}
