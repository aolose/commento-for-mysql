package main

import (
	"net/http"
	"time"
)

func domainStatistics(domain string) ([]int64, error) {
	currentTime := time.Now()
	beginTime := currentTime.Add(-time.Duration(30) * 24 * time.Hour)

	rows, err := db.Table("views").
		Select("COUNT(date(view_date))").
		Where("domain = ? AND view_date BETWEEN ? AND ?", domain, beginTime, currentTime).
		Group("date(view_date)").Rows()

	if err != nil {
		logger.Errorf("cannot get daily views: %v", err)
		return []int64{}, errorInternal
	}

	defer rows.Close()

	var last30Days []int64
	for rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			logger.Errorf("cannot get daily views for the last month: %v", err)
			return []int64{}, errorInternal
		}
		last30Days = append(last30Days, count)
	}

	return last30Days, nil
}

func domainStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	viewsLast30Days, err := domainStatistics(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	commentsLast30Days, err := commentStatistics(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "viewsLast30Days": viewsLast30Days, "commentsLast30Days": commentsLast30Days})
}
