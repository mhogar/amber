SELECT u."username", u."rank", u."password_hash"
	FROM "user" u
	WHERE u."username" = $1