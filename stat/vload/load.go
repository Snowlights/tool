package vload

import "github.com/shirou/gopsutil/load"

type LoadType int64

const (
	OneMin     LoadType = 1
	FiveMin    LoadType = 2
	FifteenMin LoadType = 3
)

func Load(loadType LoadType) (float64, error) {

	loadInfo, err := load.Avg()
	if err != nil {
		return 0, err
	}

	switch loadType {
	case OneMin:
		return loadInfo.Load1, nil
	case FiveMin:
		return loadInfo.Load5, nil
	case FifteenMin:
		return loadInfo.Load15, nil
	}

	return loadInfo.Load1, nil
}
