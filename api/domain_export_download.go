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
		SELECT domain, binData, creationDate
		FROM exports
		WHERE exportHex = $1;
	`
	row := db.QueryRow(statement, exportHex)

	var domain string
	var binData []byte
	var creationDate time.Time
	if err := row.Scan(&domain, &binData, &creationDate); err != nil {
		fmt.Fprintf(w, "Error: that exportHex does not exist\n")
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s-%v.json.gz"`, domain, creationDate.Unix()))
	w.Write(binData)
}
