INSERT INTO "session" ("token", "user_key")
	WITH
		t1 AS (SELECT u."key" FROM "user" u WHERE u."username" = $2)
	SELECT $1, t1."key"
		FROM t1