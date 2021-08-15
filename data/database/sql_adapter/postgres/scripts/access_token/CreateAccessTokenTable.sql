CREATE TABLE "public"."access_token" (
	"id" UUID NOT NULL,
	"user_key" INTEGER NOT NULL,
	"client_key" SMALLINT NOT NULL,
	CONSTRAINT "access_token_pk" PRIMARY KEY ("id"),
	CONSTRAINT "access_token_user_fk" FOREIGN KEY ("user_key") REFERENCES "public"."user"("key") ON DELETE CASCADE,
	CONSTRAINT "access_token_client_fk" FOREIGN KEY ("client_key") REFERENCES "public"."client"("key") ON DELETE CASCADE
);