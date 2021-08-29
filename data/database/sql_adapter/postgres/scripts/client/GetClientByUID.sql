SELECT c."uid", c."name", c."redirect_url", c."token_type", c."key_uri"
	FROM "client" c
WHERE c."uid" = $1