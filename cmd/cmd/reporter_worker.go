package cmd

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/platinummonkey/unifi"
	"github.com/platinummonkey/unifi/cmd/reporters"
	"go.uber.org/zap"
)

type reporterWorker struct {
	site siteConfig
	reporters reporters.Reporters
}

func (w *reporterWorker) ReportAlarmStats() {
	logger.Debug("collecting alarm stats", zap.String("site", w.site.Name))
	// load last alarm state
	alarmLastTimestamp := db.LastAlarmTimestamp(w.site.ID)

	alarmsResp, err := client.SiteAlarmsCount(w.site.ID, 24, false)
	if err == nil {
		w.reporters.ReportMetric(reporters.CountMetricType, "alarm.count", float64(alarmsResp.Meta.Count), "archived:false")
	}
	alarmsResp, err = client.SiteAlarmsCount(w.site.ID, 24, true)
	if err == nil {
		w.reporters.ReportMetric(reporters.CountMetricType, "alarm.count", float64(alarmsResp.Meta.Count), "archived:true")
	}
	offset := 0
	limit := 1000

	now := time.Now().UTC()
	// max 1 week (768 hours)
	historyHours := math.Min(math.Ceil((-time.Duration(now.Unix() - alarmLastTimestamp) * time.Second).Hours()), 768)
	newAlarmsData, err := client.SiteAlarms(w.site.ID, int(historyHours), offset, limit, unifi.EventSortOrderTimeDescending, false)
	if err == nil {
		newCount := 0
		for _, alarm := range newAlarmsData.Data {
			if alarm.Time > alarmLastTimestamp {
				alarmLastTimestamp = alarm.Time
			}
			newCount++
			w.reporters.ReportEvent(
				fmt.Sprintf("Alert: Site=%s[%s] ID=%s Action=%s Category=%s", w.site.Name, w.site.ID, alarm.ID, alarm.InnerAlertAction, alarm.InnerAlertCategory),
				fmt.Sprintf("%s\nTimestamp: %s", alarm.Message, alarm.DatetimeStr),  // TODO: add more info
				fmt.Sprintf("severity:%d", alarm.InnerAlertSeverity),
				fmt.Sprintf("signature_id:%d", alarm.InnerAlertSignatureID),
				fmt.Sprintf("revision:%d", alarm.InnerAlertRevision),
				fmt.Sprintf("gid:%d", alarm.InnerAlertGID),
			)
		}
		w.reporters.ReportMetric(reporters.CountMetricType,"alarm.new.count", float64(newCount))
		db.PersistAlarmTimestamp(w.site.ID, alarmLastTimestamp)
	}
}

func (w *reporterWorker) ReportEventStats() {

}

func (w *reporterWorker) ReportSiteStats() {

}

func (w *reporterWorker) ReportRogueAccessPoints() {

}

func (w *reporterWorker) ReportBackupInfo() {

}

func workerCollectSiteStats(wg *sync.WaitGroup, workChan chan siteConfig, reporters reporters.Reporters) {
	for {
		select {
		case site := <- workChan:
			worker := &reporterWorker{
				site: site,
				reporters: reporters,
			}

			worker.ReportAlarmStats()
			worker.ReportEventStats()
			worker.ReportSiteStats()
			worker.ReportRogueAccessPoints()
			worker.ReportBackupInfo()

			// load the latest state
			// backups, err := client.ListBackups(site.ID)
			// client.SiteReport(site.ID, )
		case <-time.After(time.Millisecond*100):
			// exit no more work
			wg.Done()
			return
		}
	}
}
