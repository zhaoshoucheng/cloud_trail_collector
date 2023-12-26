package collector

import (
	"cloud_trail_collector/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

// ESOutPutCollector 输出到ES
type ESOutPutCollector struct {
	config *config.OutPutConfig
	client *elastic.Client
}

func NewESOutPutCollector(conf *config.OutPutConfig) (OutPutCollector, error) {
	if conf.Index == "" {
		return nil, errors.New("index is empty")
	}
	collector := &ESOutPutCollector{}
	collector.config = conf
	var err error
	collector.client, err = elastic.NewClient(
		elastic.SetBasicAuth(conf.UserName, conf.PassWord),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetURL(conf.EndPoints...),
		//elastic.SetMaxRetries(5),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewConstantBackoff(2*time.Second))),
	)
	if err != nil {
		err = fmt.Errorf("init es cli fail, err: %v", err)
		return nil, err
	}
	return collector, nil
}
func (es *ESOutPutCollector) GetName() string {
	return "ESOutPutCollector"
}

func (es *ESOutPutCollector) Insert(context context.Context, body interface{}) error {
	if es.client == nil {
		return fmt.Errorf("es client is nil")
	}

	indexService := es.client.Index().Index(es.config.Index)
	b, err := json.Marshal(body)
	if err != nil {
		return errors.New(fmt.Sprintf("ESOutPutCollector Body Marshal err %v body %v", err, body))
	}
	_, err = indexService.BodyString(string(b)).Do(context)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	registerOutPutCollector("elastic", NewESOutPutCollector)
}
