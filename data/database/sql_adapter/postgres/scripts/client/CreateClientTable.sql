CREATE TABLE "public"."client" (
	"key" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	"redirect_url" VARCHAR(100) NOT NULL,
	"token_type" SMALLINT NOT NULL,
	"key_uri" VARCHAR(100) NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("key"),
	CONSTRAINT "client_uid_un" UNIQUE ("uid")
);