// Auto generated. DO NOT EDIT.

package scripts

// ScriptRepository is an implementation of the sql script repository interface that fetches scripts laoded from sql files.
type ScriptRepository struct {}

// CreateAccessTokenTableScript gets the CreateAccessTokenTable script
func (ScriptRepository) CreateAccessTokenTableScript() string {
	return `
CREATE TABLE "public"."access_token" (
	"id" UUID NOT NULL,
	"user_id" INTEGER NOT NULL,
	"client_id" SMALLINT NOT NULL,
	CONSTRAINT "access_token_pk" PRIMARY KEY ("id"),
	CONSTRAINT "access_token_user_fk" FOREIGN KEY ("user_id") REFERENCES "public"."user"("id") ON DELETE CASCADE,
	CONSTRAINT "access_token_client_fk" FOREIGN KEY ("client_id") REFERENCES "public"."client"("id") ON DELETE CASCADE
);
`
}

// DeleteAccessTokenScript gets the DeleteAccessToken script
func (ScriptRepository) DeleteAccessTokenScript() string {
	return `
DELETE FROM "access_token" tk
    WHERE tk."id" = $1
`
}

// DeleteAllOtherUserTokensScript gets the DeleteAllOtherUserTokens script
func (ScriptRepository) DeleteAllOtherUserTokensScript() string {
	return `
DELETE FROM "access_token" tk
    WHERE tk."user_id" = $1 AND tk."id" != $2
`
}

// DropAccessTokenTableScript gets the DropAccessTokenTable script
func (ScriptRepository) DropAccessTokenTableScript() string {
	return `
DROP TABLE "public"."access_token"
`
}

// GetAccessTokenByIdScript gets the GetAccessTokenById script
func (ScriptRepository) GetAccessTokenByIdScript() string {
	return `
SELECT
    tk."id",
    u."id", u."username", u."password_hash",
    c."id", c."uid", c."name"
FROM "access_token" tk
    INNER JOIN "user" u ON u."id" = tk."user_id"
    INNER JOIN "client" c ON c."id" = tk."client_id"
WHERE tk."id" = $1
`
}

// SaveAccessTokenScript gets the SaveAccessToken script
func (ScriptRepository) SaveAccessTokenScript() string {
	return `
INSERT INTO "access_token" ("id", "user_id", "client_id")
	VALUES ($1, $2, $3)
`
}

// CreateClientScript gets the CreateClient script
func (ScriptRepository) CreateClientScript() string {
	return `
INSERT INTO "client" ("uid", "name")
	VALUES ($1, $2)
`
}

// CreateClientTableScript gets the CreateClientTable script
func (ScriptRepository) CreateClientTableScript() string {
	return `
CREATE TABLE "public"."client" (
	"id" SMALLSERIAL,
	"uid" UUID NOT NULL,
	"name" VARCHAR(30) NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("id")
);
`
}

// DeleteClientScript gets the DeleteClient script
func (ScriptRepository) DeleteClientScript() string {
	return `
DELETE FROM "client" c
    WHERE c."id" = $1
`
}

// DropClientTableScript gets the DropClientTable script
func (ScriptRepository) DropClientTableScript() string {
	return `
DROP TABLE "public"."client"
`
}

// GetClientByUIDScript gets the GetClientByUID script
func (ScriptRepository) GetClientByUIDScript() string {
	return `
SELECT c."id", c."uid", c."name"
	FROM "client" c
	WHERE c."uid" = $1
`
}

// UpdateClientScript gets the UpdateClient script
func (ScriptRepository) UpdateClientScript() string {
	return `
UPDATE "client" SET
    "name" = $2
WHERE "id" = $1
`
}

// CreateMigrationTableScript gets the CreateMigrationTable script
func (ScriptRepository) CreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS "public"."migration" (
    "timestamp" CHAR(14) NOT NULL,
    CONSTRAINT "migration_pk" PRIMARY KEY ("timestamp")
);
`
}

// DeleteMigrationByTimestampScript gets the DeleteMigrationByTimestamp script
func (ScriptRepository) DeleteMigrationByTimestampScript() string {
	return `
DELETE FROM "migration"
   WHERE "timestamp" = $1
`
}

// GetLatestTimestampScript gets the GetLatestTimestamp script
func (ScriptRepository) GetLatestTimestampScript() string {
	return `
SELECT m."timestamp" FROM "migration" m
    ORDER BY m."timestamp" DESC
    LIMIT 1
`
}

// GetMigrationByTimestampScript gets the GetMigrationByTimestamp script
func (ScriptRepository) GetMigrationByTimestampScript() string {
	return `
SELECT m."timestamp" 
    FROM "migration" m 
    WHERE m."timestamp" = $1
`
}

// SaveMigrationScript gets the SaveMigration script
func (ScriptRepository) SaveMigrationScript() string {
	return `
INSERT INTO "migration" ("timestamp") 
    VALUES ($1)
`
}

// CreateUserScript gets the CreateUser script
func (ScriptRepository) CreateUserScript() string {
	return `
INSERT INTO "user" ("username", "password_hash")
	VALUES ($1, $2)
`
}

// CreateUserTableScript gets the CreateUserTable script
func (ScriptRepository) CreateUserTableScript() string {
	return `
CREATE TABLE "public"."user" (
	"id" SERIAL,
	"username" VARCHAR(30) NOT NULL,
	"password_hash" BYTEA NOT NULL,
	CONSTRAINT "user_pk" PRIMARY KEY ("id"),
	CONSTRAINT "user_username_un" UNIQUE ("username")
);
`
}

// DeleteUserScript gets the DeleteUser script
func (ScriptRepository) DeleteUserScript() string {
	return `
DELETE FROM "user" u
    WHERE u."id" = $1
`
}

// DropUserTableScript gets the DropUserTable script
func (ScriptRepository) DropUserTableScript() string {
	return `
DROP TABLE "public"."user"
`
}

// GetUserByUsernameScript gets the GetUserByUsername script
func (ScriptRepository) GetUserByUsernameScript() string {
	return `
SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."username" = $1
`
}

// UpdateUserScript gets the UpdateUser script
func (ScriptRepository) UpdateUserScript() string {
	return `
UPDATE "user" SET
    "password_hash" = $2
WHERE "id" = $1
`
}
