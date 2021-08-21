SELECT
    s."id",
    u."username", u."password_hash"
FROM "session" s
    INNER JOIN "user" u ON u."key" = s."user_key"
WHERE s."id" = $1