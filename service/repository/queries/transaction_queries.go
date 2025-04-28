package queries

const (
	InsertTransaction = `
		INSERT INTO transactions
			(
				request_id,
				amount,
				transaction_datetime,
				rrn,
				bill_number,
				currency_code,
				payment_status,
				payment_description,
				customer_id,
				merchant_id,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
		RETURNING id;
	`

	GetTransactionByRequestID = `
		SELECT
			id,
			request_id,
			amount,
			transaction_datetime,
			rrn,
			bill_number,
			currency_code,
			payment_status,
			payment_description,
			customer_id,
			merchant_id,
			created_at,
			updated_at
		FROM
			transactions
		WHERE
			request_id = $1
	`

	GetTransactionByRequestIDAndBillNumber = `
		SELECT
			id,
			request_id,
			amount,
			transaction_datetime,
			rrn,
			bill_number,
			currency_code,
			payment_status,
			payment_description,
			customer_id,
			merchant_id,
			created_at,
			updated_at
		FROM
			transactions
		WHERE
			request_id = $1
		AND bill_number = $2
	`

	GetTransactions = `
		SELECT
			t.id,
			t.request_id,
			t.amount,
			t.transaction_datetime,
			t.rrn,
			t.bill_number,
			t.currency_code,
			t.payment_status,
			t.payment_description,
			t.customer_id,
			COALESCE(m.merchant_id, '') AS merchant_id,
			t.created_at,
			t.updated_at
		FROM
			transactions t
			JOIN customers c ON t.customer_id = c.customer_id
			JOIN merchants m ON t.merchant_id = m.id
		WHERE
			t.request_id != ''
	`
)
