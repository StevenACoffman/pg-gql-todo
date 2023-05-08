-- name: CreateTodo :one
INSERT INTO todo (
    description
) VALUES (
             $1
         ) RETURNING *;

-- name: GetTodo :one
SELECT * FROM todo
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todo
ORDER BY created_at;

-- name: UpdateTodo :one
UPDATE todo
SET done = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todo
WHERE id = $1;
