// This file is auto-generated, don't edit it. Thanks.
package collector

import (
	"cloud_trail_collector/config"
	"context"
	"encoding/json"
	"fmt"
	actiontrail20200706 "github.com/alibabacloud-go/actiontrail-20200706/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"reflect"
	"sync/atomic"
	"time"
)

type AliCloudTrailInputCollector struct {
	client    *actiontrail20200706.Client
	conf      *config.InputConfig
	count     int64
	nextToken string
}

func (ali *AliCloudTrailInputCollector) GetName() string {
	return "AliCloudTrailInputCollector"
}
func (ali *AliCloudTrailInputCollector) Update(ctx context.Context, ch chan interface{}) {
	if ch == nil {
		panic("AliCloudTrailInputCollector result ch init nil ")
	}
	//**Test**
	//fmt.Println("Ali云日志采集执行时间:", time.Now())
	//ali.count = 0
	//endTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	//startTime := time.Now().UTC().Add(time.Hour * -1).Format("2006-01-02T15:04:05Z")
	//ali.GetAliCloudtrailData(startTime, endTime, ch)
	//fmt.Printf("add ali cloud trail data into es success, data's length is: %d\n", atomic.LoadInt64(&ali.count))
	//return
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Aws云日志采集执行时间:", time.Now())
			ali.count = 0
			endTime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
			startTime := time.Now().UTC().Add(time.Hour * -1).Format("2006-01-02T15:04:05Z")
			ali.GetAliCloudtrailData(startTime, endTime, ch)
			fmt.Printf("add aws cloud trail data into es success, data's length is: %d\n", atomic.LoadInt64(&ali.count))
		}
	}
}

func (ali *AliCloudTrailInputCollector) GetAliCloudtrailData(startTime string, endTime string, result chan interface{}) {
	maxResult := "50"
	lookupEventsRequest := &actiontrail20200706.LookupEventsRequest{
		StartTime:  &startTime,
		EndTime:    &endTime,
		MaxResults: &maxResult,
	}
	if ali.nextToken != "" {
		lookupEventsRequest.NextToken = &ali.nextToken
	}
	resp, _err := ali.client.LookupEventsWithOptions(lookupEventsRequest, &util.RuntimeOptions{})
	if _err != nil {
		fmt.Println("AliCloudTrailInputCollector GetAliCloudtrailData LookupEventsWithOptions err ", _err)
		return
	}
	if resp.StatusCode == nil || *resp.StatusCode != 200 {
		fmt.Println("status code is not 200")
		return
	}
	for _, event := range resp.Body.Events {
		res := transferAliEvent(event)
		result <- res
	}
	if resp.Body.NextToken != nil && *resp.Body.NextToken != "" {
		//翻页
		fmt.Println(*resp.Body.NextToken)
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			ali.nextToken = *resp.Body.NextToken
			ali.GetAliCloudtrailData(startTime, endTime, result)
		}
	} else {
		ali.nextToken = ""
	}
}

func NewAliCloudTrailInputCollector(input *config.InputConfig) (InputCollector, error) {
	conf := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &input.AccessKey,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &input.SecretKey,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Actiontrail
	conf.Endpoint = tea.String(input.EndPoint)
	client := &actiontrail20200706.Client{}
	client, err := actiontrail20200706.NewClient(conf)
	if err != nil {
		return nil, err
	}
	collector := &AliCloudTrailInputCollector{
		client: client,
		conf:   input,
	}
	return collector, nil
}
func init() {
	registerInputCollector("ali_cloud_trail", NewAliCloudTrailInputCollector)
}

type AliEvent struct {
	/*阿里云地域*/
	AcsRegion string `json:"acsRegion"`
	/*事件的补充信息 map 类型，进行json化处理成字符串*/
	AdditionalEventData string `json:"additionalEventData"`
	/*当eventType取值为ApiCall时，事件代表一个API的调用。此时，该字段为API的版本信息。*/
	ApiVersion string `json:"apiVersion"`
	/*事件分类。取值：Management（管控事件）*/
	EventCategory string `json:"eventCategory"`
	/*事件ID*/
	EventId string `json:"eventId"`
	/*事件名称。具体含义如下：

	如果eventType取值为ApiCall，该字段为API的名称。

	如果eventType取值不为ApiCall，该字段表示事件含义。*/
	EventName string `json:"eventName"`
	/*
		事件的读写类型。取值：

		Write：写类型。

		Read：读类型。*/
	EventRW string `json:"eventRW"`
	/*事件来源*/
	EventSource string `json:"eventSource"`
	/*事件的发生时间（UTC格式）*/
	EventTime string `json:"eventTime"`
	/*
	   发生的事件类型。取值：

	   ApiCall：API调用事件。

	   ConsoleOperation：部分控制台或售卖页的管控事件。

	   ConsoleSignin：控制台登录事件。

	   ConsoleSignout：控制台登出事件。

	   AliyunServiceEvent：此类事件为阿里云平台对您的资源执行的管控事件，目前主要是预付费实例的到期自动释放事件。*/
	EventType string `json:"eventType"`
	/*管控事件格式的版本，当前版本为1*/
	EventVersion string `json:"eventVersion"`
	/*云服务处理API请求发生错误时，记录的错误码。*/
	ErrorCode string `json:"errorCode"`
	/*云服务处理API请求发生错误时，记录的错误消息。*/
	ErrorMessage string `json:"errorMessage"`
	/*请求ID。*/
	RequestId string `json:"requestId"`
	/*API请求的输入参数。 json格式化*/
	RequestParameters string `json:"requestParameters"`
	/*事件的相关资源名称，是资源的唯一标识。
	同类型的资源名称（ID）之间以半角逗号（,）间隔，不同类型的资源名称（ID）之间以半角分号（;）间隔。*/
	ResourceName string `json:"resourceName"`
	/*事件的相关资源类型。
	多个资源类型之间以半角分号（;）间隔。*/
	ResourceType string `json:"resourceType"`
	/*API响应的数据。*/
	ResponseElements string `json:"responseElements"`
	/*事件影响的资源列表。*/
	ReferencedResources string `json:"referencedResources"`
	/*事件相关的阿里云服务名称。*/
	ServiceName string `json:"serviceName"`
	/*事件发起的源IP地址。*/
	SourceIpAddress string `json:"sourceIpAddress"`
	/*发送API请求的客户端代理标识。*/
	UserAgent string `json:"userAgent"`
	/*
		是否全局事件。取值：

		true：全局事件。

		false：非全局事件*/
	IsGlobal string `json:"isGlobal"`
	/*事件属性*/
	EventAttributes string `json:"eventAttributes"`
	/*请求者的身份信息*/
	UserIdentity string `json:"userIdentity"`
}

func transferAliEvent(event map[string]interface{}) *AliEvent {
	myEvent := AliEvent{}
	vo := reflect.ValueOf(&myEvent).Elem()
	for i := 0; i < reflect.TypeOf(myEvent).NumField(); i++ {
		field := vo.Type().Field(i).Tag.Get("json")
		obj, exists := event[field]
		if !exists {
			continue
		}
		name := vo.Type().Field(i).Name
		switch obj.(type) {
		case string:
			vo.FieldByName(name).Set(reflect.ValueOf(obj))
		default:
			str, err := json.Marshal(obj)
			if err != nil {
				continue
			}
			vo.FieldByName(name).Set(reflect.ValueOf(string(str)))
		}
	}
	return &myEvent
}
