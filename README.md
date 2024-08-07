# Soldank++ Lobby
[Soldank++](https://github.com/nedik/soldank-plus-plus) JSON API based lobby server. Enables registering and discovering all registered [Soldank++](https://github.com/nedik/soldank-plus-plus) servers.

The lobby removes servers after 5 minutes unless they re-register.

## Endpoints
| HTTP Method | Endpoint                      | Returned type                          | Description                                                                                                |
| ----------- | ----------------------------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| GET         |  `/servers`                   | List<[Server](types/serverList.go#L3)> | Returns a list of all registered servers.                                                                  |
| GET         |  `/servers/:ip/:port`         | [Server](types/serverList.go#L3)       | Returns information about a server specified by the given `ip` and `port`.                                 |
| GET         |  `/servers/:ip/:port/players` | List<string>                           | Returns a list of players of a server specified by the given `ip` and `port`.                              |
| POST        |  `/servers`                   | Empty                                  | Registers a new server. Requires [RegisterServerInput](types/registerServerInput.go#L3) as request's body. |

## Environment variables
The program expects a file `app.env` to be created and filled with all the [required environment variables](initializers/loadEnv.go#L7). See [example.env](example.env)

## Dependencies
The project uses the following packages:
- [Gin Web Framework](https://github.com/gin-gonic/gin) - Framework that handles http connections
- [Viper](https://github.com/spf13/viper) - Helps define and read environment variables
- [TreeMap v2](https://github.com/igrmk/treemap) - Tree-based set data structure that helps manage the list of registered servers

## Building
To build a binary run:
```
go build
```
The above command should create `spp-lobby` file.

## Running
Run the lobby in release mode:
```
GIN_MODE=release ./spp-lobby
```

## Testing
To run all tests run command:
```
go test
```

