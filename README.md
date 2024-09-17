# Note Service

This is a simple REST API service for managing notes with spell-checking functionality.

## Features

- Create notes with automatic spell-checking
- Retrieve notes for authenticated users
- User authentication
- Dockerized application

## Prerequisites

- Go 1.16+
- Docker and Docker Compose
- Postman (for testing API endpoints)

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/note-service.git
   cd note-service
   ```

2. Build and run the application using Docker Compose:
   ```
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:8080`

## API Endpoints

### Create a Note

```
POST /notes
Authorization: token1
Content-Type: application/json

{
  "title": "My Note",
  "content": "This is the content of my note."
}
```

### Get All Notes

```
GET /notes
Authorization: token1
```

## Running Tests

To run the automated tests:

```
go test ./...
```

## Postman Collection

A Postman collection is included in the repository for easy testing of the API endpoints. Import the `todo_list.postman_collection.json` file into Postman to use it.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
