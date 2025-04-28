package queries

const (
	InsertCategory = `
		INSERT INTO categories
			(
				category_code,
				category_name,
				description,
				parent_category_code,
				store_id,
				user_id,
				created_by,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $7)
		RETURNING id;
	`

	GetCategoryByID = `
		SELECT
			category_id,
			category_code,
			category_name,
			description,
			parent_category_code,
			store_id,
			user_id,
			created_at,
			updated_at
		FROM
			categories
		WHERE
			category_id = $1
		AND
			is_active = true
	`

	GetCategoryByName = `
		SELECT
			category_id,
			category_code,
			category_name,
			description,
			parent_category_code,
			store_id,
			user_id,
			created_at,
			updated_at
		FROM
			categories
		WHERE
			LOWER(category_name) = $1
		AND
			is_active = true
	`

	GetCategories = `
		SELECT
			category_id,
			category_name,
			description,
			additional_attributes,
			store_id,
			user_id,
			is_active,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			categories
		WHERE
			is_active = true
	`
)
