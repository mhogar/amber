CREATE TABLE "public"."session" (
	"id" UUID NOT NULL,
	"user_key" INTEGER NOT NULL,
	CONSTRAINT "session_pk" PRIMARY KEY ("id"),
	CONSTRAINT "session_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE
);