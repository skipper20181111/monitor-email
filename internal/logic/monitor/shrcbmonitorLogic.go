package monitor

import (
	"context"
	"encoding/json"
	"monitor/cachemodel"
	"monitor/internal/svc"
	"monitor/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShrcbmonitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	hc     *HttpConnectorLogic
}

func NewShrcbmonitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShrcbmonitorLogic {
	return &ShrcbmonitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		hc:     NewHttpConnectorLogicLogic(ctx, svcCtx),
	}
}

func (l *ShrcbmonitorLogic) Shrcbmonitor() (resp *types.ShrcbMonitorRespList, err error) {
	resp = &types.ShrcbMonitorRespList{
		Code: "10000",
		Msg:  "ok",
		Data: make([]*types.ShrcbMonitorResp, 0),
	}
	SystemList := types.SystemList{}
	get, ok := l.svcCtx.LocalCache.Get(svc.SystemListKey)
	if ok {
		SystemList = get.(types.SystemList)
	}
	for _, SystemInfo := range SystemList.SystemList {
		l.GetReportInfo(SystemInfo)
	}
	return resp, nil
}

func (l *ShrcbmonitorLogic) GetReportInfo(SystemInfo *types.System) *types.ShrcbMonitorResp {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	//l.JudgeInsertDatabase(SystemInfo)
	//l.JudgeInsertDatabaseWithLimit(SystemInfo)
	l.InsertDatabase(SystemInfo)
	return &types.ShrcbMonitorResp{}
}
func (l *ShrcbmonitorLogic) JudgeInsertDatabase(SystemInfo *types.System) {
	if SystemInfo.BlindInsertDatabase {
		lastrecord, _ := l.svcCtx.Monitor.FindOneByName(l.ctx, SystemInfo.SystemNameEn)
		if lastrecord == nil || lastrecord.GenerateTime.Before(time.Now().Add(-time.Minute*time.Duration(SystemInfo.BlindTimeMinute))) {
			l.InsertDatabase(SystemInfo)
		}
	} else {
		l.InsertDatabase(SystemInfo)
	}
}

func (l *ShrcbmonitorLogic) JudgeInsertDatabaseWithLimit(SystemInfo *types.System) {
	if SystemInfo.BlindInsertDatabase {
		if SystemInfo.ConfirmReportNumber < 1 {
			SystemInfo.ConfirmReportNumber = 1
		}
		recordList, _ := l.svcCtx.Monitor.FindOneByNameLimitN(l.ctx, SystemInfo.SystemNameEn, SystemInfo.ConfirmReportNumber)
		if IfInsert(recordList, SystemInfo) {
			l.InsertDatabase(SystemInfo)
		}
	} else {
		l.InsertDatabase(SystemInfo)
	}
}
func IfInsert(recordList []*cachemodel.ShrcbMonitor, SystemInfo *types.System) bool {
	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()
	return true
	//if len(recordList) < SystemInfo.ConfirmReportNumber {
	//	return true
	//}
	//for _, monitor := range recordList {
	//	if monitor.Reported > 0 {
	//		return true
	//	}
	//}
	//if recordList[SystemInfo.ConfirmReportNumber-1].GenerateTime.Before(time.Now().Add(-time.Minute * time.Duration(SystemInfo.BlindTimeMinute*SystemInfo.ConfirmReportNumber))) {
	//	return true
	//}
	//return false
}
func IfReport(recordList []*cachemodel.ShrcbMonitor, SystemInfo *types.System) bool {
	if len(recordList) == SystemInfo.ConfirmReportNumber {
		for _, monitor := range recordList {
			if monitor.Reported > 0 {
				return false
			}
		}
		if recordList[0].GenerateTime.After(time.Now().Add(-time.Minute * time.Duration(SystemInfo.BlindTimeMinute))) {
			return true
			//if recordList[SystemInfo.ConfirmReportNumber-1].GenerateTime.Before(time.Now().Add(-time.Minute * time.Duration(SystemInfo.BlindTimeMinute*SystemInfo.ConfirmReportNumber))){
			//	return true
			//}
		}

	}

	return false
}
func (l *ShrcbmonitorLogic) InsertDatabase(SystemInfo *types.System) {
	inittime, _ := time.Parse("2006-01-02 15:04:05", "2099-01-01 00:00:00")
	monitorRes, count, _ := l.hc.GetServicesInfo(SystemInfo)
	monitorResBytes, _ := json.Marshal(monitorRes)
	if count > 0 {
		shrcbMonitor := &cachemodel.ShrcbMonitor{
			Reported:       0,
			SystemName:     SystemInfo.SystemNameEn,
			SystemNameZh:   SystemInfo.SystemNameCn,
			ReportTitle:    monitorRes.Title,
			ReportMsg:      monitorRes.Msg,
			ReportMarshal:  string(monitorResBytes),
			ReportSeverity: monitorRes.Severity,
			GenerateTime:   time.Now(),
			ReportTime:     inittime,
		}
		if !SystemInfo.TellTheTales {
			shrcbMonitor.Reported = 1
		}
		l.svcCtx.Monitor.Insert(l.ctx, shrcbMonitor)

		//return l.hc.Report2Shrcb(monitorRes)
	}
}
