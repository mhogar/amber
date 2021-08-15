SELECT
    tk."id",
    u."id", u."username", u."password_hash",
    c."id", c."uid", c."name"
FROM "access_token" tk
    INNER JOIN "user" u ON u."id" = tk."user_id"
    INNER JOIN "client" c ON c."id" = tk."client_id"
WHERE tk."id" = $1