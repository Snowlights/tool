package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"github.com/Snowlights/tool/vprometheus/metric"
	"github.com/Snowlights/tool/vservice/common"
	"github.com/Snowlights/tool/vservice/server"
)

type TestProcessor struct{}

func (mp *TestProcessor) Engine() (string, interface{}) {

	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.POST("/p", p)

	return "", engine
}

func p(c *gin.Context) {
	st := time.Now()
	defer func() {
		metric.StatApi("/p", time.Since(st))
	}()
	number := rand.Int63n(1000)
	time.Sleep(time.Millisecond * time.Duration(number))
	c.JSON(http.StatusOK, "hello world")
}

func TestHttpServer(t *testing.T) {

	err := server.ServService(map[common.ServiceType]common.Processor{
		common.HTTP: &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}
