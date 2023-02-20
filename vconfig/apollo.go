package vconfig

import (
	"context"
	"github.com/Snowlights/tool/parse"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"strconv"
	"strings"
)

type ApolloCenter struct {
	c agollo.Client
}

func NewCenter(cc *CenterConfig) (Center, error) {

	if len(cc.IP) == 0 {
		return nil, InvalidIp
	}

	if cc.Port <= portMin || cc.Port > portMax {
		return nil, InvalidPort
	}

	if len(cc.Namespace) == 0 {
		cc.Namespace = []string{Application}
	}

	if cc.IsBackupConfig && len(cc.BackupConfigPath) == 0 {
		return nil, InvalidBackup
	}

	c := &config.AppConfig{
		AppID:             ReplaceServiceName(cc.AppID),
		Cluster:           cc.Cluster,
		NamespaceName:     strings.Join(cc.Namespace, comma),
		IP:                cc.IP + colon + strconv.FormatInt(int64(cc.Port), 10),
		IsBackupConfig:    cc.IsBackupConfig,
		BackupConfigPath:  cc.BackupConfigPath,
		Secret:            cc.SecretKey,
		SyncServerTimeout: cc.SyncServerTimeout,
		MustStart:         cc.MustStart,
	}

	agollo.SetLogger(&CenterLogger{c: context.Background()})
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		return nil, err
	}

	return &ApolloCenter{c: client}, nil
}

func (c ApolloCenter) GetValue(key string) (string, bool) {
	return c.GetValueWithNamespace(Application, key)
}

func (c ApolloCenter) GetValueWithNamespace(namespace, key string) (string, bool) {
	return c.getValue(namespace, key)
}

func (c ApolloCenter) UnmarshalWithNameSpace(namespace, tag string, v interface{}) error {
	kv := c.getAllKeyValues(namespace)
	return parse.UnmarshalKV(kv, v, tag)
}

func (c ApolloCenter) AddListener(listener storage.ChangeListener) {
	c.c.AddChangeListener(listener)
}

func (c ApolloCenter) getValue(namespace, key string) (string, bool) {
	cacheConfig := c.getConfig(namespace)
	if cacheConfig == nil {
		return "", false
	}
	val := cacheConfig.GetValue(key)
	if len(val) > 0 {
		return val, true
	}
	return "", false
}

func (c ApolloCenter) getAllKeyValues(namespace string) map[string]string {
	cacheConfig := c.getConfig(namespace)
	if cacheConfig == nil {
		return nil
	}
	keyValues := make(map[string]string)
	cacheConfig.GetCache().Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if !ok {
			return false
		}
		valStr, ok := c.getValue(namespace, keyStr)
		if !ok {
			return false
		}
		keyValues[keyStr] = valStr
		return true
	})
	return keyValues
}

func (c ApolloCenter) getConfig(namespace string) *storage.Config {
	return c.c.GetConfig(namespace)
}
