// This file is auto-generated, don't edit it. Thanks.
package collector

import (
	"cloud_trail_collector/config"
	"context"
	"testing"
)

func TestAwsTest(t *testing.T) {
	input := &config.InputConfig{
		Type:      "aws_cloud_trail",
		AccessKey: "xxxxxx",
		SecretKey: "xxxxx",
		RegionId:  "ap-southeast-1",
	}
	collector, err := NewAwsCloudTrailInputCollector(input)
	if err != nil {
		panic(err)
	}
	respChan := make(chan interface{})
	go func(resp chan interface{}) {
		for {
			select {
			case msg := <-respChan:
				_ = msg
				//	fmt.Println(msg)
			}
		}
	}(respChan)
	collector.Update(context.Background(), respChan)
	select {}
}
