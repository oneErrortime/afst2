# Backend API Deployment Guide

## Local Development

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Run the application:
```bash
python main.py
```

The API will be available at `http://localhost:8000`

## Deploy to Render.com

### 1. Prepare Your Render.com Account
- Sign up at [Render.com](https://render.com)
- Connect your GitHub/GitLab account
- Create a new Web Service

### 2. Configure Environment Variables
After connecting your repository, set these environment variables in the Render dashboard:

- `DATABASE_URL`: PostgreSQL connection string (you can create a free PostgreSQL instance on Render or use another provider like ElephantSQL)
- `SECRET_KEY`: A random secret key for JWT tokens (use a strong, unique value)

### 3. Deployment Settings
- Runtime: Python
- Build Command: `pip install -r requirements.txt`
- Start Command: `uvicorn main:app --host 0.0.0.0 --port $PORT`

## API Endpoints

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and get access token
- `GET /auth/me` - Get current user info (requires token)

### Example Requests

#### Register a new user:
```bash
curl -X POST "https://your-app-name.onrender.com/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### Login:
```bash
curl -X POST "https://your-app-name.onrender.com/auth/login" \
  -d "username=testuser&password=password123"
```

#### Access protected endpoint:
```bash
curl -X GET "https://your-app-name.onrender.com/auth/me" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Database Setup

For development, the app uses SQLite by default. For production, configure a PostgreSQL database:

1. Create a PostgreSQL database (Render offers free tier)
2. Set the `DATABASE_URL` environment variable to your database connection string
3. The format should be: `postgresql://username:password@host:port/database_name`

## Important Security Notes

- Change the default `SECRET_KEY` in production
- Use HTTPS in production
- Validate and sanitize all inputs
- Use strong passwords
- Regularly update dependencies