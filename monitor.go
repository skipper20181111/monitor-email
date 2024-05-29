package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpc"
	"monitor/internal/logic/monitor"
	"monitor/internal/types"
	"net/http"
	"time"

	"monitor/internal/config"
	"monitor/internal/handler"
	"monitor/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/monitor-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	go RefreshCache()
	go StartMonitor(ctx)
	go StartTell(ctx)
	go DeleteHistory(ctx)
	server.Start()
}
func RefreshCache() {
	for true {
		fmt.Println("开始刷新")
		time.Sleep(time.Second)
		urlPath := "http://localhost:58888/refresh/refresh"
		resp, _ := httpc.Do(context.Background(), http.MethodGet, urlPath, nil)
		if resp == nil {
			continue
		}
		fmt.Println("结束刷新", resp)
		fmt.Println(resp.Body.Close())

		time.Sleep(time.Second * 50)
	}
}

func StartMonitor(svcCtx *svc.ServiceContext) {
	time.Sleep(time.Second * 10)
	for true {
		fmt.Println("发起监控")
		Monitor(svcCtx)
		fmt.Println("结束监控")
		time.Sleep(time.Minute * time.Duration(svcCtx.Config.Monitor.TimeGapMinute))
		//time.Sleep(time.Second * 30)
	}
}
func Monitor(svcCtx *svc.ServiceContext) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	sm := monitor.NewShrcbmonitorLogic(context.Background(), svcCtx)
	sm.Shrcbmonitor()
}
func StartTell(svcCtx *svc.ServiceContext) {
	time.Sleep(time.Second * 15)
	for true {
		fmt.Println("开始轮询发送监控信息--如果有的话")
		TellTheTales(svcCtx)
		fmt.Println("轮询发送监控信息结束")
		time.Sleep(time.Minute)
	}
}
func TellTheTales(svcCtx *svc.ServiceContext) {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	SystemList := types.SystemList{}
	get, ok := svcCtx.LocalCache.Get(svc.SystemListKey)
	if ok {
		SystemList = get.(types.SystemList)
	}
	for _, SystemInfo := range SystemList.SystemList {
		_, ok1 := svcCtx.LocalCache.Get(SystemInfo.SystemNameEn)
		if ok1 {
			continue
		}
		MonitorInfoList, _ := svcCtx.Monitor.FindOneByNameLimitN(context.Background(), SystemInfo.SystemNameEn, SystemInfo.ConfirmReportNumber)
		if monitor.IfReport(MonitorInfoList, SystemInfo) {
			deadlinectx, _ := context.WithTimeout(context.Background(), time.Second*3)
			shrcbMonitorRes := &types.ShrcbMonitorRes{}
			json.Unmarshal([]byte(MonitorInfoList[0].ReportMarshal), shrcbMonitorRes)
			hc := monitor.NewHttpConnectorLogicLogic(deadlinectx, svcCtx)
			if hc.Report2Shrcb(shrcbMonitorRes) {
				svcCtx.LocalCache.SetWithExpire(SystemInfo.SystemNameEn, true, time.Minute*time.Duration(SystemInfo.BlindTimeMinute))
				svcCtx.Monitor.UpdateReported(context.Background(), 1, MonitorInfoList[0].SystemName)
			}
		}
		_, ok1 = svcCtx.LocalCache.Get(SystemInfo.SystemNameEn)

	}

}
func DeleteHistory(svcCtx *svc.ServiceContext) {
	time.Sleep(time.Second * 5)
	for true {
		fmt.Println(fmt.Sprintf("开始清理%d天之前的历史数据", svcCtx.Config.Monitor.HistoryKeepDay))
		svcCtx.Monitor.DeleteByTime(context.Background(), int64(svcCtx.Config.Monitor.HistoryKeepDay))
		fmt.Println("清理结束")
		time.Sleep(time.Hour * 24)
	}
}
