# Deployment to Render.com

This guide explains how to deploy this FastAPI application to Render.com.

## Prerequisites

- A Render.com account
- Your project pushed to a Git repository (GitHub, GitLab, or Bitbucket)

## Steps to Deploy

1. **Prepare your repository**:
   - Make sure your code is pushed to a Git repository
   - The repository should include:
     - `Dockerfile` - containerization instructions
     - `.render.yaml` - Render deployment configuration
     - `requirements.txt` - Python dependencies
     - All application source code

2. **Create a new Web Service on Render**:
   - Go to https://dashboard.render.com/
   - Click "New +" and select "Web Service"
   - Connect your Git repository
   - Select the branch you want to deploy

3. **Configure the deployment**:
   - Service name: Choose any name (e.g., "library-management-api")
   - Environment: Docker
   - Branch: main (or your default branch)
   - The build and start commands are already specified in `.render.yaml`

4. **Set environment variables**:
   - DATABASE_URL: Your database connection string (for PostgreSQL, it would be something like `postgresql://username:password@host:port/dbname`)
   - SECRET_KEY: A strong, random secret key for JWT tokens
   - ALGORITHM: (optional) Defaults to HS256
   - ACCESS_TOKEN_EXPIRE_MINUTES: (optional) Defaults to 30

5. **Deploy**:
   - Click "Create Web Service"
   - Render will build and deploy your application automatically

## Important Notes

- The application will be available at the URL provided by Render after successful deployment
- The database URL should be updated to use a PostgreSQL database for production use
- Make sure to use strong, unique values for SECRET_KEY in production
- The application uses port specified by the `$PORT` environment variable on Render