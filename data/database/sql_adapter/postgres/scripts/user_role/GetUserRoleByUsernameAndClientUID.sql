SELECT u."username", c."uid", ur."role"
    FROM "user_role" ur
        INNER JOIN "user" u on u."username" = $1 AND u."key" = ur."user_key"
        INNER JOIN "client" c on c."uid" = $2 AND c."key" = ur."client_key"