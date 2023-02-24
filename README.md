# Word of wisdom
A tcp server & client implemented to demonstrate proof of work request communication and send a random quote after successful authorization

## Build commands
First run the server and then the client
### Server:
```shell
make build-server
make run-server
```

### Client:
```shell
make build-client
make run-client
```

## TODOs:
- Add configurations like listen tcp address to dockerfile or make docker-compose.yml files
- Pass a logger as configuration and use it for logging
- Implement other algorithms for POW
- Better error handling
- Implement Request/Response mechanism