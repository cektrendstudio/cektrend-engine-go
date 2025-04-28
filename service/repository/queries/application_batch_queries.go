package queries

const (
	InsertApplicationBatch = `
		INSERT INTO application_batch
			(
				batch_name,
				file_name,
				file_url,
				application_type,
				processing_status,
				processing_started_at,
				processing_completed_at,
				total_row,
				total_success,
				properties,
				store_id,
				user_id,
				created_by,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $13)
		RETURNING batch_id;
	`

	GetApplicationBatchByBatchName = `
		SELECT
			batch_id,
			batch_name,
			file_name,
			file_url,
			application_type,
			processing_status,
			processing_started_at,
			processing_completed_at,
			total_row,
			total_success,
			properties,
			store_id,
			user_id,
			is_active,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			application_batch
		WHERE
			batch_name = $1
		AND
			is_active = true
	`
)
