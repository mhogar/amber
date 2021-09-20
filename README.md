# Amber ![CI Status](https://github.com/mhogar/amber/actions/workflows/CI.yml/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/mhogar/amber/badge.svg)](https://coveralls.io/github/mhogar/amber) [![GoDoc](https://godoc.org/github.com/mhogar/amber?status.svg)](https://pkg.go.dev/github.com/mhogar/amber)

Amber is a micro-service for user management and authentication. Manage all your users in one place, then provide secure and centralized authentication across all your client applications.

## Usage

### API

Amber is built as a REST API so it can better be integrated into any desired system. For details view the [Postman API Documentation](https://documenter.getpostman.com/view/11281814/UUxtEqag).

### Authenticating for a Client

On top of the REST API, Amber provides a login view to ensure the correct handling of user credentials when authenticating. Clients should provide a link to the view, which can be found at `/token?client_id=...` (providing their correct client id). Upon successful authentication, the view will automatically redirect to the URL configured in the client with the appended token.

Tokens are JWTs and provide information about the user including their username and role. They should not be used directly as session tokens, but instead processed by the application to create a new session using their encoded data.

## Building and Tools

Amber is a pure golang application. It can be built/run using standard go commands such as `go build` and `go run`. To run the main server, use the `main.go` file in the root directory.

Amber also provides several helper tools, all of which are located in the `tools/` directory. Run the tool with the `-h` flag for detailed input descriptions.
- __Migration Runner__: Runs the data migrations. This will need to be run before using the server.
- __Admin Creator__: Creates a new admin user. This is necessary for creating the first user in the system.
- __Config Generator__: Generates a new config file, filling it with default values.
- __Key Generator__: Generates a new private/public key pair that can be used by the create token endpoint.

## Setup and Running

The following checklist should be followed when running the application for the first time in a new environment.
1. __Create the Config File__: This can be done manually or with the Config Generator tool. The name of the file should be `config.env.yml` where `env` is the desired environment.
1. __Add Keys__: Add the required key files to the static directory. They can either be generated using the Key Generator tool or provided from another source. Note: three key files are already present in the static directory. DO NOT USE THESE IN PRODUCTION. They are for testing only and are not safe since they are checked in using the source control.
1. __Run Migrations__: Use the Migration Runner tool to run the data migrations.
1. __Create the Max Admin__: Use the Admin Creator tool to create a max admin. Note: you should later change the password using the API since passing a password via the command line with the tool may not be safe.

Once the setup has been completed, the server can be run. Set the environment variable `CFG_ENV` to whatever environment you are running in. Its name should directly match the `env` part of the config file name created earlier. The default environment is "local".
