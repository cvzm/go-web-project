# Go Web Service Demo Project

This demonstration project showcases a Go-based web service, highlighting best practices and architecture for modern Go web applications. It uses an event processing system as a concrete example to illustrate the construction of a fully functional and well-structured Go web service.

## Project Overview

This demo project centers around event handling, demonstrating the complete workflow from database operations to API design and application bootstrapping. Through this practical use case, developers can gain valuable insights into building scalable, high-performance, and maintainable web services in Go.

## Key Features

- RESTful API implemented using the Echo framework, showcasing CRUD operations for events
- PostgreSQL database operations using GORM for persistent storage of event data
- Domain-Driven Design (DDD) principles applied to clearly separate business logic from infrastructure code
- Graceful shutdown mechanism to ensure proper request handling during service termination
- Flexible configuration management through environment variables, adaptable to different deployment environments
- Support for database read replicas, optimizing query performance for high-concurrency scenarios

## Project Structure

- `adapter/storage`: Database connection and configuration, including setup for primary and read-only replicas
- `api`: Event-related API controllers and route setup
- `bootstrap`: Application initialization and configuration management
- `domain`: Event domain models and interface definitions
- `usecase`: Application use cases, implementing business logic for event operations
- `repository`: Data repositories, implementing database operations for events
- `main.go`: Application entry point for launching the web service
