package queries

const (
	CreateToken = `
		INSERT INTO oauth_tokens
			(
				user_id,
				access_token,
				refresh_token,
				access_token_expires_at,
				refresh_token_expires_at,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $6)
		RETURNING token_id;
	`

	UpdateTokenByTokenID = `
		UPDATE oauth_tokens
		SET
			access_token = $2,
			access_token_expires_at = $3,
			refresh_token = $4,
			refresh_token_expires_at = $5,
			updated_at = $6
		WHERE
			token_id = $1;
	`

	RevokedInactiveToken = `
		UPDATE oauth_tokens
		SET
			revoked = true,
			updated_at = $3
		WHERE
			user_id = $1
		AND
			token_id != $2;
	`

	ValidateAccessToken = `
		SELECT
			user_id,
			access_token_expires_at,
			revoked
		FROM
			oauth_tokens
		WHERE
			access_token = $1
		AND
			user_id = $2
		AND revoked = false
	`

	ValidateRefreshToken = `
		SELECT
			token_id,
			user_id,
			refresh_token,
			refresh_token_expires_at
		FROM
			oauth_tokens
		WHERE
			refresh_token = $1
		AND
			user_id = $2
		AND revoked = false
	`
)
