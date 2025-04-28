package queries

const (
	GetStoreByUserAndStoreID = `
		SELECT
			store_id,
			store_code,
			store_name,
			store_status,
			created_at,
			updated_at
		FROM
			stores
		WHERE
			user_id = $1
			AND store_id = $2
			AND is_active = true
	`
)
