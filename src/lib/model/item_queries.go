package model

var (
	TODO_ITEM_INSERT_QUERY = `
		INSERT INTO
			"todo_items"
		(
			"id",
			"title",
			"description",
			"done",
			"done_at",
			"created_at",
			"updated_at",
			"group_id"
		) VALUES (
			:id,
			:title,
			:description,
			:done,
			:done_at,
			:created_at,
			:updated_at,
			:group_id
		)
	`
	TODO_ITEM_UPDATE_QUERY = `
		UPDATE
			"todo_items"
		SET
			"title" = :title,
			"description" = :description,
			"done" = :done,
			"done_at" = :done_at,
			"created_at" = :created_at,
			"updated_at" = :updated_at,
			"group_id" = :group_id
		WHERE
			"id" = :id
	`
	TODO_ITEM_SELECT_QUERY = `
		SELECT
			"id",
			"title",
			"description",
			"done",
			"done_at",
			"created_at",
			"updated_at",
			"group_id"
		FROM
			"todo_items"
	`
	TODO_ITEM_SELECT_ID_QUERY = TODO_ITEM_SELECT_QUERY + `
		WHERE
			id = ?
	`
	TODO_ITEM_SELECT_WITH_GROUP_QUERY = TODO_ITEM_SELECT_QUERY + `
		WHERE
			group_id = ?
	`
	TODO_ITEM_SELECT_IDS_WITH_GROUP_QUERY = `
		SELECT 
			"id" 
		FROM 
			"todo_items"
		WHERE 
			"group_id" = ?
	`
)
