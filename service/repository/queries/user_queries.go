package queries

const (
	RegisterUser = `
		INSERT INTO users
			(
				name,
				email,
				password,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING user_id;
	`

	GetUserByID = `
		SELECT
			user_id,
			name,
			email,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE
			user_id = $1
	`

	GetUserByEmail = `
		SELECT
			user_id,
			name,
			email,
			password,
			created_at,
			updated_at
		FROM
			users
		WHERE
			email = $1
	`
)
