CREATE TABLE "public"."user_role" (
	"client_key" SMALLINT,
    "user_key" INTEGER,
	"role" VARCHAR(15) NOT NULL,
	CONSTRAINT "user_role_pk" PRIMARY KEY ("client_key", "user_key"),
	CONSTRAINT "user_role_client_fk" FOREIGN KEY ("client_key") REFERENCES "client"("key") ON DELETE CASCADE
	CONSTRAINT "user_role_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE,
);