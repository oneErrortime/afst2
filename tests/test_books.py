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

def test_create_book():
    response = client.post(
        "/books/",
        json={
            "title": "Test Book",
            "author": "Test Author",
            "year": 2023,
            "isbn": "978-1234567890",
            "copies": 3
        }
    )
    assert response.status_code == 200
    data = response.json()
    assert data["title"] == "Test Book"
    assert data["author"] == "Test Author"
    assert data["copies"] == 3

def test_get_books():
    response = client.get("/books/")
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)

def test_get_book():
    # First create a book
    create_response = client.post(
        "/books/",
        json={
            "title": "Test Book for Get",
            "author": "Test Author",
            "year": 2023,
            "isbn": "978-0987654321",
            "copies": 2
        }
    )
    assert create_response.status_code == 200
    book_id = create_response.json()["id"]
    
    # Then get the book
    response = client.get(f"/books/{book_id}")
    assert response.status_code == 200
    data = response.json()
    assert data["title"] == "Test Book for Get"
    assert data["author"] == "Test Author"