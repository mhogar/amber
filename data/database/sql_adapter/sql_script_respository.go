package sqladapter

// SQLScriptRepository is an interface for encapsulating other sql script repository.
type SQLScriptRepository interface {
	SessionScriptRepository
	ClientScriptRepository
	MigrationScriptRepository
	UserScriptRepository
}

// SessionScriptRepository is an interface for fetching session sql scripts.
type SessionScriptRepository interface {
	CreateSessionTableScript() string
	DropSessionTableScript() string
	SaveSessionScript() string
	GetSessionByIdScript() string
	DeleteSessionScript() string
	DeleteAllOtherUserSessionsScript() string
}

// ClientScriptRepository is an interface for fetching client sql scripts.
type ClientScriptRepository interface {
	CreateClientTableScript() string
	DropClientTableScript() string
	CreateClientScript() string
	GetClientByUIDScript() string
	UpdateClientScript() string
	DeleteClientScript() string
}

// MigrationScriptRepository is an interface for fetching migration sql scripts.
type MigrationScriptRepository interface {
	CreateMigrationTableScript() string
	SaveMigrationScript() string
	GetMigrationByTimestampScript() string
	GetLatestTimestampScript() string
	DeleteMigrationByTimestampScript() string
}

// UserScriptRepository is an interface for fetching user sql scripts.
type UserScriptRepository interface {
	CreateUserTableScript() string
	DropUserTableScript() string
	CreateUserScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	DeleteUserScript() string
}
