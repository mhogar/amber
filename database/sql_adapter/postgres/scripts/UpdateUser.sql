UPDATE "user" SET
    "username" = $2,
    "password_hash" = $3
WHERE "id" = $1