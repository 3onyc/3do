package model

var (
	TODO_LIST_INSERT_QUERY = `
		INSERT INTO
			"todo_lists"
		(
			"id",
			"title",
			"description",
			"created_at",
			"updated_at"
		) VALUES (
			:id,
			:title,
			:description,
			:created_at,
			:updated_at
		)
	`
	TODO_LIST_UPDATE_QUERY = `
		UPDATE
			"todo_lists"
		SET
			"title" = :title,
			"description" = :description,
			"created_at" = :created_at,
			"updated_at" = :updated_at
		WHERE
			"id" = :id
	`
	TODO_LIST_SELECT_QUERY = `
		SELECT
			"id",
			"title",
			"description",
			"created_at",
			"updated_at"
		FROM
			"todo_lists"
	`
	TODO_LIST_SELECT_ID_QUERY = TODO_LIST_SELECT_QUERY + `
		WHERE
			id = ?
	`
)
