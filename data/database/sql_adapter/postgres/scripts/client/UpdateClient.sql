UPDATE "client" SET
    "name" = $2,
    "redirect_url" = $3
WHERE "uid" = $1