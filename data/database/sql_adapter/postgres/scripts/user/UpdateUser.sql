UPDATE "user" SET
    "password_hash" = $2
WHERE "username" = $1