DELETE FROM "session" s
    WHERE s."id" != $1 AND
        s."user_key" IN (
            SELECT u."key" FROM "user" u WHERE u."username" = $2
        )