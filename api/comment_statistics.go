package main

import "time"

func commentStatistics(domain string) ([]int64, error) {
	currentTime := time.Now()
	beginTime := currentTime.Add(-time.Duration(30) * 24 * time.Hour)

	rows, err := db.Table("comments").
		Select("COUNT(date(creation_date))").
		Where("domain = ? AND creation_date BETWEEN ? AND ?", domain, beginTime, currentTime).
		Group("date(creation_date)").Rows()
	if err != nil {
		logger.Errorf("cannot get daily views: %v", err)
		return []int64{}, errorInternal
	}
	defer rows.Close()
	var last30Days []int64
	for rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			logger.Errorf("cannot get daily comments for the last month: %v", err)
			return make([]int64, 0), errorInternal
		}
		last30Days = append(last30Days, count)
	}

	return last30Days, nil
}
