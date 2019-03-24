package main

import (
	"time"
)

func viewsCleanupBegin() error {
	go func() {
		for {
			err := db.Delete(&Views{}, "view_date < ?", time.Now().UTC().AddDate(0, 0, -45)).Error
			if err != nil {
				logger.Errorf("error cleaning up views: %v", err)
				return
			}

			time.Sleep(24 * time.Hour)
		}
	}()

	return nil
}
