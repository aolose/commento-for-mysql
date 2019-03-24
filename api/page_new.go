package main

func pageNew(domain string, path string) error {
	// path can be empty
	if domain == "" {
		return errorMissingField
	}

	var pg Pages
	db.First(&Pages{}, "domain = ? AND path = ?", domain, path).Scan(&pg)
	if pg.Domain != "" {
		logger.Errorf("page exist!domain: %v path: %v", domain, path)
	} else {
		db.Create(&Pages{Domain: domain, Path: path})
	}
	return nil
}
