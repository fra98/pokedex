# Pokedex

A fun REST API that provides Pokémon information with translated descriptions.
This API serves basic Pokémon data and offers descriptions translated in Yoda-speak or Shakespearean English based on the Pokémon's habitat and legendary status.

## Features

- Get basic information about a Pokémon including name, English description, habitat, and legendary status
- Get Pokémon information with descriptions translated into Yoda-speak (for Pokémons with habitat `cave` or legendary) or Shakespearean English (for all others)
- Graceful handling of API rate limits: if the translation can not be provided (e.g., rate limit exceeded), the original description is returned
- In-memory caching to reduce external API calls and improve performance
- Command-line flags are used to configure the server address and port, timeout, enable/disable cache, etc.

## Requirements

If you want to build and run the application locally, you will need the following:

- Git
- [httpie](https://httpie.io/docs#installation)
- [Go](https://golang.org/doc/install) >=1.24
- [make](https://www.gnu.org/software/make/) utility

If you prefer to use [Docker](https://docs.docker.com/get-started/get-docker/), you need a working Docker installation.

Warning: the following instructions assume you have a Unix-like environment (e.g., Linux, macOS). If you are using Windows, you may need to adapt the commands accordingly.

## Installation

### Clone the repository

```bash
git clone https://github.com/fra98/pokedex.git
cd pokedex
```

## Building and running the API server

### Option 1) Binary execution

#### Build binary

```bash
make build
```

#### Run binary

```bash
chmod +x ./bin/pokedex
./bin/pokedex
```

By default, the API runs locally port 8080. You can access it at <http://localhost:8080>.
If you want to run the API on a different host or port use the `--address` flag (e.g., `./bin/pokedex --address :9090`).

### Option 2) Container execution

If you prefer to use Docker, you can build and run the API in a container, or pull the pre-built image at `ghcr.io/fra98/pokedex:latest`.

#### Build Docker image

```bash
make build-docker
```

#### Run Docker container

```bash
make run-docker
```

Or, if you want to pass arguments to the container:

```bash
docker run -p 8080:8080 pokedex:latest
```

## API Endpoints

### 1. Get Basic Pokémon Information

```text
GET /pokemon/<pokemon-name>
```

Example:

```bash
http GET http://localhost:8080/pokemon/mewtwo
```

Response:

```json
{
    "name": "mewtwo",
    "description": "It was created by a scientist after years of horrific gene splicing and DNA engineering experiments.",
    "habitat": "rare",
    "isLegendary": true
}
```

### 2. Get Translated Pokémon Description

```text
GET /pokemon/translated/<pokemon-name>
```

Example:

```bash
http GET http://localhost:8080/pokemon/translated/mewtwo
```

Response:

```json
{
    "name": "mewtwo",
    "description": "Created by a scientist after years of horrific gene splicing and dna engineering experiments, it was.",
    "habitat": "rare",
    "isLegendary": true
}
```

## Manual Testing

You can test the API using a web browser or tools like Postman, curl, httpie, etc.:

1. Basic information:

   ```text
   http://localhost:8080/pokemon/pikachu
   ```

2. Translated information:

   ```text
   http://localhost:8080/pokemon/translated/pikachu
   ```

## Design Decisions

### Project Structure

The project layout follows the standard Go project structure, with the main components organized as follows:

```text
pokedex
├─ Dockerfile           # Docker configuration
├─ Makefile             # build and run commands
├─ README.md            # documentation
├─ cmd                  # entry point
└─ pkg
   ├─ api               # API handlers and routes
   ├─ client            # external API clients
   │  ├─ pokeapi        # - PokeAPI client
   │  └─ translator     # - FunTranslations API client
   ├─ consts            # common constants
   ├─ errors            # custom errors
   ├─ flags             # command-line flags
   ├─ models            # shared data models
   ├─ server            # server configuration
   └─ service           # business logic
```

### Components

The application is built with a layered architecture where each component has a specific responsibility and separation of concerns. The main layers are:

- **API Layer** (`api/`): handles HTTP requests/responses, input validation, routing, and middleware.
It is implemented with a basic http server, using Gin web framework.
- **Service Layer** (`service/`): handle business logic for Pokemon data retrieval and translation
- **Client Layer** (`client/`): encapsulates external API communication.
It provides clients to interact with PokeAPI and FunTranslations APIs.

Other components include:

- **Models** (`models/`): defines data structures shared across layers
- **Errors** (`errors/`): custom error types

### Design decisions

#### Interface-driven development

The application uses interfaces to define contracts between layers:

```go
// Pokemon service interface
type Pokemon interface {
    GetPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error)
    GetTranslatedPokemonInfo(ctx context.Context, name string) (*models.PokemonResponse, error)
}

// PokeAPI client interface
type Client interface {
    GetPokemonSpecies(ctx context.Context, name string) (*PokemonSpecies, error)
}

// Translation client interface
type Client interface {
    Translate(ctx context.Context, text, translationType string) (string, error)
}
```

This approach allows for easy testing (no mocking), clear boundaries between components, and flexibility to swap implementations (e.g., swap within the cached client and non-cached client).

#### HTTP Stubbing for testing

The testing strategy uses HTTP stubbing thanks to the `httptest` package.
This allows testing the API handlers and service layer without making actual HTTP requests to external APIs.
By stubbing at the HTTP level rather than mocking interfaces, we ensure that our tests verify that our client code correctly interacts with the external APIs.

The benefits of this approach are:

- tests the actual HTTP client code
- verifies correct request construction
- resiliency to refactoring (e.g., no need to update mocks when changing implementations)
- tests HTTP error handling and response parsing
- more close to real-world scenarios (testing edge cases like rate limiting)

#### Caching

The application uses an in-memory cache to store Pokémon data and translations, using the `go-cache` library.
The cache reduces the number of external API calls and improves performance.
Caching is implemented using a decorator pattern, where the cache client wraps the actual
client and intercepts requests to check if the data is already cached.
This is achieved easily thanks to the interface-based design, as both the cache and non-cache clients implement the same interface, making it easy to conditionally enable/disable caching.
The cache expiring timeout and cleanup interval are configurable via command-line flags.

#### Stateless and containerizable

The application is designed to be stateless, making it easy to scale horizontally and deploy in containerized environments like Docker or Kubernetes, thanks to small image size and minimal dependencies.

### Considerations for Production

For a production environment, the following enhancements could be made:

#### Performance and Scalability

1. **External Caching Layer**: implement a dedicated caching service (e.g., *Redis*, *Memcached*) to obtain the following benefits (at the cost of additional complexity and external dependencies):
    - Distributed caching for multiple instances
    - Persistence across application restarts
    - Built-in TTL and eviction policies
    - Can be shared across services
2. **Horizontal Scaling**: deploy multiple replicas (horizontal scaling) to handle increased load. This is easier to achieve since the application is stateless.
3. **Load Balancing**: Deploy multiple instances of the API behind a load balancer for better performance and availability

#### Security and Reliability

1. **Rate Limiting**: Add rate limiting on the API to prevent abuse and ensure fair usage. It can be implemented using a middleware or by hiding the API behind a dedicated API Gateway or reverse proxy.
2. **Input Validation**: More robust input validation
3. **Authentication**: Add API authentication for secured endpoints
4. **TLS**: Enable HTTPS for secure communication
5. **Health Checks**: Add health check endpoints for container orchestration systems
6. **Retry Logic**: Implement retry logic for external API calls to handle transient failures

#### Observability and Monitoring

1. **Logging**: Add structured logging for better observability
2. **Metrics**: Implement Prometheus metrics for monitoring

#### Other Improvements

1. **API Documentation**: Add Swagger/OpenAPI documentation
2. **Kubernetes Deployment**: Create a Kubernetes deployment manifest for easy scaling and management
3. **API Versioning**: Implement API versioning to maintain backward compatibility

#### Testing

1. **E2E Tests**: Add real E2E integration tests to cover the whole application (run manually or using a CI/CD pipeline), simulating real-world scenarios

## Dependencies

This project uses the following external APIs:

- [PokeAPI](https://pokeapi.co/) - For Pokémon data
- [FunTranslations API](https://funtranslations.com/) - For text translations

And the following Go libraries:

- [gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [go-cache](github.com/patrickmn/go-cache) - In-memory cache
- [testify](github.com/stretchr/testify) - Testing utilities
- [pflag](github.com/spf13/pflag) - Command-line flag parsing
