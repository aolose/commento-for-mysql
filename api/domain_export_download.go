package main

import (
	"fmt"
	"net/http"
	"time"
)

func domainExportDownloadHandler(w http.ResponseWriter, r *http.Request) {
	exportHex := r.FormValue("exportHex")
	if exportHex == "" {
		fmt.Fprintf(w, "Error: empty exportHex\n")
		return
	}

	statement := `
		SELECT domain, bin_data, creation_Date
		FROM exports
		WHERE export_hex = ?;
	`
	row := db.Raw(statement, exportHex).Row()

	var domain string
	var binData []byte
	var creationDate time.Time
	if err := row.Scan(&domain, &binData, &creationDate); err != nil {
		fmt.Fprintf(w, "Error: that exportHex does not exist\n")
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s-%v.gz"`, domain, creationDate.Unix()))
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(binData)
}
