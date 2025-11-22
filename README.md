# Library Management API

A comprehensive library management system built with FastAPI that provides RESTful API endpoints for managing books, readers, and book borrowing processes. The system implements JWT-based authentication for librarians and includes comprehensive business logic for library operations.




ДВА ОСНОВНЫХ ИНСТАНСА:
https://ooouwuooo.github.io/afst/  - это фронтенд, тут ты можешь протестировать функционал проекта. 
https://afst.onrender.com/  - это временный бэкенд, тут ты можешь самостоятельно сгенерировать и задеплоить проект. Поддомен выбирай какой захочешь, не обязателньо afst. 

**Регистрируешься в render.com или используешь мой временный бэкенд, далее вводишь его во фронтенд-панели и тестируешь функциональность софта**





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

After starting the application, you can use the following pre-created admin account to access the system:

**Admin Credentials:**
- Email: admin@example.com
- Password: admin123

Alternatively, you can register a new user through the `/auth/register` endpoint. Once registered, you can log in via `/auth/login` to obtain a JWT token for accessing protected endpoints.

### Web Dashboard

The application includes a web-based dashboard for easier interaction with the API. After starting the application, navigate to `http://localhost:8000/dashboard` to access the user-friendly interface for managing books, readers, and borrowing operations.

## Using the GitPages Panel

The frontend for this application is deployed to GitHub Pages at `https://ooouwuooo.github.io/afst/`. To use the GitPages panel with your backend:

1. Navigate to `https://ooouwuooo.github.io/afst/`
2. In the API Configuration section, enter your backend API URL (e.g., `https://your-backend-url.com` or `http://localhost:8000` for local development)
3. Click "Save URL"
4. Use the login credentials provided below to access the admin panel:
   - Email: admin@example.com
   - Password: admin123
5. After successful login, you'll have full access to the library management system

**Note:** For the GitPages version to work properly, your backend API must be accessible from the internet and have proper CORS configuration to allow requests from the GitHub Pages domain.

### Local Development with GitPages Frontend

If you want to run the frontend locally for development:
```bash
python local_test_server.py
```

This will serve the frontend files from the `/docs` directory at `http://localhost:8000`, where you can connect to your backend API.

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

## Deployment to Render.com

This application can be easily deployed to Render.com. Follow these steps:

1. Make sure your code is pushed to a Git repository (GitHub, GitLab, or Bitbucket)
2. Go to https://dashboard.render.com/ and create a new Web Service
3. Connect your Git repository
4. Use the provided `Dockerfile` and `.render.yaml` configuration files
5. Set the required environment variables:
   - `DATABASE_URL`: Your PostgreSQL database connection string
   - `SECRET_KEY`: A strong, random secret key for JWT tokens
   - `ALGORITHM`: (optional) Defaults to HS256
   - `ACCESS_TOKEN_EXPIRE_MINUTES`: (optional) Defaults to 30

For detailed instructions, see the `DEPLOYMENT_RENDER.md` file in the repository.

**Note**: For production use, make sure to:
- Use a PostgreSQL database instead of SQLite
- Set a strong, unique `SECRET_KEY`
- Configure proper SSL certificates if using custom domains
