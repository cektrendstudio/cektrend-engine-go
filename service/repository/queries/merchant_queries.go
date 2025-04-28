package queries

const (
	InsertMerchant = `
		INSERT INTO merchants
			(
				merchant_id,
				merchant_name,
				merchant_city,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $4)
		RETURNING id;
	`

	GetMerchants = `
		SELECT
			id,
			merchant_id,
			merchant_name,
			merchant_city,
			created_at,
			updated_at
		FROM
			merchants
	`

	GetMerchantByID = `
		SELECT
			id,
			merchant_id,
			merchant_name,
			merchant_city,
			created_at,
			updated_at
		FROM
			merchants
		WHERE
			id = $1
	`

	GetMerchantByMerchantID = `
		SELECT
			id,
			merchant_id,
			merchant_name,
			merchant_city,
			created_at,
			updated_at
		FROM
			merchants
		WHERE
			merchant_id = $1
	`

	UpdateMerchantByMerchantID = `
		UPDATE merchants
		SET
			merchant_name = $2,
			updated_at = $3
		WHERE
			merchant_id = $1
		RETURNING merchant_id;
	`

	DeleteMerchantByMerchantID = `
		DELETE merchants WHERE merchant_id = $1;
	`
)
