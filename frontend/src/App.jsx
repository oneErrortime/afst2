import React, { useState, useEffect } from 'react';
import './App.css';
import { AuthContext } from './context/AuthContext';
import { AuthProvider } from './context/AuthContext';
import Header from './components/Header';
import LoginForm from './components/LoginForm';
import Dashboard from './components/Dashboard';

function App() {
  return (
    <AuthProvider>
      <div className="App">
        <Header />
        <main className="main-content">
          <LoginForm />
          <Dashboard />
        </main>
      </div>
    </AuthProvider>
  );
}

export default App;