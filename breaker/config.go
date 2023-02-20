package breaker

const (
	timeoutErr    = "timeout"
	connectionErr = "invalid connection"

	underLine = "_"
	slash     = "/"
)

type Config struct {
	Ticker      int64 `json:"ticker" properties:"ticker"`
	Granularity int64 `json:"granularity" properties:"granularity"`
	Threshold   int64 `json:"threshold" properties:"threshold"`
	BreakerGap  int64 `json:"breaker_gap" properties:"breaker_gap"`
}

func (c *Config) ticker() int64 {
	return c.Ticker
}

func (c *Config) granularity() int64 {
	return c.Granularity
}

func (c *Config) threshold() int64 {
	return c.Threshold
}

func (c *Config) breakerGap() int64 {
	return c.BreakerGap
}
