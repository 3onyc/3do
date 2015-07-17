package model

var (
	TODO_GROUP_INSERT_QUERY = `
		INSERT INTO
			"todo_groups"
		(
			"id",
			"title",
			"created_at",
			"updated_at",
			"list_id"
		) VALUES (
			:id,
			:title,
			:created_at,
			:updated_at,
			:list_id
		)
	`
	TODO_GROUP_SELECT_QUERY = `
		SELECT
			"id",
			"title",
			"created_at",
			"updated_at",
			"list_id"
		FROM
			"todo_groups"
	`
	TODO_GROUP_SELECT_ID_QUERY = TODO_GROUP_SELECT_QUERY + `
		WHERE
			id = ?
	`
	TODO_GROUP_SELECT_WITH_LIST_QUERY = TODO_GROUP_SELECT_QUERY + `
		WHERE
			list_id = ?
	`
	TODO_GROUP_SELECT_IDS_WITH_LIST_QUERY = `
		SELECT 
			"id" 
		FROM 
			"todo_groups"
		WHERE 
			"list_id" = ?
	`
)
