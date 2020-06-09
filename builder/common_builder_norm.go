// +build !rhel5

package builder

import (
	"context"
	"sync"

	"github.com/ercole-io/ercole-agent/model"
	"github.com/ercole-io/ercole-agent/utils"
)

func (b *CommonBuilder) getOpenDatabase(entry model.OratabEntry, hostType string) *model.Database {
	dbVersion := b.fetcher.GetDbVersion(entry)

	statsCtx, cancelStatsCtx := context.WithCancel(context.Background())
	if b.configuration.Forcestats {
		utils.RunRoutine(b.configuration, func() {
			b.fetcher.RunStats(entry)

			cancelStatsCtx()
		})
	} else {
		cancelStatsCtx()
	}

	database := b.fetcher.GetOpenDb(entry)

	var wg sync.WaitGroup

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Tablespaces = b.fetcher.GetTablespaces(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Schemas = b.fetcher.GetSchemas(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Patches = b.fetcher.GetPatches(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Features2 = b.fetcher.GetFeatures2(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		<-statsCtx.Done()

		database.Licenses = b.fetcher.GetLicenses(entry, dbVersion, hostType)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.ADDMs = b.fetcher.GetADDMs(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.SegmentAdvisors = b.fetcher.GetSegmentAdvisors(entry)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.LastPSUs = b.fetcher.GetLastPSUs(entry, dbVersion)
	}, &wg)

	utils.RunRoutineInGroup(b.configuration, func() {
		database.Backups = b.fetcher.GetBackups(entry)
	}, &wg)

	wg.Wait()

	return &database
}
