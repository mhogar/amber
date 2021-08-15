DELETE FROM "access_token" tk
    WHERE tk."id" != $1 AND
        tk."user_key" IN (
            SELECT u."key" FROM "user" u WHERE u."username" = $2
        )