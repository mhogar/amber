package sqladapter

// SQLScriptRepository is an interface for encapsulating other sql script repositories.
type SQLScriptRepository interface {
	SessionScriptRepository
	ClientScriptRepository
	MigrationScriptRepository
	UserScriptRepository
	UserRoleScriptRepository
}

// SessionScriptRepository is an interface for fetching session sql scripts.
type SessionScriptRepository interface {
	CreateSessionTableScript() string
	DropSessionTableScript() string
	SaveSessionScript() string
	GetSessionByTokenScript() string
	DeleteSessionScript() string
	DeleteAllUserSessionsScript() string
	DeleteAllOtherUserSessionsScript() string
}

// ClientScriptRepository is an interface for fetching client sql scripts.
type ClientScriptRepository interface {
	CreateClientTableScript() string
	DropClientTableScript() string
	CreateClientScript() string
	GetClientsScript() string
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
	GetUsersWithLesserRankScript() string
	GetUserByUsernameScript() string
	UpdateUserScript() string
	UpdateUserPasswordScript() string
	DeleteUserScript() string
}

// UserRoleScriptRepository is an interface for fetching user-role sql scripts.
type UserRoleScriptRepository interface {
	CreateUserRoleTableScript() string
	DropUserRoleTableScript() string
	CreateUserRoleScript() string
	GetUserRolesWithLesserRankByClientUIDScript() string
	GetUserRoleByClientUIDAndUsernameScript() string
	UpdateUserRoleScript() string
	DeleteUserRoleScript() string
}
