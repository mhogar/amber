SELECT u."username", u."password_hash", u."rank"
	FROM "user" u
	WHERE u."username" = $1