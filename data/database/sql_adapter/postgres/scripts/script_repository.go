// Auto generated. DO NOT EDIT.

package scripts

// ScriptRepository is an implementation of the sql script repository interface that fetches scripts laoded from sql files.
type ScriptRepository struct {}

// CreateAccessTokenTableScript gets the CreateAccessTokenTable script
func (ScriptRepository) CreateAccessTokenTableScript() string {
	return `
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
    c."id",
    s."id", s."name"
FROM "access_token" tk
    INNER JOIN "user" u ON u."id" = tk."user_id"
    INNER JOIN "client" c ON c."id" = tk."client_id"
    INNER JOIN "scope" s ON s."id" = tk."scope_id"
WHERE tk."id" = $1
`
}

// SaveAccessTokenScript gets the SaveAccessToken script
func (ScriptRepository) SaveAccessTokenScript() string {
	return `
INSERT INTO "access_token" ("id", "user_id", "client_id", "scope_id")
	VALUES ($1, $2, $3, $4)
`
}

// CreateClientTableScript gets the CreateClientTable script
func (ScriptRepository) CreateClientTableScript() string {
	return `
CREATE TABLE "public"."client" (
	"id" uuid NOT NULL,
	CONSTRAINT "client_pk" PRIMARY KEY ("id")
);
`
}

// DropClientTableScript gets the DropClientTable script
func (ScriptRepository) DropClientTableScript() string {
	return `
DROP TABLE "public"."client"
`
}

// GetClientByIdScript gets the GetClientById script
func (ScriptRepository) GetClientByIdScript() string {
	return `
SELECT c."id"
	FROM "client" c
	WHERE c."id" = $1
`
}

// SaveClientScript gets the SaveClient script
func (ScriptRepository) SaveClientScript() string {
	return `
INSERT INTO "client" ("id")
	VALUES ($1)
`
}

// CreateMigrationTableScript gets the CreateMigrationTable script
func (ScriptRepository) CreateMigrationTableScript() string {
	return `
CREATE TABLE IF NOT EXISTS "public"."migration" (
    "timestamp" char(14) NOT NULL,
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

// CreateScopeTableScript gets the CreateScopeTable script
func (ScriptRepository) CreateScopeTableScript() string {
	return `
CREATE TABLE "public"."scope" (
	"id" uuid NOT NULL,
	"name" varchar(15) NOT NULL,
	CONSTRAINT "scope_pk" PRIMARY KEY ("id"),
	CONSTRAINT "scope_name_un" UNIQUE ("name")
);
`
}

// DropScopeTableScript gets the DropScopeTable script
func (ScriptRepository) DropScopeTableScript() string {
	return `
DROP TABLE "public"."scope"
`
}

// GetScopeByNameScript gets the GetScopeByName script
func (ScriptRepository) GetScopeByNameScript() string {
	return `
SELECT s."id", s."name"
	FROM "scope" s
	WHERE s."name" = $1
`
}

// SaveScopeScript gets the SaveScope script
func (ScriptRepository) SaveScopeScript() string {
	return `
INSERT INTO "scope" ("id", "name")
	VALUES ($1, $2)
`
}

// CreateUserTableScript gets the CreateUserTable script
func (ScriptRepository) CreateUserTableScript() string {
	return `
CREATE TABLE "public"."user" (
	"id" uuid NOT NULL,
	"username" varchar(30) NOT NULL,
	"password_hash" bytea NOT NULL,
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

// GetUserByIdScript gets the GetUserById script
func (ScriptRepository) GetUserByIdScript() string {
	return `
SELECT u."id", u."username", u."password_hash"
	FROM "user" u
	WHERE u."id" = $1
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

// SaveUserScript gets the SaveUser script
func (ScriptRepository) SaveUserScript() string {
	return `
INSERT INTO "user" ("id", "username", "password_hash")
	VALUES ($1, $2, $3)
`
}

// UpdateUserScript gets the UpdateUser script
func (ScriptRepository) UpdateUserScript() string {
	return `
UPDATE "user" SET
    "username" = $2,
    "password_hash" = $3
WHERE "id" = $1
`
}
