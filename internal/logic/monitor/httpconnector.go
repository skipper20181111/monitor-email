package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"monitor/internal/svc"
	"monitor/internal/types"
	"net/http"
	"strings"
	"time"
)

type HttpConnectorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHttpConnectorLogicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpConnectorLogic {
	return &HttpConnectorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func GetCacheSession(SystemInfo *types.System) string {
	passwd, _ := Decrypt(SystemInfo.Passwd, svc.Keystr)
	return GetSession(SystemInfo.OuterIP, SystemInfo.User, passwd)
}
func GetSession(ip, username, userpasswd string) string {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	cdp_session := ""
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 阻止重定向
			return errors.New("禁止重定向")
		},
	}
	GetSessionUrl := fmt.Sprintf("http://%s:7180/j_spring_security_check?j_username=%s&j_password=%s", ip, username, userpasswd)
	req, _ := http.NewRequest("POST", GetSessionUrl, nil)
	noredirectresp, _ := c.Do(req)
	splited := strings.Split(noredirectresp.Header.Get("Set-Cookie"), "; ")
	for _, setcookies := range splited {
		if strings.Contains(setcookies, "SESSION") {
			cdp_session = setcookies
		}
	}
	noredirectresp.Body.Close()
	PostLoginUrl := fmt.Sprintf("http://%s:7180/cmf/postLogin", ip)
	sreq, _ := http.NewRequest("GET", PostLoginUrl, nil)
	sreq.Header.Add("Cookie", cdp_session)
	sresp, _ := c.Do(sreq)
	sresp.Body.Close()
	return cdp_session
}
func (l *HttpConnectorLogic) Report2Shrcb(shrcbMonitorRes *types.ShrcbMonitorRes) bool {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	shrcbresp := &types.ShrcbMonitorResp{}
	marshal, _ := json.Marshal(shrcbMonitorRes)
	resp, _ := http.Post(l.svcCtx.Config.ServerInfo.Url+shrcbMonitorRes.SysNameEn, "application/json", bytes.NewReader(marshal))
	//resp, _ := http.Post("http://localhost:58888/monitor/shrcbtest?source="+shrcbMonitorRes.SysNameEn, "application/json", bytes.NewReader(marshal))
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, shrcbresp)
	resp.Body.Close()
	fmt.Println("#####################################################################################################", shrcbresp, shrcbMonitorRes.SysNameCn, shrcbMonitorRes.Title)
	if strings.EqualFold("ok", shrcbresp.Msg) {
		return true
	}
	if shrcbresp.Msg == "Ok" {
		return true
	}
	return false
}
func (l *HttpConnectorLogic) GetServicesInfo(SystemInfo *types.System) (*types.ShrcbMonitorRes, int, string) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	majorService := ""
	ServiceProblemDetail := make(map[string]map[string]int)
	ServiceProblemCount := 0
	monitorRes := &types.ShrcbMonitorRes{
		Datetime:        time.Now().Format("2006-01-02 15:04:05"),
		HostName:        SystemInfo.HostName,
		IpAddress:       SystemInfo.InnerIP,
		Severity:        SystemInfo.Severity,
		SysNameCn:       SystemInfo.SystemNameCn,
		SysNameEn:       SystemInfo.SystemNameEn,
		BlindTimeMinute: SystemInfo.BlindTimeMinute,
	}

	urlPath := fmt.Sprintf(l.svcCtx.Config.CDP.ServicesUrl, SystemInfo.OuterIP, SystemInfo.ClusterName)
	client := &http.Client{}
	ServiceNameList := make([]string, 0)
	for serviceName, _ := range SystemInfo.NeedReport {
		if serviceName != "Default" {
			ServiceNameList = append(ServiceNameList, serviceName)
		}
	}
	sreq, _ := http.NewRequest("GET", urlPath, nil)
	sreq.Header.Add("Cookie", GetCacheSession(SystemInfo))
	sresp, _ := client.Do(sreq)
	res := &types.ServicesList{}
	body, _ := ioutil.ReadAll(sresp.Body)
	json.Unmarshal(body, res)
	sresp.Body.Close()
	for _, item := range res.Items {
		fmt.Println(SystemInfo.SystemNameCn, item.Name, item.HealthSummary, item.ServiceState, item.DisplayName, "____________________________________________________________")
		if GetServiceInfo(SystemInfo, item.HealthSummary, item.Name, ServiceNameList) {
			majorService = item.Name
			monitorRes.Title = fmt.Sprintf("%s等服务状态为%s,详细信息:", item.Name, item.HealthSummary)
		}
		for _, check := range item.HealthChecks {
			if GetServiceInfo(SystemInfo, check.Summary, item.Name, ServiceNameList) {
				monitorRes.Msg = fmt.Sprintf("%s等服务状态为%s,%s,详细信息:", check.Name, check.Summary, check.Explanation)
				ServiceProblemCount = ServiceProblemCount + 1
				if _, ok := ServiceProblemDetail[item.Name]; ok {
					ServiceProblemDetail[item.Name][check.Summary] = ServiceProblemDetail[item.Name][check.Summary] + 1
				} else {
					ServiceProblemDetail[item.Name] = make(map[string]int)
					ServiceProblemDetail[item.Name][check.Summary] = 1
				}

			}
		}
	}
	marshalServiceProblemDetail, _ := json.Marshal(ServiceProblemDetail)
	monitorRes.Title = monitorRes.Title + string(marshalServiceProblemDetail)
	monitorRes.Msg = monitorRes.Msg + string(marshalServiceProblemDetail)
	fmt.Printf("=============================================================================================", monitorRes)
	return monitorRes, ServiceProblemCount, majorService
}
func GetServiceInfo(SystemInfo *types.System, Summary string, itemName string, ServiceNameList []string) bool {
	if _, ok := SystemInfo.NeedReport["Default"][Summary]; ok {
		return true
	} else {
		for _, serviceName := range ServiceNameList {
			if strings.Contains(strings.ToLower(itemName), strings.ToLower(serviceName)) {
				if _, ok := SystemInfo.NeedReport[serviceName][Summary]; ok {
					return true
				}
			}
		}
	}
	return false
}
