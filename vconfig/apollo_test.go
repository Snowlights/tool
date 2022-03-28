package vconfig

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"testing"
	"time"
	"vtool/parse"
)

func TestNewConfig(t *testing.T) {

	c := &config.AppConfig{
		AppID:             "SampleApp",
		Cluster:           "dev",
		NamespaceName:     "application,client",
		IP:                "http://127.0.0.1:8080",
		IsBackupConfig:    true,
		BackupConfigPath:  "/Users/zhangwei/Desktop/apollo/tmp",
		Secret:            "",
		SyncServerTimeout: 0,
		MustStart:         true,
	}

	loadAppConfig := func() (*config.AppConfig, error) {
		return c, nil
	}

	client, err := agollo.StartWithConfig(loadAppConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("init apollo success")

	clientCache := client.GetConfig("client")
	val := clientCache.GetValue("f")
	fmt.Println(val)
	go func() {
		client.AddChangeListener(&Listener{})
	}()

	apolloCenter := &ApolloCenter{
		c: client,
	}
	cfg := &Cfg{}

	err = apolloCenter.UnmarshalWithNameSpace(Application, parse.PropertiesTagName, cfg)
	if err != nil {
		fmt.Println(err)
	}

	apolloCenter.AddListener(&Listener{})

	time.Sleep(time.Hour)
}

type Ency string

type Cfg struct {
	F        Ency          `json:"ddf" properties:"ddf"`
	Body     Body          `json:"ajj" properties:"ajj"`
	Timeouts time.Duration `json:"timeout" properties:"timeout"`
}

type Body struct {
	Values  []string  `json:"values" properties:"values"`
	Values1 []float64 `json:"values1" properties:"values1"`
}

type Listener struct {
}

func (l *Listener) OnChange(event *storage.ChangeEvent) {
	fmt.Println("apollo config changed", event)
}

func (l *Listener) OnNewestChange(event *storage.FullChangeEvent) {
	fmt.Println("apollo config changed", event)
}

func TestNew(t *testing.T) {

}
