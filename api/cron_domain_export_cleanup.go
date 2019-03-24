package main

import (
	"time"
)

func domainExportCleanupBegin() error {
	go func() {
		for {
			err := db.Delete(&Exports{}, "creation_date < ?", time.Now().UTC().AddDate(0, 0, -7)).Error
			if err != nil {
				logger.Errorf("error cleaning up export rows: %v", err)
				return
			}
			time.Sleep(2 * time.Hour)
		}
	}()

	return nil
}
