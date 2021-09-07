CREATE TABLE "public"."user" (
	"key" SERIAL,
	"username" VARCHAR(30) NOT NULL,
	"rank" SMALLINT NOT NULL,
	"password_hash" BYTEA NOT NULL,
	CONSTRAINT "user_pk" PRIMARY KEY ("key"),
	CONSTRAINT "user_username_un" UNIQUE ("username")
);