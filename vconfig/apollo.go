package vconfig

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/storage"
	"vtool/parse"
)

type ApolloCenter struct {
	c agollo.Client
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
	config := c.getConfig(namespace)
	if config == nil {
		return "", false
	}
	val := config.GetValue(key)
	if len(val) > 0 {
		return val, true
	}
	return "", false
}

func (c ApolloCenter) getAllKeyValues(namespace string) map[string]string {
	config := c.getConfig(namespace)
	if config == nil {
		return nil
	}
	keyValues := make(map[string]string)
	config.GetCache().Range(func(key, value interface{}) bool {
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
