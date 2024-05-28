package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-instagram-user-info-requests-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-instagram-user-info-requests-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-instagram-user-info-requests-rmq-kube/config"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient

	//configure    *existence_conf.ExistenceConf
	//complementer *sub_func_complementer.SubFuncComplementer
}

func NewDPFMAPICaller(
	conf *config.Conf,
	rmq *rabbitmq.RabbitmqClient,
	// confirmor *existence_conf.ExistenceConf,
	// complementer *sub_func_complementer.SubFuncComplementer,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
		//configure:    confirmor,
		//complementer: complementer,
	}
}

func (c *DPFMAPICaller) AsyncRequests(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (interface{}, []error) {
	var instagramUserInfo *[]dpfm_api_output_formatter.InstagramUserInfoResponse
	var errs []error

	for _, fn := range accepter {
		switch fn {
		case "InstagramUserInfo":
			instagramUserInfo = c.InstagramUserInfo(input, &errs, log, c.conf)
		default:
		}
	}

	data := &dpfm_api_output_formatter.Message{
		InstagramUserInfoResponse: instagramUserInfo,
	}

	return data, errs
}
