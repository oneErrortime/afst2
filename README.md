# Library Management API

A comprehensive library management system built with FastAPI that provides RESTful API endpoints for managing books, readers, and book borrowing processes. The system implements JWT-based authentication for librarians and includes comprehensive business logic for library operations.

## Features

- JWT-based authentication for librarians
- CRUD operations for books, readers, and users
- Book borrowing and return functionality with business rules
- Comprehensive validation and error handling
- Database migrations with Alembic
- Unit tests for business logic
- Support for book descriptions (added via migration)
- Web-based dashboard for user-friendly interaction

## Project Structure

```
library_api/
├── alembic/
│   ├── versions/
│   └── env.py
├── app/
│   ├── __init__.py
│   ├── main.py
│   ├── config.py
│   ├── database.py
│   ├── models/
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── book.py
│   │   ├── reader.py
│   │   └── borrow.py
│   ├── schemas/
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── book.py
│   │   ├── reader.py
│   │   └── borrow.py
│   ├── auth/
│   │   ├── __init__.py
│   │   └── jwt_handler.py
│   └── api/
│       ├── __init__.py
│       ├── auth.py
│       ├── books.py
│       ├── readers.py
│       └── borrows.py
├── alembic.ini
├── requirements.txt
├── tests/
└── README.md
```

## Database Schema Design

The system implements the following database structure:

1. **users table**: Stores librarian information with hashed passwords and JWT authentication
2. **books table**: Contains book information including title, author, publication year, ISBN, available copies, and description
3. **readers table**: Stores reader information with name and email
4. **borrows table**: Tracks book borrowing records with borrow/return dates

The design ensures data integrity through foreign key relationships and unique constraints where appropriate.

## Business Logic Implementation

The system implements the following business rules:

1. **Book Availability**: Books can only be borrowed if available copies > 0. When borrowed, the available copies decrease by 1.
2. **Reader Limit**: A single reader cannot borrow more than 3 books simultaneously. This is validated during the borrowing process.
3. **Return Validation**: Books can only be returned if they were actually borrowed by the specified reader and haven't been returned yet.
4. **Unique ISBN**: Each book must have a unique ISBN if provided.
5. **Data Validation**: Comprehensive validation for all input data.

These rules are enforced through validation logic in the API endpoints and database constraints where appropriate.

## Authentication Implementation

The system implements JWT-based authentication using the following components:

- **python-jose**: For JWT token creation and validation
- **passlib[bcrypt]**: For password hashing and verification
- **Security dependencies**: FastAPI dependencies that verify token validity

The authentication flow includes:
1. User registration with email and password (password is hashed)
2. Login endpoint that validates credentials and returns JWT token
3. Protected endpoints that require valid JWT tokens
4. Token expiration and refresh mechanisms

Protected endpoints include all book and reader management operations, as well as borrowing and returning functionality. The authentication ensures that only authenticated librarians can perform administrative tasks.

## Data Validation

The system implements comprehensive data validation:

- **Book Validation**: Title and author required, length limits, year validation, ISBN format, copies validation, description length limits
- **Reader Validation**: Name and email required, character validation, email format, length limits
- **Borrow Validation**: Positive integer validation for book and reader IDs

## Migration Strategy

The system uses Alembic for database migrations with a two-step approach:
1. Initial migration creates all tables except the description field in books
2. Second migration adds the description field to the books table

## Getting Started

### Prerequisites

- Python 3.8+
- PostgreSQL (or SQLite for development)

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

3. Set up environment variables:
   ```bash
   export SECRET_KEY="your-secret-key-here"
   export ALGORITHM="HS256"
   export ACCESS_TOKEN_EXPIRE_MINUTES=30
   ```

4. Set up the database:
   ```bash
   # For PostgreSQL, update DATABASE_URL in config.py
   # Run migrations
   alembic upgrade head
   ```

5. Run the application:
   ```bash
   uvicorn app.main:app --reload
   ```

### Initial Setup

After starting the application, you can register a new user through the `/auth/register` endpoint. Once registered, you can log in via `/auth/login` to obtain a JWT token for accessing protected endpoints.

### Web Dashboard

The application includes a web-based dashboard for easier interaction with the API. After starting the application, navigate to `http://localhost:8000/dashboard` to access the user-friendly interface for managing books, readers, and borrowing operations.

## Testing

Run the test suite with:
```bash
pytest tests/
```

The tests include validation of business logic, authentication, and API endpoints.

## Deployment

This application can be deployed to various platforms including Render.com, Heroku, or any cloud platform that supports Python applications.

For production use, make sure to:
- Use a PostgreSQL database instead of SQLite
- Set a strong, unique `SECRET_KEY`
- Configure proper SSL certificates if using custom domains
- Implement proper logging and monitoring