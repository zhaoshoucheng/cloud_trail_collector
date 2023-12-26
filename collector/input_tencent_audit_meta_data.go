package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"github.com/spf13/cast"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"sync/atomic"
	"time"
)

type TCMetaDataInputCollector struct {
	client *cloudaudit.Client
	marker *uint64
	count  int64
}

func NewTCMetaDataInputCollector(input *config.InputConfig) (InputCollector, error) {
	credential := common.NewCredential(
		input.AccessKey,
		input.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = input.EndPoint
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cloudaudit.NewClient(credential, input.RegionId, cpf)
	tc := &TCMetaDataInputCollector{
		client: client,
		marker: nil,
		count:  0,
	}
	return tc, nil
}

func init() {
	registerInputCollector("tc_audit_meta_data", NewTCMetaDataInputCollector)
}
func (tc *TCMetaDataInputCollector) GetName() string {
	return "TCMetaDataInputCollector"
}
func (tc *TCMetaDataInputCollector) Update(ctx context.Context, ch chan interface{}) {
	if ch == nil {
		panic("TCMetaDataInputCollector result ch init nil ")
	}
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			fmt.Println("腾讯云日志采集执行时间:", time.Now())
			tc.count = 0
			endTime := uint64(time.Now().UnixMilli())
			startTime := uint64(time.Now().Add(time.Hour * -1).UnixMilli())
			tc.GetTCAuditMetaData(true, &startTime, &endTime, ch)
			fmt.Printf("add tencent cloud audit meta data into es success, data's length is: %d\n", atomic.LoadInt64(&tc.count))
		}
	}
}

func (tc *TCMetaDataInputCollector) GetTCAuditMetaData(first bool, startTime, endTime *uint64, resultChan chan interface{}) {
	if !first && tc.marker == nil {
		return
	}

	request := cloudaudit.NewDescribeEventsRequest()
	request.MaxResults = common.Uint64Ptr(50) //单次请求 返回日志数量最大值50
	request.StartTime = startTime
	request.EndTime = endTime
	request.NextToken = tc.marker
	response, err := tc.client.DescribeEvents(request)
	if err == nil {
		//fmt.Printf("%+v\n", response)
	} else {
		//重试一次
		if response, err = tc.client.DescribeEvents(request); err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, item := range response.Response.Events {
		traces := genMyEvent(item)
		resultChan <- traces
	}

	atomic.AddInt64(&tc.count, int64(len(response.Response.Events)))

	if !*response.Response.ListOver {
		//翻页频率控制
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			tc.marker = response.Response.NextToken
			tc.GetTCAuditMetaData(false, startTime, endTime, resultChan)
		}
	} else {
		tc.marker = nil
	}
}

type MyEvent struct {
	// 日志ID
	EventId string `json:"event_id,omitempty"`

	// 用户名
	Username string `json:"username,omitempty"`

	// 事件时间
	EventTime time.Time `json:"event_time,omitempty"`

	RecordTime string `json:"record_time"`

	// 日志详情
	CloudAuditEvent string `json:"cloud_audit_event"`

	// 资源类型中文描述（此字段请按需使用，如果您是其他语言使用者，可以忽略该字段描述）
	ResourceTypeCn string `json:"resource_type_cn,omitempty"`

	// 鉴权错误码
	ErrorCode int64 `json:"error_code,omitempty"`

	// 事件名称
	EventName string `json:"event_name,omitempty"`

	// 证书ID
	// 注意：此字段可能返回 null，表示取不到有效值。
	SecretId string `json:"secret_id,omitempty"`

	// 请求来源
	EventSource string `json:"event_source,omitempty"`

	// 请求ID
	RequestID string `json:"request_id,omitempty"`

	// 资源地域
	ResourceRegion string `json:"resource_region,omitempty"`

	// 主账号ID
	AccountID int64 `json:"account_id,omitempty"`

	// 源IP
	// 注意：此字段可能返回 null，表示取不到有效值。
	SourceIPAddress string `json:"source_ip_address,omitempty"`

	// 事件名称中文描述（此字段请按需使用，如果您是其他语言使用者，可以忽略该字段描述）
	EventNameCn string `json:"event_name_cn,omitempty"`

	// 资源对
	//Resources *Resource `json:"Resources,omitempty" name:"Resources"`
	// 资源类型
	ResourceType string `json:"resource_type,omitempty"`

	// 资源名称
	// 注意：此字段可能返回 null，表示取不到有效值。
	ResourceName string `json:"resource_name,omitempty"`
	// 事件地域
	EventRegion string `json:"event_region,omitempty"`

	// IP 归属地
	Location string `json:"location,omitempty"`
}

func genMyEvent(msg *cloudaudit.Event) *MyEvent {
	myEvent := &MyEvent{}
	if msg.EventId != nil {
		myEvent.EventId = *msg.EventId
	}
	if msg.Username != nil {
		myEvent.Username = *msg.Username
	}
	if msg.EventTime != nil {
		timestamp := cast.ToInt64(*msg.EventTime)
		myEvent.ParseTime(timestamp)
	}
	if msg.CloudAuditEvent != nil {
		myEvent.CloudAuditEvent = *msg.CloudAuditEvent
	}
	if msg.ResourceTypeCn != nil {
		myEvent.ResourceTypeCn = *msg.ResourceTypeCn
	}
	if msg.ErrorCode != nil {
		myEvent.ErrorCode = *msg.ErrorCode
	}
	if msg.EventName != nil {
		myEvent.EventName = *msg.EventName
	}

	if msg.SecretId != nil {
		myEvent.SecretId = *msg.SecretId
	}
	if msg.EventSource != nil {
		myEvent.EventSource = *msg.EventSource
	}
	if msg.RequestID != nil {
		myEvent.RequestID = *msg.RequestID
	}
	if msg.ResourceRegion != nil {
		myEvent.ResourceRegion = *msg.ResourceRegion
	}
	if msg.AccountID != nil {
		myEvent.AccountID = *msg.AccountID
	}
	if msg.SourceIPAddress != nil {
		myEvent.SourceIPAddress = *msg.SourceIPAddress
	}
	if msg.EventNameCn != nil {
		myEvent.EventNameCn = *msg.EventNameCn
	}
	if msg.Resources != nil {
		myEvent.ResourceType = *msg.Resources.ResourceType
		myEvent.ResourceName = *msg.Resources.ResourceName
	}

	if msg.EventRegion != nil {
		myEvent.EventRegion = *msg.EventRegion
	}
	if msg.Location != nil {
		myEvent.Location = *msg.Location
	}
	return myEvent
}

func (event *MyEvent) ParseTime(timestamp int64) {
	event.EventTime = time.Unix(timestamp, 0)
	event.RecordTime = time.Unix(timestamp, 0).Format(Format)
	return
}
