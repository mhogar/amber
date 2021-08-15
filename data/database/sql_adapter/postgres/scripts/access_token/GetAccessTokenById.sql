SELECT
    tk."id",
    u."username", u."password_hash",
    c."uid", c."name"
FROM "access_token" tk
    INNER JOIN "user" u ON u."key" = tk."user_key"
    INNER JOIN "client" c ON c."key" = tk."client_key"
WHERE tk."id" = $1