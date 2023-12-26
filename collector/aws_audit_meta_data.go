// This file is auto-generated, don't edit it. Thanks.
package collector

import (
	"cloud_trail_collector/config"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"strings"
	"sync/atomic"
	"time"
)

type AwsCloudTrailInputCollector struct {
	svc       *cloudtrail.CloudTrail
	conf      *config.InputConfig
	count     int64
	nextToken string
}

func NewAwsCloudTrailInputCollector(input *config.InputConfig) (InputCollector, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(input.RegionId),
		Credentials: credentials.NewStaticCredentials(input.AccessKey, input.SecretKey, ""),
	}))
	svc := cloudtrail.New(sess)

	awsCollector := &AwsCloudTrailInputCollector{
		svc:  svc,
		conf: input,
	}
	return awsCollector, nil
}
func init() {
	registerInputCollector("aws_cloud_trail", NewAwsCloudTrailInputCollector)
}
func (a *AwsCloudTrailInputCollector) GetName() string {
	return "AwsCloudTrailInputCollector"
}
func (a *AwsCloudTrailInputCollector) Update(ctx context.Context, ch chan interface{}) {
	if ch == nil {
		panic("AwsCloudTrailInputCollector result ch init nil ")
	}
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Aws云日志采集执行时间:", time.Now())
			a.count = 0
			endTime := time.Now()
			startTime := time.Now().Add(time.Hour * -1)
			a.GetAwsCloudtrailData(startTime, endTime, ch)
			fmt.Printf("add aws cloud trail data into es success, data's length is: %d\n", atomic.LoadInt64(&a.count))
		}
	}
}
func (a *AwsCloudTrailInputCollector) GetAwsCloudtrailData(startTime time.Time, endTime time.Time, result chan interface{}) {
	input := &cloudtrail.LookupEventsInput{StartTime: aws.Time(startTime), EndTime: aws.Time(endTime)}
	if a.nextToken != "" {
		input.NextToken = &a.nextToken
	}
	resp, err := a.svc.LookupEvents(input)
	if err != nil {
		fmt.Println("AwsCloudTrailInputCollector LookupEvents err ", err)
		return
	}
	for _, event := range resp.Events {
		data := transferAwsEvent(event)
		result <- data
	}
	atomic.AddInt64(&a.count, int64(len(resp.Events)))
	if resp.NextToken != nil && *resp.NextToken != "" {
		//翻页
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			a.nextToken = *resp.NextToken
			a.GetAwsCloudtrailData(startTime, endTime, result)
		}
	} else {
		a.nextToken = ""
	}

}

type Event struct {
	// The Amazon Web Services access key ID that was used to sign the request.
	// If the request was made with temporary security credentials, this is the
	// access key ID of the temporary credentials.
	AccessKeyId *string `type:"string"`

	// A JSON string that contains a representation of the event returned.
	CloudTrailEvent *string `type:"string"`

	// The CloudTrail ID of the event returned.
	EventId *string `type:"string"`

	// The name of the event returned.
	EventName *string `type:"string"`

	// The Amazon Web Services service to which the request was made.
	EventSource *string `type:"string"`

	// The date and time of the event returned.
	EventTime *time.Time `type:"timestamp"`

	// Information about whether the event is a write event or a read event.
	ReadOnly *string `type:"string"`
	//资源类型 "，"分隔
	ResourcesType string `type:"string"`
	//资源名称 "，"分隔
	ResourcesName string `type:"string"`
	// A user name or role name of the requester that called the API in the event
	// returned.
	Username *string `type:"string"`
}

func transferAwsEvent(event *cloudtrail.Event) *Event {
	res := &Event{}
	res.AccessKeyId = event.AccessKeyId
	res.CloudTrailEvent = event.CloudTrailEvent
	res.EventId = event.EventId
	res.EventName = event.EventName
	res.EventSource = event.EventSource
	res.EventTime = event.EventTime
	res.ReadOnly = event.ReadOnly
	var resourceTypes []string
	var resourceName []string
	for _, resource := range event.Resources {
		tp := "-"
		if resource.ResourceType != nil {
			tp = *resource.ResourceType
		}
		resourceTypes = append(resourceTypes, tp)

		name := "-"
		if resource.ResourceName != nil {
			name = *resource.ResourceName
		}
		resourceName = append(resourceName, name)
	}
	res.ResourcesType = strings.Join(resourceTypes, ",")
	res.ResourcesType = strings.Join(resourceTypes, ",")
	res.Username = event.Username
	return res
}
