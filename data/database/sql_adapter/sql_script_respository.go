package sqladapter

// SQLScriptRepository is an interface for encapsulating other sql script repository.
type SQLScriptRepository interface {
	AccessTokenScriptRepository
	ClientScriptRepository
	MigrationScriptRepository
	ScopeScriptRepository
	UserScriptRepository
}

// AccessTokenScriptRepository is an interface for fetching access token sql scripts.
type AccessTokenScriptRepository interface {
	CreateAccessTokenTableScript() string
	DropAccessTokenTableScript() string
	SaveAccessTokenScript() string
	GetAccessTokenByIdScript() string
	DeleteAccessTokenScript() string
	DeleteAllOtherUserTokensScript() string
}

// ClientScriptRepository is an interface for fetching client sql scripts.
type ClientScriptRepository interface {
	CreateClientTableScript() string
	DropClientTableScript() string
	SaveClientScript() string
	GetClientByIdScript() string
}

// MigrationScriptRepository is an interface for fetching migration sql scripts.
type MigrationScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string
}

// ScopeScriptRepository is an interface for fetching scope sql scripts.
type ScopeScriptRepository interface {
	CreateScopeTableScript() string
	DropScopeTableScript() string
	SaveScopeScript() string
	GetScopeByNameScript() string
}

// UserScriptRepository is an interface for fetching user sql scripts.
type UserScriptRepository interface {
	CreateUserTableScript() string
	DropUserTableScript() string
	SaveUserScript() string
	GetUserByIdScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	DeleteUserScript() string
}
