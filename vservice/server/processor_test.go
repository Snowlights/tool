package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
	"vtool/vprometheus/vmetric"
	"vtool/vservice/common"
)

type TestProcessor struct{}

func (mp *TestProcessor) Prepare() error {
	// set default metric
	return nil
}

func (mp *TestProcessor) Engine() (string, interface{}) {

	engine := gin.New()
	engine.Use(gin.Recovery())

	return "", engine
}

func TestServ(t *testing.T) {

	namespace, subsystem, name := "namespace", "subsystem", "histogram"

	bucketMin, bucketMax, bucketStep := 0, 100, 10
	buckets := []float64{}
	for i := bucketMin; i <= bucketMax; i += bucketStep {
		buckets = append(buckets, float64(i))
	}

	histogram := vmetric.NewHistogram(&vmetric.VecOpts{
		NameSpace:  namespace,
		SubSystem:  subsystem,
		Name:       name,
		Help:       "help",
		LabelNames: []string{"a"},
		Buckets:    buckets,
	}).With("a", "1")

	go func() {
		for {
			for i := 0; i < 100; i++ {
				histogram.Observe(float64(i))
			}
			time.Sleep(time.Second * 15)
		}
	}()

	err := Serv(context.Background(), &common.RegisterConfig{
		RegistrationType: common.ZOOKEEPER,
		ServName:         "censor",
		Group:            "/group/base",
	}, map[string]common.Processor{
		"test": &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}
