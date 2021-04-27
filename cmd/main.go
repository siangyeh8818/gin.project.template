import (
	server "github.com/siangyeh8818/gin.project.template/internal/http"
)

func main() {
	go server.MetricRun()
	server.ServerRun()
}