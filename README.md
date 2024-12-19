# GO-CQRS

This project is an implementation of the CQRS (Command Query Responsibility Segregation) pattern in Go. It provides a framework for handling commands and queries separately, allowing for more scalable and maintainable code.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Commands](#commands)
- [Queries](#queries)
- [Events](#events)
- [Aggregates](#aggregates)
- [Middlewares](#middlewares)
- [License](#license)

## Installation

To install the project, clone the repository and run the following command:

```sh
go mod tidy
```

## Usage

To run the project, use the following command:

```sh
go run examples/main.go
```

## Project Structure

The project is structured as follows:

- `examples/`: Contains example implementations of commands, queries, and event handlers.
- `pkg/`: Contains the core CQRS framework, including command and query buses, middlewares, and aggregate root implementations.

## Commands

Commands are used to change the state of the application. They are handled by command handlers. Example command handlers can be found in the `examples/application/command-handlers` directory.

Example command:

```go
package commands

type CreateUserCommand struct {
    Username    string `json:"username"`
    Password    string `json:"password"`
    DateOfBirth string `json:"date_of_birth"`
}
```

## Queries

Queries are used to retrieve data from the application. They are handled by query handlers. Example query handlers can be found in the `examples/application/query-handlers` directory.

Example query:

```go
package queries

type GetUserByIdQuery struct {
    ID string `json:"id"`
}
```

## Events

Events represent something that has happened in the application. They are handled by event handlers. Example events can be found in the `examples/domain/events` directory.

Example event:

```go
package domain_events

import (
    "log"
    "github.com/hashemihossein/GO-CQRS/examples/domain"
)

type UserCreatedEvent struct {
    User *domain.User
}

func (event *UserCreatedEvent) Handle() error {
    log.Println("User created: ", event.User)
    return nil
}
```

## Aggregates

Aggregates are the root entities that encapsulate the state and behavior of the application. They handle events and apply changes to the state. Example aggregates can be found in the `examples/domain` directory.

Example aggregate:

```go
package domain

import "github.com/hashemihossein/GO-CQRS/pkg/aggregate"

type User struct {
    *aggregate.AggregateRoot
    ID          string
    Username    string
    Password    string
    DateOfBirth string
}
```

## Middlewares

Middlewares are used to add additional behavior to commands, queries, and events. They can be used for logging, validation, etc. Example middlewares can be found in the `pkg/` directory.

Example middleware:

```go
package command

import "context"

type Middleware interface {
    Execute(ctx context.Context, cmd Command, next NextFunc) error
}

type NextFunc func(ctx context.Context, cmd Command) error
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
