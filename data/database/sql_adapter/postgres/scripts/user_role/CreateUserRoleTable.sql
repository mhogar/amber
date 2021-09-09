CREATE TABLE "public"."user_role" (
    "user_key" INTEGER,
	"client_key" SMALLINT,
	"role" VARCHAR(15) NOT NULL,
	CONSTRAINT "user_role_pk" PRIMARY KEY ("user_key", "client_key"),
	CONSTRAINT "user_role_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE,
	CONSTRAINT "user_role_client_fk" FOREIGN KEY ("client_key") REFERENCES "client"("key") ON DELETE CASCADE
);