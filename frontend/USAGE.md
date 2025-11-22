# Using the Library Management Frontend with the Backend

This document explains how to connect the frontend to your deployed backend API.

## Backend API Endpoint Configuration

The frontend is configured to connect to the backend API at `/api/` by default. During development, this is proxied to `http://localhost:8000`.

When your backend is deployed (e.g., on Render), you'll need to update the API endpoint.

### Option 1: Update the Vite proxy configuration (for development)

Modify `vite.config.js`:
```js
export default defineConfig({
  // ... other config
  server: {
    host: true,
    port: 3000,
    proxy: {
      '/api': {
        target: 'https://your-backend-app-name.onrender.com', // Replace with your actual backend URL
        changeOrigin: true,
      }
    }
  },
  // ... rest of config
})
```

### Option 2: Update the fetch calls directly

In production, you can update the API URLs in each component to point to your backend:

For example, in `Login.jsx`:
```js
const response = await fetch('https://your-backend-app-name.onrender.com/api/auth/login', {
  // ... rest of the request
});
```

## Authentication Flow

1. The user logs in via the login page with their username and password
2. The backend returns a JWT token on successful authentication
3. The frontend stores the token in localStorage
4. The token is included in the Authorization header for all subsequent API requests
5. The user is redirected to the dashboard

## API Endpoints Used

The frontend communicates with the following backend endpoints:

- `POST /api/auth/login` - User authentication
- `GET /api/books/` - Retrieve all books
- `POST /api/books/` - Create a new book
- `PUT /api/books/{id}` - Update a book
- `DELETE /api/books/{id}` - Delete a book
- `GET /api/users/` - Retrieve all users
- `POST /api/users/` - Create a new user
- `PUT /api/users/{id}` - Update a user
- `DELETE /api/users/{id}` - Delete a user
- `GET /api/borrows/` - Retrieve all borrow records
- `POST /api/borrows/` - Create a new borrow record
- `POST /api/borrows/{id}/return` - Mark a book as returned

## Environment Configuration for Production

When deploying to GitHub Pages, you can use environment variables to configure the API endpoint:

1. Create a `.env.production` file:
```
VITE_API_URL=https://your-backend-app-name.onrender.com
```

2. Update your fetch calls to use the environment variable:
```js
const response = await fetch(`${import.meta.env.VITE_API_URL}/api/auth/login`, {
  // ... rest of the request
});
```

3. Update your `vite.config.js` to include environment variables:
```js
export default defineConfig({
  // ... other config
  define: {
    'process.env': process.env
  }
  // ... rest of config
})
```

## Deploying to GitHub Pages

1. Make sure your repository is set up for GitHub Pages in the repository settings
2. Push your code to the main branch
3. The GitHub Actions workflow will automatically build and deploy the frontend
4. Access your frontend at `https://your-username.github.io/your-repo-name`

## Troubleshooting

### CORS Issues
If you encounter CORS errors, make sure your backend is configured to allow requests from your frontend domain.

### API Connection Issues
- Check that your backend is running and accessible
- Verify the API endpoint URLs in your frontend components
- Check browser developer tools for specific error messages

### Authentication Issues
- Verify that the login credentials are correct
- Check that the JWT token is being stored properly in localStorage
- Ensure the Authorization header is being sent with requests