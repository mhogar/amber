CREATE TABLE "public"."user" (
	"id" SERIAL,
	"username" VARCHAR(30) NOT NULL,
	"password_hash" BYTEA NOT NULL,
	CONSTRAINT "user_pk" PRIMARY KEY ("id"),
	CONSTRAINT "user_username_un" UNIQUE ("username")
);