The queries should also be preceded by a comment with the name the return type hint such as `:one`, `:many` or `:exec` using a comment like:
```
-- name: UpsertAuthorName :one
UPDATE author
SET
  name = CASE WHEN @set_name::bool
    THEN @name::text
    ELSE name
    END
RETURNING *;
```
