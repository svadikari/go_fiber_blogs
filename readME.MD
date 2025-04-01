# go_fiber_blogs

A Go web application for blog management built with Fiber framework.

## Features

- RESTful API for blog management (CRUD operations)
- User authentication
- Swagger API documentation
- HTML templating with layouts
- Database integration


## Tech Stack

| Tech            | Usage                                                     |
|-----------------|-----------------------------------------------------------|
| Go              | [SDK](https://golang.org/)                                |
| Fiber           | [Web Framework](https://docs.gofiber.io/)                 |
| golang-jwt      | JWT token library                                         |
| joho/godotenv   | To read environment variable/file                         |
| Swagger         | [API Documentation](https://github.com/gofiber/swagger)   |
| swaggo/swag     | Swagger cli to generat swagger documentation              |
| x/crypto        | Encryption/Decrytion library                              |
| driver/postgres | PostGresql driver                                         |
| gorm            | Go ORM for SQL drivers                                    |
| HTML            | [HTML Templates](https://docs.gofiber.io/guide/templates) |

## Getting Started

### Prerequisites

- Go 1.16+
- Database (configuration in environment variables)

### Installation

1. Clone the repository
  ```bash
  git clone https://github.com/svadikari/go_fiber_blogs.git
  cd go_fiber_blogs
  ```
2. Install dependencies
```bash
go mod download
```

3. Configure environment variables (create a `.env` file)
```
APPLICATION_PORT=8080
DB_CONNECTION_STRING=your_db_connection_string
```

4. Run the application
```bash
go run src/main.go
```


The application will be available at `http://localhost:8080`

## API Documentation

Swagger UI is available at `http://localhost:8080/swagger/` endpoint.


## Web Application

Web UI is available at `http://localhost:8080/login/` endpoint.

## Project Structure

- `/src` - Source code
    - `/database` - Database connection and models
    - `/docs` - Swagger documentation
    - `/middleware` - Application middleware
    - `/routers` - Route definitions
    - `/utils` - Utility functions
- `/templates` - HTML templates

## License

Open

## Contact

Shyam - shyam@shyam.com