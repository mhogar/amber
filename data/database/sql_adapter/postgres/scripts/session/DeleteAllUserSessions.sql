DELETE FROM "session" s
    WHERE s."user_key" IN (
        SELECT u."key" FROM "user" u WHERE u."username" = $1
    )