SELECT
    tk."id",
    u."id", u."username", u."password_hash",
    c."id",
    s."id", s."name"
FROM "access_token" tk
    INNER JOIN "user" u ON u."id" = tk."user_id"
    INNER JOIN "client" c ON c."id" = tk."client_id"
    INNER JOIN "scope" s ON s."id" = tk."scope_id"
WHERE tk."id" = $1