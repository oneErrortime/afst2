# Library Management System Frontend

This is the frontend for the library management system. It's built with React and Vite and connects to the backend API to manage books, users, and book borrowings.

## Features

- User authentication and authorization
- Book management (CRUD operations)
- User management (CRUD operations)
- Book borrowing and return functionality
- Search and filter capabilities
- Responsive design for various screen sizes

## How to Run Locally

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open your browser to `http://localhost:3000`

## How to Build for Production

To build the application for production deployment:

```bash
npm run build
```

This will create a `dist` folder with the production-ready files.

## Deployment to GitHub Pages

1. Install the gh-pages package:
   ```bash
   npm install --save-dev gh-pages
   ```

2. Add the following scripts to your `package.json`:
   ```json
   {
     "scripts": {
       "predeploy": "npm run build",
       "deploy": "gh-pages -d dist"
     }
   }
   ```

3. Update your `vite.config.js` to set the correct base path:
   ```js
   export default {
     // ... other config
     base: '/your-repo-name/'
   }
   ```

4. Deploy to GitHub Pages:
   ```bash
   npm run deploy
   ```

## API Configuration

The frontend connects to the backend API at `/api/` by default. During development, the Vite server proxies API requests to `http://localhost:8000` (the default backend port).

For production deployment, you'll need to configure the API endpoint in the components or via environment variables.

## Technologies Used

- React 18
- React Router DOM
- Vite
- CSS Modules