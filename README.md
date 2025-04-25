# CodePush Server

A Go-based server for managing CodePush deployments with support for multiple databases and containerized deployment.

## Features

- User authentication and authorization
- App ID and token management
- Database-agnostic architecture
- Docker and Kubernetes support
- Environment-based configuration

## Prerequisites

- Go 1.21 or later
- Docker (for containerization)
- Kubernetes (for deployment)
- PostgreSQL/MySQL (for database)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/piyushsharma67/codepushserver.git
cd codepushserver
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

4. Update the `.env` file with your configuration.

## Running Locally

```bash
go run main.go
```

## Docker Deployment

1. Build the Docker image:
```bash
docker build -t codepush-server .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env codepush-server
```

## Kubernetes Deployment

1. Create Kubernetes secrets and config maps:
```bash
kubectl create secret generic codepush-secrets --from-env-file=.env
kubectl create configmap codepush-config --from-env-file=.env
```

2. Apply the Kubernetes deployment:
```bash
kubectl apply -f kubernetes/deployment.yaml
```

## API Endpoints

### Authentication
- POST `/auth/register` - Register a new user
- POST `/auth/login` - Login and get JWT token

### Protected Routes
- GET `/api/profile` - Get user profile
- POST `/api/app` - Create a new app
- GET `/api/app/:id` - Get app details
- PUT `/api/app/:id` - Update app
- DELETE `/api/app/:id` - Delete app

## Database Support

The server supports multiple databases through a common interface. Currently supported:
- PostgreSQL
- MySQL

## License

MIT 