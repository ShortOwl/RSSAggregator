-- name: SavePasswordReset :exec
INSERT INTO password_resets (user_id, otp_hash,created_at,expires_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET otp_hash = EXCLUDED.otp_hash,
expires_at = EXCLUDED.expires_at,
created_at = EXCLUDED.created_at;

-- name: GetPasswordResetByUser :one
SELECT user_id, otp_hash, expires_at, created_at
FROM password_resets
WHERE user_id = $1;

-- name: CompletePasswordReset :execrows
WITH valid_reset AS (
    DELETE FROM password_resets
    WHERE user_id = $1
      AND otp_hash = $2
      AND expires_at > $3
    RETURNING user_id
)
UPDATE users
SET password_hash = $4,
    updated_at = $5
WHERE id = (SELECT user_id FROM valid_reset);
