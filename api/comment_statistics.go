package main

import ()

func commentStatistics(domain string) ([]int64, error) {
	statement := `
		SELECT COUNT(comments.creationDate)
		FROM (
			SELECT to_char(date_trunc('day', (current_date - offs)), 'YYYY-MM-DD') AS date
			FROM generate_series(0, 30, 1) AS offs
		) gen LEFT OUTER JOIN comments
		ON gen.date = to_char(date_trunc('day', comments.creationDate), 'YYYY-MM-DD') AND
		   comments.domain=$1
		GROUP BY gen.date
		ORDER BY gen.date;
	`
	rows, err := db.Query(statement, domain)
	if err != nil {
		logger.Errorf("cannot get daily views: %v", err)
		return []int64{}, errorInternal
	}

	defer rows.Close()

	last30Days := []int64{}
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
