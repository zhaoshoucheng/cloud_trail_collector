package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
)

// ConsoleOutPutCollector 输出到控制台
type ConsoleOutPutCollector struct {
}

func NewConsoleOutputCollector(conf *config.OutPutConfig) (OutPutCollector, error) {
	return &ConsoleOutPutCollector{}, nil
}

func (t *ConsoleOutPutCollector) GetName() string {
	return "ConsoleOutPutCollector"
}
func (t *ConsoleOutPutCollector) Insert(ctx context.Context, data interface{}) error {
	fmt.Println(data)
	return nil
}

func init() {
	registerOutPutCollector("console", NewConsoleOutputCollector)
}
