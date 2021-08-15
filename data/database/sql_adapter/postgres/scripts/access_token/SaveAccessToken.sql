INSERT INTO "access_token" ("id", "user_key", "client_key")
	WITH
		t1 AS (SELECT u."key" FROM "user" u WHERE u."username" = $2),
		t2 AS (SELECT c."key" FROM "client" c WHERE c."uid" = $3)
	SELECT $1, t1."key", t2."key"
		FROM t1, t2