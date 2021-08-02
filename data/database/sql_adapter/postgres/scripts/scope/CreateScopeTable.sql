CREATE TABLE "public"."scope" (
	"id" uuid NOT NULL,
	"name" varchar(15) NOT NULL,
	CONSTRAINT "scope_pk" PRIMARY KEY ("id"),
	CONSTRAINT "scope_name_un" UNIQUE ("name")
);