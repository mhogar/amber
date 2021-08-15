CREATE TABLE "public"."access_token" (
	"id" UUID NOT NULL,
	"user_id" INTEGER NOT NULL,
	"client_id" SMALLINT NOT NULL,
	CONSTRAINT "access_token_pk" PRIMARY KEY ("id"),
	CONSTRAINT "access_token_user_fk" FOREIGN KEY ("user_id") REFERENCES "public"."user"("id") ON DELETE CASCADE,
	CONSTRAINT "access_token_client_fk" FOREIGN KEY ("client_id") REFERENCES "public"."client"("id") ON DELETE CASCADE
);