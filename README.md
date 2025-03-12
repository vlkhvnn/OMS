# OMS

OMS (Order Management System) is a microservices-based application written in Go. It is designed to manage various aspects of order processing, including order creation, payments, stock management, and integration via a gateway. The project is organized into several services that each handle distinct responsibilities.

## Overview

OMS is a modular solution intended to simplify and streamline order management. The system is broken down into several components that communicate with each other to perform operations such as order processing, payment handling, inventory management, and service coordination via a gateway.

## Project Structure

The repository is organized as follows:

- **common/**  
  Shared utilities, models, and configurations used across different services.

- **gateway/**  
  An API gateway that routes requests to the appropriate services and handles cross-cutting concerns such as authentication and logging.

- **kitchen/**  
  The backend service responsible for processing order-related business logic.

- **orders/**  
  A dedicated service that handles the creation, tracking, and management of orders.

- **payments/**  
  A service to manage payment processing, integration with payment providers, and related financial transactions.

- **stock/**  
  A service responsible for managing inventory levels and stock updates.

- **docker-compose.yml**  
  Docker Compose configuration to run the services locally, simplifying development and testing.

- **go.work** & **go.work.sum**  
  Workspace configuration files for managing multi-module Go projects.

## Features

- **Microservices Architecture:**  
  Each service handles a distinct aspect of order management, allowing for scalable and maintainable code.

- **API Gateway:**  
  A central gateway facilitates communication between the client and microservices, managing routing and common functionality.

- **Containerization:**  
  Docker Compose configuration for easy local setup and orchestration of multiple services.

- **Go-Based Implementation:**  
  Written primarily in Go, taking advantage of its concurrency model and performance benefits.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/vlkhvnn/OMS.git
   cd OMS
2. **Running with Docker Compose:**

   Ensure Docker is running on your system. Then, run:

   ```bash
   docker-compose up
3. **Running Locally each service:**

   To run individual services navigate to the respective service folder (for example, the `orders/` directory) and execute:

   ```bash
   air
   ```

Make sure you have set up any necessary environment variables or configuration files as required by the service.

## Configuration

   Each service may have its own configuration file or use environment variables. Check the individual service folders for specific configuration details. Common settings might include:

   - Database connection strings
   - API keys for payment gateways
   - Service-specific settings

   Ensure these configurations are properly set up before running the services.

## How to work with this Project?

  First of all, lets create an order. To enable to just make a POST request to this url: http://localhost:8080/api/customers/2/orders with the JSON below:
  ```
  [
  {
    "id": "1",
    "quantity": 2
  },
  {
    "id": "2",
    "quantity": 1
  }
]
  ```
There will be a payment link in responce. Just fill the gaps in stripe with 4242...
Make sure that you run the stripe CLI and listen to the webhook!!!
```
  stripe listen --forward-to localhost:8081/webhook
```
### Monitoring
Jaeger UI link: http://localhost:16686
MongoDB UI link: http://localhost:8082
RabbitMQ UI link: http://localhost:15672
Consul UI link: http://localhost:8500
