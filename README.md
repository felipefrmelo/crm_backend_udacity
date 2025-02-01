# CRM Backend

This project represents the backend of a Customer Relationship Management (CRM) web application. It provides a RESTful API to support the following functionalities:

- **Get a list of all customers**
- **Get data for a single customer**
- **Add a customer**
- **Update a customer's information**
- **Remove a customer**

The server is built with Go and leverages both Fiber and Gorilla Mux frameworks for routing and HTTP handling.

---

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
---

## Features

- **RESTful API Endpoints:**

  - `GET /customers` – Retrieve all customers.
  - `GET /customers/:id` – Retrieve a specific customer.
  - `POST /customers` – Add a new customer.
  - `PUT /customers/:id` – Update an existing customer's information.
  - `DELETE /customers/:id` – Remove a customer.

- **Data Models:**

  - `Customer`
  - `CreateCustomer`
  - `UpdateCustomer` (supports partial updates using pointer fields)

- **Adapters:**
  - Provides abstraction layers to work with different HTTP engines (Fiber and Gorilla Mux).

---

## Getting Started

These instructions will help you set up and run the CRM backend on your local machine for development and testing purposes.

### Prerequisites

- **Go:** Make sure you have Go installed (version specified in `go.mod` is `go 1.23.4` or above).
- **Git:** To clone the repository.

### Installation

1. **Clone the repository:**

   ```bash

   git clone https://github.com/felipefrmelo/crm_backend_udacity.git
   cd crm_backend_udacity
   ```

2. **Download the dependencies:**

   ```bash

   go mod download
   ```

### Usage

**Running the Server**
To run the server locally, execute:

```bash
go run cmd/main.go
or
go run cmd/main.go --server (fiber||gorilla)
```

By default, the server listens on port 3000. You can modify the listening address in your configuration if needed.
