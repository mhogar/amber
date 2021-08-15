UPDATE "user" SET
    "password_hash" = $2
WHERE "id" = $1