-- name: CreateApplication :one
INSERT INTO applications (
  job_title,
  company,
  user_id,
  status,
  location,
  notes,
  application_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: GetApplication :one
SELECT * FROM applications
WHERE id = $1 LIMIT 1;


-- name: ListApplications :many
SELECT * FROM applications
WHERE user_id = $1 AND id > $2
ORDER BY id
LIMIT $3;


-- name: UpdateApplication :one
UPDATE applications
  set job_title = $2,
  company = $3,
  application_date = $4,
  status = $5,
  location = $6,
  notes = $7
WHERE id = $1 RETURNING *;


-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = $1;

