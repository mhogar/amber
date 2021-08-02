CREATE TABLE "public"."user" (
	"id" uuid NOT NULL,
	"username" varchar(30) NOT NULL,
	"password_hash" bytea NOT NULL,
	CONSTRAINT "user_pk" PRIMARY KEY ("id"),
	CONSTRAINT "user_username_un" UNIQUE ("username")
);