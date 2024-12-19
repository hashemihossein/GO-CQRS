# GO-CQRS

This package is an implementation of the CQRS (Command Query Responsibility Segregation) pattern in Go. It provides a framework for handling commands and queries separately, allowing for more scalable and maintainable code.

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

You can download this package, use the following command:

```sh
go get github.com/hashemihossein/GO-CQRS
```

## Usage

If you clone the project you can run the examples, use the following command:

```sh
go run examples/main.go
```

## Commands

Commands are used to change the state of the application. They are handled by command handlers. Example command handlers can be found in the `examples/application/command-handlers` directory.

## Queries

Queries are used to retrieve data from the application. They are handled by query handlers. Example query handlers can be found in the `examples/application/query-handlers` directory.

## Events

Events represent something that has happened in the application. They are handled by event handlers. Example events can be found in the `examples/domain/events` directory.

## Aggregates

Aggregates are the root entities that encapsulate the state and behavior of the application. They handle events and apply changes to the state. Example aggregates can be found in the `examples/domain` directory.

## Middlewares

Middlewares are used to add additional behavior to commands, queries, and events. They can be used for logging, validation, etc. Example middlewares can be found in the `pkg/` directory.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
