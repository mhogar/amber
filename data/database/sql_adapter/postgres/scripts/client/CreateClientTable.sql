CREATE TABLE "public"."client" (
	"key" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("key"),
	CONSTRAINT "client_uid_un" UNIQUE ("uid")
);