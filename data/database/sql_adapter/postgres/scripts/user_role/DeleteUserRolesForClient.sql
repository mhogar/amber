DELETE FROM "user_role" ur
    INNER JOIN "client" c ON c."uid" = $1 AND c."key" = ur."client_key"