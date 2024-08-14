# Go Scaffold: A Templating Framework

Go Scaffold is a templating framework designed to generate generic code based on a provided configuration JSON file.

## Features

Go Scaffold can currently generate the following:

- **App Config Files**: Basic configuration files for your application.
- **Database Structs**: Structs for PostgreSQL, MySQL, and SQLite.
- **MongoDB Connection**: Code to establish a connection with MongoDB.
- **CRUD Code Generation**: CRUD operations using models defined in the template config JSON file, with support for:
  - **GORM and Fiber** (mostly stable) with sample testing
  - **GORM and Echo** (early stages) with sample testing
- **Jaeger Tracing**: Integration with OpenTelemetry (otel) for tracing.
- **Dockerfile**: A basic Dockerfile for containerization.
- **RabbitMQ Communication**: Basic consumer and publisher for RabbitMQ in Go.
- **RPC Structure**: Basic structure for RPC communication.
- **Additional Features**: Other features that are not yet documented.

## Quick Start

1. Create a new directory and navigate into it:
    ```bash
    mkdir frame-play && cd frame-play
    ```

2. Copy `config.json` into the `frame-play` directory and set the `project_name` in the `config.json` file.

3. Generate the basic structure:
    ```bash
    ./frame basic
    ```

4. Use the help command to see available options:
    ```bash
    ./frame help
    ```

## Example Usage

To generate and run a basic project with Fiber:

1. Generate the Fiber project:
    ```bash
    ./frame genfiber
    ```

2. Build and run the project:
    ```bash
    go mod tidy
    go run main.go dev
    ```
