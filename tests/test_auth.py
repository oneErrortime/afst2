import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from app.main import app
from app.database import get_db, Base
from app import models  # Import models to register them with Base

# Create test database
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Create tables
Base.metadata.create_all(bind=engine)

def override_get_db():
    try:
        db = TestingSessionLocal()
        yield db
    finally:
        db.close()

app.dependency_overrides[get_db] = override_get_db

client = TestClient(app)

def test_user_registration():
    """Test user registration"""
    response = client.post(
        "/auth/register",
        json={
            "email": "test@example.com",
            "password": "password123"
        }
    )
    assert response.status_code == 200
    data = response.json()
    assert data["email"] == "test@example.com"
    assert "id" in data

def test_user_login():
    """Test user login"""
    # First register a user
    register_response = client.post(
        "/auth/register",
        json={
            "email": "login@example.com",
            "password": "password123"
        }
    )
    assert register_response.status_code == 200
    
    # Then try to login
    login_response = client.post(
        "/auth/login",
        json={
            "email": "login@example.com",
            "password": "password123"
        }
    )
    assert login_response.status_code == 200
    data = login_response.json()
    assert "access_token" in data
    assert data["token_type"] == "bearer"

def test_login_invalid_credentials():
    """Test login with invalid credentials"""
    login_response = client.post(
        "/auth/login",
        json={
            "email": "nonexistent@example.com",
            "password": "wrongpassword"
        }
    )
    assert login_response.status_code == 401

def test_protected_endpoint_without_token():
    """Test that protected endpoints require authentication"""
    # Try to create a book without authentication
    response = client.post(
        "/books/",
        json={
            "title": "Unauthorized Book",
            "author": "Unauthorized Author"
        }
    )
    assert response.status_code == 403  # Should require authentication

def test_protected_endpoint_with_token():
    """Test that protected endpoints work with valid token"""
    # Register and login to get a token
    register_response = client.post(
        "/auth/register",
        json={
            "email": "protected@example.com",
            "password": "password123"
        }
    )
    assert register_response.status_code == 200
    
    login_response = client.post(
        "/auth/login",
        json={
            "email": "protected@example.com",
            "password": "password123"
        }
    )
    assert login_response.status_code == 200
    token_data = login_response.json()
    access_token = token_data["access_token"]
    headers = {"Authorization": f"Bearer {access_token}"}
    
    # Now try to create a book with the token
    response = client.post(
        "/books/",
        json={
            "title": "Authorized Book",
            "author": "Authorized Author"
        },
        headers=headers
    )
    assert response.status_code == 200
    data = response.json()
    assert data["title"] == "Authorized Book"