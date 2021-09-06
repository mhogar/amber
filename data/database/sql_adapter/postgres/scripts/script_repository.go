// Auto generated. DO NOT EDIT.

package scripts

// ScriptRepository is an implementation of the sql script repository interface that fetches scripts laoded from sql files.
type ScriptRepository struct {}

// CreateClientScript gets the CreateClient script.
func (ScriptRepository) CreateClientScript() string {
	return `
INSERT INTO "client" ("uid", "name", "redirect_url", "token_type", "key_uri")
	VALUES ($1, $2, $3, $4, $5)
`
}

// CreateClientTableScript gets the CreateClientTable script.
func (ScriptRepository) CreateClientTableScript() string {
	return `
CREATE TABLE "public"."client" (
	"key" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	"redirect_url" VARCHAR(100) NOT NULL,
	"token_type" SMALLINT NOT NULL,
	"key_uri" VARCHAR(100) NOT NULL,
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
SELECT c."uid", c."name", c."redirect_url", c."token_type", c."key_uri"
	FROM "client" c
WHERE c."uid" = $1
`
}

// UpdateClientScript gets the UpdateClient script.
func (ScriptRepository) UpdateClientScript() string {
	return `
UPDATE "client" SET
    "name" = $2,
    "redirect_url" = $3,
    "token_type" = $4,
    "key_uri" = $5
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
	"token" UUID NOT NULL,
	"user_key" INTEGER NOT NULL,
	CONSTRAINT "session_pk" PRIMARY KEY ("token"),
	CONSTRAINT "session_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE
);
`
}

// DeleteAllOtherUserSessionsScript gets the DeleteAllOtherUserSessions script.
func (ScriptRepository) DeleteAllOtherUserSessionsScript() string {
	return `
DELETE FROM "session" s
    WHERE s."token" != $1 AND
        s."user_key" IN (
            SELECT u."key" FROM "user" u WHERE u."username" = $2
        )
`
}

// DeleteSessionScript gets the DeleteSession script.
func (ScriptRepository) DeleteSessionScript() string {
	return `
DELETE FROM "session" s
    WHERE s."token" = $1
`
}

// DropSessionTableScript gets the DropSessionTable script.
func (ScriptRepository) DropSessionTableScript() string {
	return `
DROP TABLE "public"."session"
`
}

// GetSessionByTokenScript gets the GetSessionByToken script.
func (ScriptRepository) GetSessionByTokenScript() string {
	return `
SELECT
    s."token",
    u."username",
    u."rank"
FROM "session" s
    INNER JOIN "user" u ON u."key" = s."user_key"
WHERE s."token" = $1
`
}

// SaveSessionScript gets the SaveSession script.
func (ScriptRepository) SaveSessionScript() string {
	return `
INSERT INTO "session" ("token", "user_key")
	WITH
		t1 AS (SELECT u."key" FROM "user" u WHERE u."username" = $2)
	SELECT $1, t1."key"
		FROM t1
`
}

// CreateUserScript gets the CreateUser script.
func (ScriptRepository) CreateUserScript() string {
	return `
INSERT INTO "user" ("username", "rank", "password_hash")
	VALUES ($1, $2, $3)
`
}

// CreateUserTableScript gets the CreateUserTable script.
func (ScriptRepository) CreateUserTableScript() string {
	return `
CREATE TABLE "public"."user" (
	"key" SERIAL,
	"username" VARCHAR(30) NOT NULL,
	"rank" SMALLINT NOT NULL,
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
SELECT u."username", u."rank", u."password_hash"
	FROM "user" u
	WHERE u."username" = $1
`
}

// UpdateUserScript gets the UpdateUser script.
func (ScriptRepository) UpdateUserScript() string {
	return `
UPDATE "user" SET
    "rank" = $2
WHERE "username" = $1
`
}

// UpdateUserPasswordScript gets the UpdateUserPassword script.
func (ScriptRepository) UpdateUserPasswordScript() string {
	return `
UPDATE "user" SET
    "password_hash" = $2
WHERE "username" = $1
`
}

// AddUserRoleForClientScript gets the AddUserRoleForClient script.
func (ScriptRepository) AddUserRoleForClientScript() string {
	return `
INSERT INTO "user_role" ("client_key", "user_key", "role")
    WITH
        t1 AS (SELECT c."key" FROM "client" c WHERE c."uid" = $1),
		t2 AS (SELECT u."key" FROM "user" u WHERE u."username" = $2)
	SELECT t1."key", t2."key", $3
		FROM t1, t2
`
}

// CreateUserRoleTableScript gets the CreateUserRoleTable script.
func (ScriptRepository) CreateUserRoleTableScript() string {
	return `
CREATE TABLE "public"."user_role" (
	"client_key" SMALLINT,
    "user_key" INTEGER,
	"role" VARCHAR(15) NOT NULL,
	CONSTRAINT "user_role_pk" PRIMARY KEY ("client_key", "user_key"),
	CONSTRAINT "user_role_client_fk" FOREIGN KEY ("client_key") REFERENCES "client"("key") ON DELETE CASCADE,
	CONSTRAINT "user_role_user_fk" FOREIGN KEY ("user_key") REFERENCES "user"("key") ON DELETE CASCADE
);
`
}

// DeleteUserRolesForClientScript gets the DeleteUserRolesForClient script.
func (ScriptRepository) DeleteUserRolesForClientScript() string {
	return `
DELETE FROM "user_role" ur
    WHERE ur."client_key" IN (
        SELECT c."key" FROM "client" c WHERE c."uid" = $1
    )
`
}

// DropUserRoleTableScript gets the DropUserRoleTable script.
func (ScriptRepository) DropUserRoleTableScript() string {
	return `
DROP TABLE "public"."user_role"
`
}

// GetUserRoleForClientScript gets the GetUserRoleForClient script.
func (ScriptRepository) GetUserRoleForClientScript() string {
	return `
SELECT
    u."username",
    ur."role"
FROM "user_role" ur
    INNER JOIN "client" c on c."uid" = $1 AND c."key" = ur."client_key"
    INNER JOIN "user" u on u."username" = $2 AND u."key" = ur."user_key"
`
}

// GetUserRolesForClientScript gets the GetUserRolesForClient script.
func (ScriptRepository) GetUserRolesForClientScript() string {
	return `
SELECT
    u."username",
    ur."role"
FROM "user_role" ur
    INNER JOIN "client" c on c."uid" = $1 AND c."key" = ur."client_key"
    INNER JOIN "user" u on u."key" = ur."user_key"
`
}
