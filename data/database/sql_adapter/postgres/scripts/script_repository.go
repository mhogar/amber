// Auto generated. DO NOT EDIT.

package scripts

// ScriptRepository is an implementation of the sql script repository interface that fetches scripts laoded from sql files.
type ScriptRepository struct {}

// CreateClientScript gets the CreateClient script.
func (ScriptRepository) CreateClientScript() string {
	return `
INSERT INTO "client" ("uid", "name")
	VALUES ($1, $2)
`
}

// CreateClientTableScript gets the CreateClientTable script.
func (ScriptRepository) CreateClientTableScript() string {
	return `
CREATE TABLE "public"."client" (
	"key" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("key"),
	CONSTRAINT "client_uid_un" UNIQUE ("uid")
);
`
}

// DeleteClientScript gets the DeleteClient script.
func (ScriptRepository) DeleteClientScript() string {
	return `
DELETE FROM "client" c
    WHERE c."uid" = $1
`
}

// DropClientTableScript gets the DropClientTable script.
func (ScriptRepository) DropClientTableScript() string {
	return `
DROP TABLE "public"."client"
`
}

// GetClientByUIDScript gets the GetClientByUID script.
func (ScriptRepository) GetClientByUIDScript() string {
	return `
SELECT c."uid", c."name"
	FROM "client" c
WHERE c."uid" = $1
`
}

// UpdateClientScript gets the UpdateClient script.
func (ScriptRepository) UpdateClientScript() string {
	return `
UPDATE "client" SET
    "name" = $2
WHERE "uid" = $1
`
}

// CreateMigrationTableScript gets the CreateMigrationTable script.
func (ScriptRepository) CreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS "public"."migration" (
    "timestamp" CHAR(3) NOT NULL,
    CONSTRAINT "migration_pk" PRIMARY KEY ("timestamp")
);
`
}

// DeleteMigrationByTimestampScript gets the DeleteMigrationByTimestamp script.
func (ScriptRepository) DeleteMigrationByTimestampScript() string {
	return `
DELETE FROM "migration"
   WHERE "timestamp" = $1
`
}

// GetLatestTimestampScript gets the GetLatestTimestamp script.
func (ScriptRepository) GetLatestTimestampScript() string {
	return `
SELECT m."timestamp" FROM "migration" m
    ORDER BY m."timestamp" DESC
    LIMIT 1
`
}

// GetMigrationByTimestampScript gets the GetMigrationByTimestamp script.
func (ScriptRepository) GetMigrationByTimestampScript() string {
	return `
SELECT m."timestamp" 
    FROM "migration" m 
    WHERE m."timestamp" = $1
`
}

// SaveMigrationScript gets the SaveMigration script.
func (ScriptRepository) SaveMigrationScript() string {
	return `
INSERT INTO "migration" ("timestamp") 
    VALUES ($1)
`
}

// CreateSessionTableScript gets the CreateSessionTable script.
func (ScriptRepository) CreateSessionTableScript() string {
	return `
CREATE TABLE "public"."session" (
	"id" UUID NOT NULL,
	"user_key" INTEGER NOT NULL,
	CONSTRAINT "session_pk" PRIMARY KEY ("id"),
	CONSTRAINT "session_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE
);
`
}

// DeleteAllOtherUserSessionsScript gets the DeleteAllOtherUserSessions script.
func (ScriptRepository) DeleteAllOtherUserSessionsScript() string {
	return `
DELETE FROM "session" s
    WHERE s."id" != $1 AND
        s."user_key" IN (
            SELECT u."key" FROM "user" u WHERE u."username" = $2
        )
`
}

// DeleteSessionScript gets the DeleteSession script.
func (ScriptRepository) DeleteSessionScript() string {
	return `
DELETE FROM "session" s
    WHERE s."id" = $1
`
}

// DropSessionTableScript gets the DropSessionTable script.
func (ScriptRepository) DropSessionTableScript() string {
	return `
DROP TABLE "public"."session"
`
}

// GetSessionByIdScript gets the GetSessionById script.
func (ScriptRepository) GetSessionByIdScript() string {
	return `
SELECT
    s."id",
    u."username", u."password_hash"
FROM "session" s
    INNER JOIN "user" u ON u."key" = s."user_key"
WHERE s."id" = $1
`
}

// SaveSessionScript gets the SaveSession script.
func (ScriptRepository) SaveSessionScript() string {
	return `
INSERT INTO "session" ("id", "user_key")
	WITH
		t1 AS (SELECT u."key" FROM "user" u WHERE u."username" = $2)
	SELECT $1, t1."key"
		FROM t1
`
}

// CreateUserScript gets the CreateUser script.
func (ScriptRepository) CreateUserScript() string {
	return `
INSERT INTO "user" ("username", "password_hash")
	VALUES ($1, $2)
`
}

// CreateUserTableScript gets the CreateUserTable script.
func (ScriptRepository) CreateUserTableScript() string {
	return `
CREATE TABLE "public"."user" (
	"key" SERIAL,
	"username" VARCHAR(30) NOT NULL,
	"password_hash" BYTEA NOT NULL,
	CONSTRAINT "user_pk" PRIMARY KEY ("key"),
	CONSTRAINT "user_username_un" UNIQUE ("username")
);
`
}

// DeleteUserScript gets the DeleteUser script.
func (ScriptRepository) DeleteUserScript() string {
	return `
DELETE FROM "user" u
    WHERE u."username" = $1
`
}

// DropUserTableScript gets the DropUserTable script.
func (ScriptRepository) DropUserTableScript() string {
	return `
DROP TABLE "public"."user"
`
}

// GetUserByUsernameScript gets the GetUserByUsername script.
func (ScriptRepository) GetUserByUsernameScript() string {
	return `
SELECT u."username", u."password_hash"
	FROM "user" u
	WHERE u."username" = $1
`
}

// UpdateUserScript gets the UpdateUser script.
func (ScriptRepository) UpdateUserScript() string {
	return `
UPDATE "user" SET
    "password_hash" = $2
WHERE "username" = $1
`
}
