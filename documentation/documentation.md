# Kebabot Docs

> [!WARNING]
> This is a personal project! The documentation here is for personal reference and will likely not be entirely thorough as this bot has never been intended to be used by anyone else. There may be missing features here and there, and missing documentation on certain features.

## Structure

- `cmd/kebabot/`: This is where `main.go` resides. This file loads the configuration, connects to the postgresql database and starts the bot.

- `internal/`: Where all of the bot's actual code resides, and what `main.go` calls.
    - `bot/`: This is where the bot is created and all of the commands and listeners are send to discord itself.
    - `commands/`: This is where every command of the discord bot is listed, with `commands.go` aggregating all of them together.
    - `listeners/`: This is where the message listeners are. Rather than being application commands, these listen for raw text.
    - `config/`: Where the configuration from the `.env` is pulled.
    - `database/`: WORK IN PROGRESS...

- `sql/`:
    - `migrations/`: WORK IN PROGRESS...
    - `queries/`: WORK IN PROGRESS...

`sqlc.yaml`: Where the configuration for sqlc resides.

`.env`: Where all of the environment variables for bot configuration are stored. `env.example` will be the one committed to git, where you can enter your own configurations.

## Getting Started

### Prerequisites

> [!WARNING]
> The database side of the bot has not yet been started. The database connection is still required to omit errors during startup. PostgreSQL is recommended, however any other database is usable. For small personal projects, you may wish to use the simpler SQLite. Make sure to update the `sqlc.yaml` and `.env` files accordingly.

- Go (>= 1.27.1 Recommended)
- An SQL Database (PostgreSQL Recommended, any others work too.)
- SQLC (>= 1.31.1 Recommended)
- Goose (>= 3.28.0 Recommended)

### Dependencies

Use the below command for all dependencies of the project. The latest version will automatically be installed.

```bash
go get github.com/jogo/godotenv
go get github.com/disgoorg/disgo
go get github.com/disgoorg/snowflake/v2
go get github.com/jackc/pgx/v5
```

### Setup

1. Create a database using your preferred SQL database.
2. Add the database URL to the `env.example` with format: `postgres://username:password@localhost:5432/kebabot`.

> [!IMPORTANT]
> Replace `username:password` with the username and password you have for your SQL database. Replace `localhost:5432` with whatever address and port you are running your database server off of. Replace kebabot with the name of the database. Replace `postgres` with whatever database you are using. Keep in mind, SQLite uses: `sqlite:///kebabot.db` format instead; no need for username, password, address or port!

3. Complete the rest of the env and rename `env.example` to `.env`.
4. Run the below commands.

 ```bash
go mod tidy
```

> [!NOTE]
> The below command is optional. It copies all dependencies into the `vendor/` folder, which has been excluded from git as dependencies can be pretty large in size at times, and it helps keep the repository cleaner. Keep in mind to run the command every so often to keep dependencies up to date.

```bash
go mod vendor
```

### Running

1. Go into the cmd/kebabot directory.
2. `go build`
3. Go back to the main directory.
4. `./cmd/kebabot/kebabot`
5. Within a few seconds, if no errors have occured. The bot should be up and running.

Below is a copyable command:

```bash
cd cmd/kebabot
go build
cd ../..
./cmd/kebabot/kebabot
```

- Q: Why are we moving back into the main folder and not just running it from the `cmd/kebabot/` folder?
- A: If you run it from the `cmd/kebabot/` folder, you'll likely get a "DISCORD_TOKEN is not set" error. This is because the `.env` is in the main folder, and cannot be found otherwise.

## Deep Dive Into The Project

> [!NOTE]
> This section will summarise each part of a file in the repository. If you click the headings, you will be redirected to a `.md` of that file which explains in further detail the code itself.

### [cmd/kebabot/main.go](main.go.md)

This is the file where the configuration for the bot is loaded, the connection to the database is established, and the bot is started.

#### [func main()](main.go.md#func-main)

The main function serves a small purpose. It calls the `run` function and handles any errors produced by it.

#### [func run()](main.go.md#func-run)

The run function actually runs the bot. If any error occurs, it returns an error to the main function. It loads the configuration from the `.env`; it registers the context so that the bot stops when a kill signal is recieved; it connects to the database, if applicable; it creates a new bot and `closeCtx` to make sure it shuts down gracefully; it then starts the bot and will forever wait for a kill command to stop the bot, which will call all deferred methods.

### [internal/bot/bot.go](bot.go.md)

This is where the functions for the bot to be created, started and closed are stored.

#### [type Bot struct](bot.go.md#type-bot-struct)

This struct holds the client which the bot uses to speak to Discord, and the configuration which is loaded from the `.env`.

#### [func New()](bot.go.md#func-new)

This function allows the creation of the Bot, all of its intents and the commands & handlers of said commands.

#### [func Start()](bot.go.md#func-start)

This function syncs the commands to Discord, and then opens the gateway connection, allowing the bot to go online.

#### [func Close()](bot.go.md#func-close)

This function simply closes the gateway connection gracefully.