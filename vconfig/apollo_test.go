package vconfig

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"testing"
)

func TestNewConfig(t *testing.T) {

	c := &config.AppConfig{
		AppID:             "",
		Cluster:           "",
		NamespaceName:     "",
		IP:                "",
		IsBackupConfig:    false,
		BackupConfigPath:  "",
		Secret:            "",
		SyncServerTimeout: 0,
		MustStart:         false,
	}

	loadAppConfig := func() (*config.AppConfig, error) {
		return c, nil
	}

	client, _ := agollo.StartWithConfig(loadAppConfig)
	fmt.Println("初始化Apollo配置成功")

	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	value, _ := cache.Get("key")

	fmt.Println(value)
}
