SELECT c."uid", c."name", c."redirect_url"
	FROM "client" c
WHERE c."uid" = $1