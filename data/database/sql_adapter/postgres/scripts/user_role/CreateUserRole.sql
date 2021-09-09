INSERT INTO "user_role" ("user_key", "client_key", "role")
    WITH
		t1 AS (SELECT u."key" FROM "user" u WHERE u."username" = $1),
        t2 AS (SELECT c."key" FROM "client" c WHERE c."uid" = $2)
	SELECT t1."key", t2."key", $3
		FROM t1, t2