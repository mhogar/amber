SELECT
    s."token",
    u."username"
FROM "session" s
    INNER JOIN "user" u ON u."key" = s."user_key"
WHERE s."token" = $1