CREATE TABLE "public"."access_token" (
	"id" uuid NOT NULL,
	"user_id" uuid NOT NULL,
	"client_id" uuid NOT NULL,
	"scope_id" uuid NOT NULL,
	CONSTRAINT "access_token_pk" PRIMARY KEY ("id"),
	CONSTRAINT "access_token_user_fk" FOREIGN KEY ("user_id") REFERENCES "public"."user"("id") ON DELETE CASCADE,
	CONSTRAINT "access_token_client_fk" FOREIGN KEY ("client_id") REFERENCES "public"."client"("id") ON DELETE CASCADE,
	CONSTRAINT "access_token_scope_fk" FOREIGN KEY ("scope_id") REFERENCES "public"."scope"("id") ON DELETE CASCADE
);