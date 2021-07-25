DELETE FROM "access_token" tk
    WHERE tk."user_id" = $1 AND tk."id" != $2