package collector

import (
	"cloud_trail_collector/config"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	inputFactories  = make(map[string]func(config *config.InputConfig) (InputCollector, error))
	outputFactories = make(map[string]func(config *config.OutPutConfig) (OutPutCollector, error))
)

type InputCollector interface {
	GetName() string
	Update(context context.Context, ch chan interface{})
}
type OutPutCollector interface {
	GetName() string
	Insert(context context.Context, body interface{}) error
}

func registerInputCollector(collector string, factory func(config *config.InputConfig) (InputCollector, error)) {
	inputFactories[collector] = factory
}
func NewInputCollector(cfg *config.InputConfig) (InputCollector, error) {
	if cfg.Type == "" {
		return nil, errors.New("input config type is empty")
	}
	factory, exists := inputFactories[cfg.Type]
	if !exists {
		return nil, errors.New(fmt.Sprintf("not support input collector %s", cfg.Type))
	}
	return factory(cfg)
}
func registerOutPutCollector(collector string, factory func(config *config.OutPutConfig) (OutPutCollector, error)) {
	outputFactories[collector] = factory
}
func NewOutputCollector(cfg *config.OutPutConfig) (OutPutCollector, error) {
	if cfg.Type == "" {
		return nil, errors.New("output config type is empty")
	}
	factory, exists := outputFactories[cfg.Type]
	if !exists {
		return nil, errors.New(fmt.Sprintf("not support output collector %s", cfg.Type))
	}
	return factory(cfg)
}

type Pipeline struct {
	Input       InputCollector
	InputChan   chan interface{}
	OutPuts     []OutPutCollector
	OutPutsChan []chan interface{}
}

var Pipelines []*Pipeline

func MakePipeLines() error {
	cfg := config.GetConfig()
	if cfg == nil {
		return errors.New("config not init ")
	}
	if len(cfg.Inputs) == 0 || len(cfg.Outputs) == 0 {
		return errors.New("input config is nil or output config is nil")
	}
	for _, inputCfg := range cfg.Inputs {
		inputCollector, err := NewInputCollector(inputCfg)
		if err != nil {
			continue
		}
		fmt.Printf("init %s input success ~ \n", inputCollector.GetName())
		pip := &Pipeline{}
		pip.InputChan = make(chan interface{})
		for _, outputCfg := range cfg.Outputs {
			if outputCfg.Condition != "" {
				match := expr(outputCfg.Condition, *inputCfg)
				//条件不匹配
				if !match {
					continue
				}
			}
			outputCollector, err := NewOutputCollector(outputCfg)
			if err != nil {
				continue
			}
			fmt.Printf("init %s output success ~ \n", outputCollector.GetName())
			pip.OutPuts = append(pip.OutPuts, outputCollector)
			pip.OutPutsChan = append(pip.OutPutsChan, make(chan interface{}))
		}
		if len(pip.OutPuts) == 0 {
			continue
		}
		pip.Input = inputCollector
		Pipelines = append(Pipelines, pip)
	}
	if len(Pipelines) == 0 {
		return errors.New("pipline init fail, input or output is nil")
	}
	return nil
}
func StartAllPipelines(ctx context.Context) error {
	if len(Pipelines) == 0 {
		return errors.New("there is no pipelines ")
	}
	for _, pip := range Pipelines {
		go pip.Start(ctx)
	}
	fmt.Println("pipelines start success ~")
	return nil
}

// Start 管道开始采集，阻塞，协程调用
func (pip *Pipeline) Start(ctx context.Context) {
	if len(pip.OutPutsChan) != len(pip.OutPuts) {
		panic(fmt.Sprintf("Pipeline init err ,%d , %d", len(pip.OutPuts), len(pip.OutPuts)))
	}
	pip.handleMessage(ctx)
	go pip.Input.Update(ctx, pip.InputChan)
	for {
		select {
		case msg := <-pip.InputChan:
			//广播，一条消息广播到所有OutPut通道
			for index, OutPutChan := range pip.OutPutsChan {
				timeOut := time.NewTimer(time.Second)
				ch := OutPutChan
				select {
				case ch <- msg:
					//fmt.Println(msg) //debug
				case <-timeOut.C: //监控超时
					fmt.Println("Broadcast timeout type: input: ", pip.Input.GetName(), " output:", pip.OutPuts[index].GetName())
				}
			}
		}
	}
}

func (pip *Pipeline) handleMessage(ctx context.Context) {
	for index, outPut := range pip.OutPuts {
		outPutHeader(ctx, pip.OutPutsChan[index], outPut)
	}
}

func outPutHeader(ctx context.Context, ch chan interface{}, collector OutPutCollector) {
	for i := 0; i < config.GetConfig().Worker; i++ {
		go func() {
			for {
				select {
				case msg := <-ch:
					err := collector.Insert(ctx, msg)
					if err != nil {
						fmt.Println("outPutHeader err ", err)
					}
				}
			}
		}()
	}
}

func expr(condition string, data interface{}) bool {
	args := strings.Split(strings.ReplaceAll(condition, " ", ""), "==")
	if len(args) != 2 {
		return false
	}
	tag := args[0]
	value := args[1]
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)
		if field.Tag.Get("toml") == tag && v.Field(index).String() == value {
			return true
		}
	}
	return false
}
