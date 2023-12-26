package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/region"
	"github.com/spf13/cast"
	"sync/atomic"
	"time"
)

const (
	Format = "2006-01-02 15:04:05"
)

type HWMetaDataInputCollector struct {
	client *cts.CtsClient
	marker *string
	count  int64
	conf   *config.InputConfig
}

func NewHwInputCollector(input *config.InputConfig) (InputCollector, error) {
	auth := basic.NewCredentialsBuilder().
		WithProjectId(input.EndPoint).
		WithAk(input.AccessKey).
		WithSk(input.SecretKey).
		Build()
	hwClient := cts.NewCtsClient(
		cts.CtsClientBuilder().
			WithRegion(region.ValueOf(input.RegionId)).
			WithCredential(auth).
			Build())
	hw := &HWMetaDataInputCollector{
		client: hwClient,
		marker: nil,
		count:  0,
		conf:   input,
	}
	return hw, nil
}

func (hw *HWMetaDataInputCollector) GetName() string {
	return "TimerInputCollector"
}
func (hw *HWMetaDataInputCollector) Update(ctx context.Context, ch chan interface{}) {
	if ch == nil {
		panic("HWMetaDataInputCollector result ch init nil ")
	}
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			fmt.Println("华为云日志采集执行时间:", time.Now())
			hw.count = 0
			//hw.StartGetMetadata(ch)
			endTime := time.Now()
			startTime := time.Now().Add(time.Hour * -1)
			hw.GetHWAuditMetaData(true, startTime.UnixMilli(), endTime.UnixMilli(), ch)
			fmt.Printf("add hua wei cloud audit meta data into es success, data's length is: %d\n", atomic.LoadInt64(&hw.count))
		}
	}
}

func init() {
	registerInputCollector("hw_audit_meta_data", NewHwInputCollector)
}

func (hw *HWMetaDataInputCollector) GetHWAuditMetaData(first bool, startTime, endTime int64, resultChan chan interface{}) {
	if !first && hw.marker == nil {
		return
	}

	limit := int32(200)
	request := &model.ListTracesRequest{
		TraceType: model.GetListTracesRequestTraceTypeEnum().SYSTEM,
		From:      &startTime,
		To:        &endTime,
		Limit:     &limit,
		Next:      hw.marker,
	}
	response, err := hw.client.ListTraces(request)
	if err == nil {
		//fmt.Printf("%+v\n", response)
	} else {
		fmt.Println(err)
	}

	for _, item := range *response.Traces {
		traces := genMyTraces(&item)
		resultChan <- traces
	}

	atomic.AddInt64(&hw.count, int64(len(*response.Traces)))

	if response.MetaData.Marker != nil {
		//翻页频率控制
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			hw.marker = response.MetaData.Marker
			hw.GetHWAuditMetaData(false, startTime, endTime, resultChan)
		}
	} else {
		hw.marker = nil
	}
}

type MyTraces struct {

	// 标识事件对应的云服务资源ID。
	ResourceId string `json:"resource_id,omitempty"`

	// 标识查询事件列表对应的事件名称。由0-9,a-z,A-Z,'-','.','_',组成，长度为1～64个字符，且以首字符必须为字母。
	TraceName string `json:"trace_name,omitempty"`

	// 标识事件等级，目前有三种：正常（normal），警告（warning），事故（incident）。
	TraceRating string `json:"trace_rating,omitempty"`

	// 标识事件发生源头类型，管理类事件主要包括API调用（ApiCall），Console页面调用（ConsoleAction）和系统间调用（SystemAction）。 数据类事件主要包括ObsSDK，ObsAPI。
	TraceType string `json:"trace_type,omitempty"`

	// 标识事件对应接口请求内容，即资源操作请求体。
	Request string `json:"request,omitempty"`

	// 记录用户请求的响应，标识事件对应接口响应内容，即资源操作结果返回体。
	Response string `json:"response,omitempty"`

	// 记录用户请求的响应，标识事件对应接口返回的HTTP状态码。
	Code int `json:"code,omitempty"`

	// 标识事件对应的云服务接口版本。
	ApiVersion string `json:"api_version,omitempty"`

	// 标识其他云服务为此条事件添加的备注信息。
	Message string `json:"message,omitempty"`

	// 标识事件的ID，由系统生成的UUID。
	TraceId string `json:"trace_id,omitempty"`

	User string `json:"user,omitempty"`

	// 标识查询事件列表对应的云服务类型。必须为已对接CTS的云服务的英文缩写，且服务类型一般为大写字母。
	ServiceType string `json:"service_type,omitempty"`

	// 查询事件列表对应的资源类型。
	ResourceType string `json:"resource_type,omitempty"`

	// 标识触发事件的租户IP。
	SourceIp string `json:"source_ip,omitempty"`

	// 标识事件对应的资源名称。
	ResourceName string `json:"resource_name,omitempty"`

	// 记录本次请求的request id
	RequestId string `json:"request_id,omitempty"`

	// 记录本次请求出错后，问题定位所需要的辅助信息。
	LocationInfo string `json:"location_info,omitempty"`

	// 云资源的详情页面
	Endpoint string `json:"endpoint,omitempty"`

	// 云资源的详情页面的访问链接（不含endpoint）
	ResourceUrl string `json:"resource_url,omitempty"`

	RecordTime string `json:"record_time"`

	Time time.Time `json:"time"`
}

func genMyTraces(trace *model.Traces) *MyTraces {
	itemTrace := &MyTraces{}
	if trace.ResourceId != nil {
		itemTrace.ResourceId = *trace.ResourceId
	}
	if trace.TraceName != nil {
		itemTrace.TraceName = *trace.TraceName
	}
	if trace.TraceRating != nil {
		itemTrace.TraceRating = trace.TraceRating.Value()
	}
	if trace.TraceType != nil {
		itemTrace.TraceType = *trace.TraceType
	}
	if trace.Request != nil {
		itemTrace.Request = *trace.Request
	}
	if trace.Response != nil {
		itemTrace.Response = *trace.Response
	}
	if trace.Code != nil {
		itemTrace.Code = cast.ToInt(*trace.Code)
	}
	if trace.ApiVersion != nil {
		itemTrace.ApiVersion = *trace.ApiVersion
	}
	if trace.Message != nil {
		itemTrace.Message = *trace.Message
	}
	if trace.TraceId != nil {
		itemTrace.TraceId = *trace.TraceId
	}
	if trace.User != nil {
		itemTrace.User = trace.User.String()
	}
	if trace.ServiceType != nil {
		itemTrace.ServiceType = *trace.ServiceType
	}
	if trace.ResourceType != nil {
		itemTrace.ResourceType = *trace.ResourceType
	}
	if trace.ResourceName != nil {
		itemTrace.ResourceName = *trace.ResourceName
	}
	if trace.SourceIp != nil {
		itemTrace.SourceIp = *trace.SourceIp
	}
	if trace.RequestId != nil {
		itemTrace.RequestId = *trace.RequestId
	}
	if trace.LocationInfo != nil {
		itemTrace.LocationInfo = *trace.LocationInfo
	}
	if trace.Endpoint != nil {
		itemTrace.Endpoint = *trace.Endpoint
	}
	if trace.ResourceUrl != nil {
		itemTrace.ResourceUrl = *trace.ResourceUrl
	}

	itemTrace.RecordTime = itemTrace.ParseTime(*trace.RecordTime)

	itemTrace.Time = time.Unix(*trace.Time/1000, 0)

	return itemTrace
}

func (traces *MyTraces) ParseTime(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).Format(Format)
}
