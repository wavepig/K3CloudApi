package sdk

import (
	"encoding/json"
)

type K3CloudApiBOS string

const (
	GetDataCenterList K3CloudApiBOS = "Kingdee.BOS.ServiceFacade.ServicesStub.Account.AccountService.GetDataCenterList"
	// 操作接口
	ExcuteOperation K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.ExcuteOperation"
	// 保存
	Save K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Save"
	// 批量保存
	BatchSave K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.BatchSave"
	// 审核
	Audit K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Audit"
	// 反审核
	UnAudit K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.UnAudit"
	// 查看
	View K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.View"
	// 提交
	Submit K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Submit"
	// 删除
	Delete K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Delete"
	// 单据查询
	ExecuteBillQuery K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.ExecuteBillQuery"
	// 单据查询(json)
	BillQuery K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.BillQuery"
	// 暂存
	Draft K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Draft"
	// 分配
	Allocate K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Allocate"
	// 弹性域保存
	FlexSave K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.FlexSave"
	// 发送消息
	SendMsg K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.SendMsg"
	// 下推
	Push K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Push"
	// 分组保存
	GroupSave K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.GroupSave"
	// 拆单
	UnPack K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.Disassembly"
	// 查询单据接口
	QueryBusinessInfo_ K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.QueryBusinessInfo"
	// 查询分组信息接口
	QueryGroupInfo K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.QueryGroupInfo"
	// 工作流审批
	Approve K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.WorkflowAudit"
	// 分组删除接口
	GroupDelete K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.GroupDelete"
	// 切换组织
	SwitchOrg K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.SwitchOrg"
	// 取消分配
	CancelAllocate K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.CancelAllocate"
	// 撤销服务接口
	CancelAssign K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.CancelAssign"
	// 获取报表数据
	GetSysReportData K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.GetSysReportData"
	// 上传附件
	AttachmentUpload K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.AttachmentUpload"
	// 下载附件
	AttachmentDownLoad K3CloudApiBOS = "Kingdee.BOS.WebApi.ServicesStub.DynamicFormService.AttachmentDownLoad"

	URL = ".common.kdsvc"
)

type K3CloudApiSdk struct {
	client *client
}

func NewK3CloudApiSdk(at AuthType, acctID, username, password, serverUrl, appid, appSecret string) *K3CloudApiSdk {
	c := newClient(InitConfig(at, acctID, username, password, serverUrl, appid, appSecret))

	return &K3CloudApiSdk{
		client: c,
	}
}

func (c *K3CloudApiSdk) Request(bos K3CloudApiBOS, formId string, data map[string]string) ([]byte, error) {
	postData := map[string]interface{}{
		"formid": formId,
		"data":   data,
	}
	reqData, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	return c.client.request("POST", string(bos)+URL, reqData)
}

func (c *K3CloudApiSdk) RequestAny(bos K3CloudApiBOS, formId string, data map[string]any) ([]byte, error) {
	postData := map[string]interface{}{
		"formid": formId,
		"data":   data,
	}
	reqData, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}
	return c.client.request("POST", string(bos)+URL, reqData)
}

func (c *K3CloudApiSdk) RequestByBos(bos string, data map[string]any) ([]byte, error) {
	reqData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.client.request("POST", string(bos)+URL, reqData)
}
