CREATE TABLE "public"."client" (
	"id" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("id")
);