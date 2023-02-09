package breaker

import (
	"bytes"
	"fmt"
	"github.com/Snowlights/tool/parse"
	"github.com/Snowlights/tool/vconfig"
	"github.com/Snowlights/tool/vservice/common"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//
const (
	ticker             = time.Millisecond * 100
	defaultGranularity = time.Second * 1
	defaultThreshold   = 10
	defaultBreakerGap  = 10 // seconds
)

var (
	ErrTriggerBreaker = func(serv, funcName string) error {
		return fmt.Errorf("trigger serv breaker, serv info is %s, function is %s", serv, funcName)
	}
)

// 计数熔断
type Manager struct {
	group, name string

	center vconfig.Center
	cMu    sync.RWMutex
	cfg    *Config

	bMu      sync.Mutex
	Breakers map[string]*Breaker
}

type Breaker struct {
	m *Manager

	Rejected      int32
	RejectedStart int64
	Count         int32
}

var bm *Manager

func InitBreakerManager(group, name string) error {
	bm = &Manager{
		group:    group,
		name:     name,
		Breakers: make(map[string]*Breaker),
	}

	centerConfig, err := bm.parseConfigEnv()
	if err != nil {
		return err
	}

	center, err := vconfig.NewCenter(centerConfig)
	if err != nil {
		return err
	}

	bm.center = center

	bm.center.AddListener(&common.ClientListener{Change: bm.reloadConfig})
	bm.reloadConfig()

	return nil
}

// StatBreaker state errors for breaker
func StatBreaker(cluster, funcName string, err error) {
	if err != nil && (strings.Contains(err.Error(), timeoutErr) || strings.Contains(err.Error(), connectionErr)) {
		key := concat(cluster, underLine, funcName)
		bm.bMu.Lock()
		if _, ok := bm.Breakers[key]; !ok {
			breaker := new(Breaker)
			breaker.m = bm
			breaker.run()
			bm.Breakers[key] = breaker
		}
		breaker := bm.Breakers[key]
		bm.bMu.Unlock()
		atomic.AddInt32(&breaker.Count, 1)
	}
}

// Entry check if allow request
func Entry(cluster, funcName string) bool {
	key := concat(cluster, underLine, funcName)
	bm.bMu.Lock()
	breaker := bm.Breakers[key]
	bm.bMu.Unlock()
	if breaker != nil {
		return atomic.LoadInt32(&breaker.Rejected) != 1
	}
	return true
}

func (m *Manager) parseConfigEnv() (*vconfig.CenterConfig, error) {
	centerConfig, err := vconfig.ParseConfigEnv()
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(centerConfig.Port, 10, 64)
	if err != nil {
		return nil, err
	}

	return &vconfig.CenterConfig{
		AppID:            m.group + common.Slash + m.name,
		Cluster:          centerConfig.Cluster,
		Namespace:        []string{vconfig.Breaker},
		IP:               centerConfig.IP,
		Port:             int(port),
		IsBackupConfig:   false,
		BackupConfigPath: "",
		MustStart:        centerConfig.MustStart,
	}, nil
}

func (m *Manager) reloadConfig() {
	cfg := new(Config)
	err := m.center.UnmarshalWithNameSpace(vconfig.Breaker, parse.PropertiesTagName, cfg)
	if err != nil {
		return
	}

	m.cMu.Lock()
	defer m.cMu.Unlock()
	m.cfg = cfg
}

func (m *Manager) ticker() time.Duration {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	if m.cfg.Ticker > 0 {
		return time.Duration(m.cfg.Ticker) * time.Millisecond
	}
	return ticker
}

func (m *Manager) granularity() time.Duration {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	if m.cfg.Granularity > 0 {
		return time.Duration(m.cfg.Granularity) * time.Millisecond
	}
	return defaultGranularity
}

func (m *Manager) threshold() int64 {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	if m.cfg.Threshold > 0 {
		return m.cfg.Threshold
	}
	return defaultThreshold
}

func (m *Manager) breakerGap() int64 {
	m.cMu.RLock()
	defer m.cMu.RUnlock()

	if m.cfg.BreakerGap > 0 {
		return m.cfg.BreakerGap
	}
	return defaultBreakerGap
}

func (b *Breaker) run() {
	go func() {
		granularityTickC := time.Tick(b.m.granularity())
		checkTickC := time.Tick(b.m.ticker())
		for {
			select {
			case <-granularityTickC:
				atomic.StoreInt32(&b.Count, 0)
				// check 1s/checkTick times in 1s
			case <-checkTickC:
				threshold := b.m.threshold()
				breakerGap := b.m.breakerGap()
				if atomic.LoadInt32(&b.Count) > int32(threshold) {
					atomic.StoreInt32(&b.Rejected, 1)
					b.RejectedStart = time.Now().Unix()
				} else {
					now := time.Now().Unix()
					if now-b.RejectedStart > int64(breakerGap) {
						atomic.StoreInt32(&b.Rejected, 0)
					}
				}
			}
		}
	}()
}

func concat(str ...string) string {
	var buffer bytes.Buffer
	for _, s := range str {
		buffer.WriteString(s)
	}
	return buffer.String()
}
