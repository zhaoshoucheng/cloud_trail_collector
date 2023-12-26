package collector

import (
	"cloud_trail_collector/config"
	"context"
	"testing"
)

func TestNewESOutPutCollector(t *testing.T) {
	conf := &config.OutPutConfig{
		Type:      "elastic",
		Condition: "",
		Index:     "upstream-2023.12.19",
		EndPoints: []string{"http://127.0.0.1:9200"},
		UserName:  "xxx",
		PassWord:  "xxx",
	}
	collector, err := NewESOutPutCollector(conf)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"@timestamp": "2023-12-19T18:02:13+08:00",
		"fall":       7668,
		"name":       "10.218.22.75:80",
		"port":       80,
		"rise":       0,
		"status":     "down",
		"timestamp":  "19/Dec/2023:18:02:13 +0800",
		"upstream":   "UFERSC.test.UFERSC",
	}
	err = collector.Insert(context.Background(), data)
	if err != nil {
		panic(err)
	}
}
