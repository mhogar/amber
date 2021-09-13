UPDATE "user_role" SET
    "role" = $3
WHERE "client_key" IN (SELECT c."key" FROM "client" c WHERE c."uid" = $1) AND
      "user_key" IN (SELECT u."key" FROM "user" u WHERE u."username" = $2)