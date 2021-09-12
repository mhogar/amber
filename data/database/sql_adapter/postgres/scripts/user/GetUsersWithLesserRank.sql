SELECT u."username", u."rank", u."password_hash"
	FROM "user" u
	WHERE u."rank" < $1
	ORDER BY u."username"