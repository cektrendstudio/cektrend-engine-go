package queries

const (
	InsertProduct = `
		INSERT INTO products
			(
				product_name,
				description,
				additional_attributes,
				store_id,
				user_id,
				created_by,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $6)
		RETURNING product_id;
	`

	GetProductByID = `
		SELECT
			product_id,
			product_name,
			description,
			additional_attributes,
			store_id,
			user_id,
			is_active,
			created_at,
			updated_at,
		FROM
			products
		WHERE
			product_id = $1
		AND
			is_active = true
	`

	GetProducts = `
		SELECT
			product_id,
			product_name,
			description,
			additional_attributes,
			store_id,
			user_id,
			is_active,
			created_at,
			updated_at,
		FROM
			products
		WHERE
			is_active = true
	`

	GetProductVariants = `
		SELECT
			pv.variant_id AS variant_id,
			COALESCE(pv.name, '') AS name,
			COALESCE(pv.sku, '') AS sku,
			COALESCE(pv.barcode, '') AS barcode,
			pv.is_unlimited_stock AS is_unlimited_stock,
			COALESCE(pv.stock, 0) AS stock,
			COALESCE(pv.min_stock, 0) AS min_stock,
			COALESCE(pv.cost_price, 0.0) AS cost_price,
			COALESCE(pv.price, 0.0) AS price,
			COALESCE(pv.additional_attributes::json, '{}'::json) AS additional_attributes,
			pv.is_enabled AS is_enabled,
			pv.created_at AS created_at,
			pv.updated_at AS updated_at
		FROM
			product_variants pv
			JOIN products p ON pv.product_id = p.product_id
		WHERE
			p.store_id = ?
			AND p.user_id = ?
			AND pv.is_active = true
	`

	GetProductCategoryIdsNotInCategoryIDs = `
		SELECT
			id
		FROM
			product_categories
		WHERE
			product_id = ?
			AND category_id NOT IN (?)
			AND is_active = true
	`

	GetCategoryId = `
		SELECT
			category_id
		FROM
			product_categories
		WHERE
			product_id = $1
			AND is_active = TRUE
	`

	DeleteBulkProductCategories = `
		DELETE FROM product_categories WHERE id IN (?);
	`

	GetProductAttributeOptionByOptionValue = `
		SELECT
			pao.option_id,
			pao.attribute_id,
			pao.option_value
		FROM
			product_attribute_options pao
			JOIN product_attributes pa ON pao.attribute_id = pa.attribute_id
		WHERE
			pa.user_id = ?
			AND pa.store_id = ?
			AND pa.reference_id = ?
			AND pa.reference_type = ?
			AND pao.option_value IN (?)
			AND is_active = true
	`

	GetExistingAttributeID = `
		SELECT
			attribute_id
		FROM
			product_attributes
		WHERE
			attribute_name = $1
			AND reference_id = $2
			AND reference_type = $3
			AND store_id = $4
			AND user_id = $5
			AND is_active = true
	`

	InsertProductAttribute = `
		INSERT INTO product_attributes
			(
				attribute_name,
				reference_id,
				reference_type,
				store_id,
				user_id,
				created_by,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $6)
		RETURNING attribute_id;
	`

	InsertProductVariant = `
		INSERT INTO product_variants
			(
				product_id,
				name,
				sku,
				barcode,
				is_unlimited_stock,
				stock,
				min_stock,
				cost_price,
				price,
				additional_attributes,
				is_enabled,
				created_by,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12)
		RETURNING variant_id;
	`

	InsertBulkProductVariants = `
		INSERT INTO product_variants
			(
				product_id,
				name,
				sku,
				barcode,
				is_unlimited_stock,
				stock,
				min_stock,
				cost_price,
				price,
				additional_attributes,
				is_enabled,
				created_by,
				updated_by
			)
		VALUES
			?;
	`
)
