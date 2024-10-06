# Payment Gateway Microservice

## Overview
This microservice integrates with multiple payment gateways for deposit and withdrawal operations.

## Getting Started
1. Clone the repository.
2. Run `go run ./cmd` to start the service on port 8080.

## API Endpoints
- **POST /deposit**: Deposit money.
- **POST /withdraw**: Withdraw money.
- **POST /callback/gateway_a**: Handle callbacks from gateway A.
- **POST /callback/gateway_b**: Handle callbacks from gateway B.

## Logging
The service uses standard logging for transactions and errors.