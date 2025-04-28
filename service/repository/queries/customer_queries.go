package queries

const (
	InsertCustomer = `
		INSERT INTO customers
			(
				customer_pan,
				customer_name,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $3)
		RETURNING customer_id;
	`

	GetCustomers = `
		SELECT
			customer_id,
			customer_pan,
			customer_name,
			created_at,
			updated_at
		FROM
			customers
	`

	GetCustomerByID = `
		SELECT
			customer_id,
			customer_pan,
			customer_name,
			created_at,
			updated_at
		FROM
			customers
		WHERE
			customer_id = $1
	`

	GetCustomerByCustomerPAN = `
		SELECT
			customer_id,
			customer_pan,
			customer_name,
			created_at,
			updated_at
		FROM
			customers
		WHERE
			customer_pan = $1
	`

	UpdateCustomerByID = `
		UPDATE customers
		SET
			customer_pan = $2,
			customer_name = $2,
			updated_at = $3
		WHERE
			customer_id = $1
		RETURNING customer_id;
	`

	DeleteCustomerByID = `
		DELETE customers WHERE customer_id = $1;
	`
)
