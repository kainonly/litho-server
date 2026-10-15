package tencent

import (
	"context"

	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func (x *Service) Sms(ctx context.Context, tid string, params []string, phone []string) (err error) {
	credential := tcommon.NewCredential(
		x.V.Sms.SecretId,
		x.V.Sms.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = new(x.V.Sms.AppId)
	request.SignName = new(x.V.Sms.Sign)
	request.TemplateId = new(tid)
	request.TemplateParamSet = tcommon.StringPtrs(params)
	request.PhoneNumberSet = tcommon.StringPtrs(phone)
	request.SetContext(ctx)
	if _, err = client.SendSms(request); err != nil {
		return
	}
	return
}
