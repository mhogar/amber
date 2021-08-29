UPDATE "client" SET
    "name" = $2,
    "redirect_url" = $3,
    "token_type" = $4,
    "key_uri" = $5
WHERE "uid" = $1