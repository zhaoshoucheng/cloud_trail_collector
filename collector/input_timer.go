package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"time"
)

// TimerInputCollector timer input 每2s中发送一个字符串，可用于测试
type TimerInputCollector struct {
}

func (t *TimerInputCollector) GetName() string {
	return "TimerInputCollector"
}
func (t *TimerInputCollector) Update(ctx context.Context, ch chan interface{}) {
	timer := time.NewTicker(time.Second * 2)
	count := 0
	for {
		select {
		case <-timer.C:
			count++
			if count >= 10000 {
				count = 0
			}
			ch <- fmt.Sprintf("timer input ch %d", count)
		}
	}
}

func NewTimerInputCollector(input *config.InputConfig) (InputCollector, error) {
	return &TimerInputCollector{}, nil
}

func init() {
	registerInputCollector("timer", NewTimerInputCollector)
}
