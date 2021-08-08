package sqladapter

// SQLScriptRepository is an interface for encapsulating other sql script repository
type SQLScriptRepository interface {
	AccessTokenScriptRepository
	ClientScriptRepository
	MigrationScriptRepository
	UserScriptRepository
}

// AccessTokenScriptRepository is an interface for fetching access token sql scripts
type AccessTokenScriptRepository interface {
	CreateAccessTokenTableScript() string
	DropAccessTokenTableScript() string
	SaveAccessTokenScript() string
	GetAccessTokenByIdScript() string
	DeleteAccessTokenScript() string
	DeleteAllOtherUserTokensScript() string
}

// ClientScriptRepository is an interface for fetching client sql scripts
type ClientScriptRepository interface {
	CreateClientTableScript() string
	DropClientTableScript() string
	SaveClientScript() string
	UpdateClientScript() string
	DeleteClientScript() string
	GetClientByIdScript() string
}

// MigrationScriptRepository is an interface for fetching migration sql scripts
type MigrationScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string
}

// UserScriptRepository is an interface for fetching user sql scripts
type UserScriptRepository interface {
	CreateUserTableScript() string
	DropUserTableScript() string
	SaveUserScript() string
	GetUserByIdScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	DeleteUserScript() string
}
