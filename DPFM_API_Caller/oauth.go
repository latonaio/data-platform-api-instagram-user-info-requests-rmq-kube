package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-instagram-user-info-requests-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-instagram-user-info-requests-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-instagram-user-info-requests-rmq-kube/config"
	"encoding/json"
	"fmt"
	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
)

func (c *DPFMAPICaller) InstagramUserInfo(
	input *dpfm_api_input_reader.SDC,
	errs *[]error,
	log *logger.Logger,
	conf *config.Conf,
) *[]dpfm_api_output_formatter.InstagramUserInfoResponse {
	var instagramUserInfo []dpfm_api_output_formatter.InstagramUserInfoResponse

	accessToken := input.InstagramUserInfo.AccessToken

	userInfoBaseURL := conf.OAuth.UserInfoURL
	userInfoURL := fmt.Sprintf(
		"%s?access_token=%s&fields=id,username",
		userInfoBaseURL,
		accessToken,
	)

	req, err := http.NewRequest("GET", userInfoURL, nil)

	if err != nil {
		*errs = append(*errs, xerrors.Errorf("NewRequest error: %d", err))
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request error: %d", err))
		return nil
	}
	defer resp.Body.Close()

	userInfoBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request response read error: %d", err))
		return nil
	}

	var response map[string]interface{}
	err = json.Unmarshal(userInfoBody, &response)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("Response response error: %d", err))
		return nil
	}

	errorObj, ok := response["error"].(map[string]interface{})
	if ok {
		code, ok := errorObj["code"].(float64)
		if ok {
			errMsg, _ := errorObj["message"].(string)
			*errs = append(*errs, xerrors.Errorf("Status code error: %v %v", code, errMsg))
			return nil
		}
	}

	var instagramUserInfoResponseBody dpfm_api_output_formatter.InstagramUserInfoResponseBody
	err = json.Unmarshal(userInfoBody, &instagramUserInfoResponseBody)
	if err != nil {
		*errs = append(*errs, xerrors.Errorf("User info request response unmarshal error: %d", err))
		return nil
	}

	userInfo := dpfm_api_output_formatter.ConvertToInstagramUserInfoRequestsFromResponse(instagramUserInfoResponseBody)

	instagramUserInfo = append(
		instagramUserInfo,
		userInfo,
	)

	return &instagramUserInfo
}
