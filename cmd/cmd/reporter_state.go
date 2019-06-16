package cmd

import (
	"os"
	"path"

	"github.com/platinummonkey/unifi"
	"github.com/spf13/viper"
	"github.com/timshannon/badgerhold"
	"go.uber.org/zap"
)

// reporter state maintains an on-disk
type reporterState struct {
	db *badgerhold.Store
}

var db *reporterState

func initReporterState() {
	opts := badgerhold.DefaultOptions
	stateDir := viper.GetString("state_dir")
	if stateDir == "" {
		stateDir = path.Join(homeDir, ".unifi.state")
	}
	opts.Dir = stateDir
	opts.ValueDir = stateDir


	// check to make sure the directory exists
	_, err := os.Stat(stateDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(stateDir, 0755)
			if err != nil {
				logger.Error("unable to create state directory", zap.String("directory", stateDir), zap.Error(err))
				os.Exit(-1)
			}
			os.Chown(stateDir, os.Getuid(), os.Getgid())
			// if err != nil {
				logger.Error("unable to create state directory", zap.String("directory", stateDir), zap.Error(err))
				os.Exit(-1)
			}
		}
	}
	badgerDB, err := badgerhold.Open(opts)
	if err != nil {
		logger.Error("unable to open state", zap.String("directory", stateDir), zap.Error(err))
	}
	db = &reporterState{
		db: badgerDB,
	}
}

func (d *reporterState) keyFor(site string, path string) string {
	return site + "_" + path
}

type LastEventState struct {
	ID string `badgerhold:"key"`
	Site string `badgerholdIndex:"siteIdx"`
	Timestamp int64
}

func (d *reporterState) LastEventTimestamp(site string) int64 {
	var s LastEventState
	err := d.db.Get(d.keyFor(site, "event_last"), &s)
	if err != nil {
		if err != badgerhold.ErrNotFound {
			logger.Warn("unable to query event timestamp", zap.String("site", site), zap.Error(err))
		}
		return 0
	}
	return s.Timestamp
}

func (d *reporterState) PersistEventTimestamp(site string, ts int64) {
	s := LastEventState{Site: site, Timestamp: ts}
	err := d.db.Upsert(d.keyFor(site, "event_last"), &s)
	if err != nil {
		logger.Warn("unable to persist event timestamp", zap.String("site", site), zap.Error(err))
	}
}

func (d *reporterState) LastAlarmTimestamp(site string) int64 {
	var s LastEventState
	key := d.keyFor(site, "alarm_last")
	err := d.db.Get(key, &s)
	if err != nil {
		if err != badgerhold.ErrNotFound {
			logger.Warn("unable to query alarm timestamp", zap.String("site", site), zap.Error(err))
		}
		return 0
	}
	return s.Timestamp
}

// PersistAlarmTimestamp will persist the alarm timestamp
func (d *reporterState) PersistAlarmTimestamp(site string, ts int64) {
	s := LastEventState{Site: site, Timestamp: ts}
	err := d.db.Upsert(d.keyFor(site, "alarm_last"), &s)
	if err != nil {
		logger.Warn("unable to persist alarm timestamp", zap.String("site", site), zap.Error(err))
	}
}

type RogueAccessPointsState struct {
	ID string `badgerhold:"key"`
	Site string `badgerholdIndex:"siteIdx"`
	AccessPoint unifi.SiteRougeAccessPoint
}

func (d *reporterState) LastRogueAccessPoints(site string) []unifi.SiteRougeAccessPoint {
	var accessPoints []unifi.SiteRougeAccessPoint
	err := d.db.Find(accessPoints, badgerhold.Where("Site").Eq(site).Index("siteIdx"))
	if err != nil {
		if err != badgerhold.ErrNotFound {
			logger.Warn("unable to query last rogue access points", zap.String("site", site), zap.Error(err))
		}
	}
	return accessPoints
}

func (d *reporterState) PersistRogueAccessPoints(site string, accessPoints []unifi.SiteRougeAccessPoint) {
	for _, ap := range accessPoints {
		key := d.keyFor(site, "rogue_access_point_" + ap.ID)
		d.db.Upsert(key, &RogueAccessPointsState{
			Site: site,
			AccessPoint: ap,
		})
	}
}

func (d *reporterState) LastReportDataPointByReport(site string, interval unifi.ReportInterval, reportType unifi.ReportType) {
	panic("not implemented")
}

func (d *reporterState) PersistReportDataPointByReport(site string, interval unifi.ReportInterval, reportType unifi.ReportType) {
	panic("not implemented")
}