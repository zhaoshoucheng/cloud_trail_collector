package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"testing"
)

func TestPipeline_Start(t *testing.T) {
	_, err := config.NewConfig("../config.toml")
	if err != nil {
		panic(err)
	}
	pipline := &Pipeline{}
	input, err := NewInputCollector(&config.InputConfig{Type: "timer"})
	if err != nil {
		panic(err)
	}
	output, err := NewOutputCollector(&config.OutPutConfig{Type: "console"})
	if err != nil {
		panic(err)
	}
	pipline.Input = input
	pipline.InputChan = make(chan interface{})
	pipline.OutPuts = []OutPutCollector{output, output}
	pipline.OutPutsChan = []chan interface{}{make(chan interface{}), make(chan interface{})}
	pipline.Start(context.Background())
}

func TestExpr(t *testing.T) {
	con := "type == test"
	conf := config.InputConfig{
		Type: "test",
	}
	fmt.Println(expr(con, conf))

}
