package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
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

func TestStop(t *testing.T) {

}

func TestServ(t *testing.T) {

	//namespace, subsystem, name := "namespace", "subsystem", "histogram"
	//
	//bucketMin, bucketMax, bucketStep := 0, 100, 10
	//buckets := []float64{}
	//for i := bucketMin; i <= bucketMax; i += bucketStep {
	//	buckets = append(buckets, float64(i))
	//}
	//
	//histogram := vmetric.NewHistogram(&vmetric.VecOpts{
	//	NameSpace:  namespace,
	//	SubSystem:  subsystem,
	//	Name:       name,
	//	Help:       "help",
	//	LabelNames: []string{"a"},
	//	Buckets:    buckets,
	//}).With("a", "1")
	//
	//go func() {
	//	for {
	//		for i := 0; i < 100; i++ {
	//			histogram.Observe(float64(i))
	//		}
	//		time.Sleep(time.Second * 15)
	//	}
	//}()

	err := ServService(map[common.ServiceType]common.Processor{
		common.ServiceTypeGin: &TestProcessor{},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

}
