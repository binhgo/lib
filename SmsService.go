package lib

import (
	"github.com/spf13/viper"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"fmt"
	"lib/util"
)

type SmsService struct {
	smsUrl string
}

func (this *SmsService) init() {
	// 修改为您的apikey(https://www.yunpian.com)登录官网后获取
	this.smsUrl = viper.GetString("sms_url")

}

/**
	// 修改为您要发送的手机号码，多个号码用逗号隔开
 */
func (this *SmsService) TplBatchSend(mobile string, tplId int, tplValues url.Values) {

	if this.smsUrl == "" {
		this.init()
	}

	this.httpsPostForm(this.smsUrl, mobile, tplId, tplValues)
}

func (this *SmsService) httpsPostForm(url string, mobile string, tplId int, data url.Values) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {

				Log.Error(e.Error())
			}

			Log.Error(util.GetStackTrace(r))
		}
	}()

	type PushParams struct {
		BusinessId    int    `json:"businessId"`
		BusinessName  string `json:"businessName"`
		Target        int    `json:"target"`
		ExtendCode    string `json:"extendCode"`
		Mobile        string `json:"mobile"`
		ChannelId     int    `json:"channelId"`
		TemplateId    int    `json:"templateId"`
		TemplateParam string `json:"templateParam"`
		OperatorID    int    `json:"operatorID"`
		OperatorName  string `json:"operatorName"`
		IsDelayed     int    `json:"isDelayed"`
		UserType      int    `json:"userType"`
	}

	pushParams := PushParams{}
	pushParams.BusinessId = 60000
	pushParams.ChannelId = 1000
	pushParams.BusinessName = "派单系统"
	pushParams.Target = 1
	pushParams.TemplateId = tplId
	pushParams.ExtendCode = data.Get("templateId")
	pushParams.Mobile = mobile
	pushParams.OperatorID = 0
	pushParams.OperatorName = "派单系统"
	pushParams.IsDelayed = 0
	pushParams.UserType = 5

	params := map[string]interface{}{}
	for k, v := range data {
		params[k] = v[0]
	}

	bytesNotify, err := json.Marshal(params)
	pushParams.TemplateParam = string(bytesNotify)

	pushParamsVal, _ := json.Marshal(pushParams)

	Log.Infof("url:%v,params:%v", url, string(pushParamsVal))
	//下发APP通知
	req, err := http.NewRequest("POST", viper.GetString("sms_url"), bytes.NewBuffer(pushParamsVal))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
