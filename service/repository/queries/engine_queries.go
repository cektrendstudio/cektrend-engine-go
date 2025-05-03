package queries

const (
	InsertPhisingWebReport = `
		INSERT INTO phishing_web_report
			(
				site_url,
				image_url,
				version,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id;
	`

	GetPhishingReportByURL = `
		SELECT
			id,
			site_url,
			image_url,
			created_at,
			updated_at
		FROM
			phishing_web_report
		WHERE
			site_url = $1
	`

	GetPhishingReportByURLs = `
		SELECT
			id,
			site_url,
			image_url,
			created_at,
			updated_at
		FROM
			phishing_web_report
		WHERE
			version = ? AND site_url IN (?)
	`
)
