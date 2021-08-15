SELECT c."id", c."uid", c."name"
	FROM "client" c
	WHERE c."uid" = $1