# Library Management API

A comprehensive library management system built with FastAPI that provides RESTful API endpoints for managing books, readers, and book borrowing processes. The system implements JWT-based authentication for librarians and includes comprehensive business logic for library operations.

## Features

- JWT-based authentication for librarians
- CRUD operations for books and readers
- Book borrowing and return functionality with business rules
- Comprehensive validation and error handling
- Database migrations with Alembic
- Unit tests for business logic

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
2. **books table**: Contains book information including title, author, publication year, ISBN, and available copies
3. **readers table**: Stores reader information with name and email
4. **borrows table**: Tracks book borrowing records with borrow/return dates

The design ensures data integrity through foreign key relationships and unique constraints where appropriate.

## Business Logic Implementation

The system implements the following business rules:

1. **Book Availability**: Books can only be borrowed if available copies > 0. When borrowed, the available copies decrease by 1.
2. **Reader Limit**: A single reader cannot borrow more than 3 books simultaneously. This is validated during the borrowing process.
3. **Return Validation**: Books can only be returned if they were actually borrowed by the specified reader and haven't been returned yet.

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

After starting the application, you'll need to register the first librarian user through the `/auth/register` endpoint. Once registered, you can log in via `/auth/login` to obtain a JWT token for accessing protected endpoints.

### Web Dashboard

The application includes a web-based dashboard for easier interaction with the API. After starting the application, navigate to `http://localhost:8000/dashboard` to access the user-friendly interface for managing books, readers, and borrowing operations.

## GitHub Pages Demo

A frontend-only version of the dashboard is available for deployment on GitHub Pages. This static version can connect to any backend API instance and provides the same functionality without requiring server-side deployment.

### Deployment to GitHub Pages

There are two ways to deploy the frontend to GitHub Pages:

**Method 1: Using GitHub Actions (Recommended)**
1. The repository already includes a GitHub Actions workflow in `.github/workflows/deploy.yml`
2. Push your code to the main branch
3. Go to repository Settings -> Pages
4. Select source as "Deploy from a branch" and choose `gh-pages` branch
5. The site will be available at `https://<username>.github.io/<repository-name>`

**Method 2: Manual Deployment**
1. Push the code in the `/docs` directory to your GitHub repository
2. Go to repository Settings -> Pages
3. Select source as your main branch with `/docs` folder
4. The site will be available at `https://<username>.github.io/<repository-name>`

The GitHub Pages version includes:
- Responsive design with Bootstrap 5
- Modern UI with Font Awesome icons
- Dashboard statistics
- API URL configuration
- Full library management functionality

### Local Testing

To test the GitHub Pages version locally, use the provided script:

```bash
python local_test_server.py
```

This will start a local server at http://localhost:8000 to serve the frontend files from the `/docs` directory.

## Additional Feature Proposal

An additional feature that could be implemented is **overdue book notifications**. This would involve:

1. Adding a `due_date` field to the borrow records
2. Implementing a background task that checks for overdue books
3. Sending notifications (email or SMS) to readers with overdue books
4. Optionally implementing a fine system for overdue returns

This would improve the library management system by helping track and manage overdue books more effectively.

## Testing

Run the test suite with:
```bash
pytest tests/
```

The tests include validation of business logic, authentication, and API endpoints.