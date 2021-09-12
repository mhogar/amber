SELECT u."username", c."uid", ur."role"
    FROM "user_role" ur
        INNER JOIN "client" c on c."uid" = $1 AND c."key" = ur."client_key"
        INNER JOIN "user" u on u."key" = ur."user_key"
    ORDER BY u."username"