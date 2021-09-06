UPDATE "user" SET
    "password_hash" = $2,
    "rank" = $3
WHERE "username" = $1