import React from 'react';
import { useAuth } from '../context/AuthContext';

const Header = () => {
  const { user, logout, isAuthenticated } = useAuth();

  return (
    <header className="header">
      <div className="container">
        <h1 className="logo">Library Management System</h1>
        {isAuthenticated && (
          <div className="user-info">
            <span>Welcome, {user?.email}</span>
            <button className="btn btn-danger" onClick={logout}>
              Logout
            </button>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;