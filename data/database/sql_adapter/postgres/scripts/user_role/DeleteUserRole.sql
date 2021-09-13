DELETE FROM "user_role" ur
    WHERE ur."client_key" IN (SELECT c."key" FROM "client" c WHERE c."uid" = $1) AND
          ur."user_key" IN (SELECT u."key" FROM "user" u WHERE u."username" = $2)