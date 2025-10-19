-- name: CreateFunction :one
INSERT INTO functions (
  function_name,
  embedding,
  js
) VALUES (
  $1, $2, $3
)
ON CONFLICT (function_name) DO NOTHING
RETURNING *;

-- name: CosineSimilarity :many
SELECT
  function_id,
  function_name,
  embedding,
  js,
  embedding <=> $1 AS distance
FROM functions
ORDER BY distance ASC
LIMIT $2;
