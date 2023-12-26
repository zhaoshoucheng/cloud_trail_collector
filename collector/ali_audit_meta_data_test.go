package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"testing"
)

func TestAliTest(t *testing.T) {
	conf := &config.InputConfig{
		Type:      "ali_cloud_trail",
		EndPoint:  "actiontrail.ap-southeast-1.aliyuncs.com",
		AccessKey: "xxxxxxxx",
		SecretKey: "xxxxxx",
	}

	collector, err := NewAliCloudTrailInputCollector(conf)
	if err != nil {
		panic(err)
	}
	respChan := make(chan interface{})
	go func(resp chan interface{}) {
		for {
			select {
			case msg := <-respChan:
				fmt.Println(msg)
			}
		}
	}(respChan)
	collector.Update(context.Background(), respChan)
	select {}

}
