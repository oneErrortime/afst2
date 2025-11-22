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

def test_borrowing_business_logic():
    """Test the business logic for borrowing books"""
    # Create a user first (we'll skip authentication for this test by using a direct DB approach)
    # Actually, let me implement proper authentication for the tests
    
    # First, register a user
    user_response = client.post(
        "/auth/register",
        json={
            "email": "librarian@test.com",
            "password": "password123"
        }
    )
    assert user_response.status_code == 200
    
    # Then login to get token
    login_response = client.post(
        "/auth/login",
        json={
            "email": "librarian@test.com",
            "password": "password123"
        }
    )
    assert login_response.status_code == 200
    token_data = login_response.json()
    access_token = token_data["access_token"]
    headers = {"Authorization": f"Bearer {access_token}"}
    
    # Create a book
    book_response = client.post(
        "/books/",
        json={
            "title": "Test Book for Borrowing",
            "author": "Test Author",
            "year": 2023,
            "isbn": "978-1111111111",
            "copies": 2
        },
        headers=headers
    )
    assert book_response.status_code == 200
    book_data = book_response.json()
    book_id = book_data["id"]
    
    # Create a reader
    reader_response = client.post(
        "/readers/",
        json={
            "name": "Test Reader",
            "email": "reader@test.com"
        },
        headers=headers
    )
    assert reader_response.status_code == 200
    reader_data = reader_response.json()
    reader_id = reader_data["id"]
    
    # Test 1: Borrow a book successfully
    borrow_response = client.post(
        "/borrows/borrow",
        json={
            "book_id": book_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert borrow_response.status_code == 200
    
    # Check that book copies decreased
    book_response = client.get(f"/books/{book_id}")
    assert book_response.status_code == 200
    updated_book = book_response.json()
    assert updated_book["copies"] == 1  # Should be 1 less than original
    
    # Test 2: Try to borrow when no copies are available
    # First, borrow all remaining copies
    book_update_response = client.put(
        f"/books/{book_id}",
        json={"copies": 0},
        headers=headers
    )
    assert book_update_response.status_code == 200
    
    # Try to borrow again - should fail
    borrow_again_response = client.post(
        "/borrows/borrow",
        json={
            "book_id": book_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert borrow_again_response.status_code == 400
    assert "No available copies" in borrow_again_response.json()["detail"]
    
    # Test 3: Reader limit (max 3 books)
    # Update book copies back to 10
    book_update_response = client.put(
        f"/books/{book_id}",
        json={"copies": 10},
        headers=headers
    )
    assert book_update_response.status_code == 200
    
    # Create 2 more books
    book2_response = client.post(
        "/books/",
        json={
            "title": "Test Book 2 for Borrowing",
            "author": "Test Author 2",
            "year": 2023,
            "isbn": "978-2222222222",
            "copies": 5
        },
        headers=headers
    )
    assert book2_response.status_code == 200
    book2_id = book2_response.json()["id"]
    
    book3_response = client.post(
        "/books/",
        json={
            "title": "Test Book 3 for Borrowing",
            "author": "Test Author 3",
            "year": 2023,
            "isbn": "978-3333333333",
            "copies": 5
        },
        headers=headers
    )
    assert book3_response.status_code == 200
    book3_id = book3_response.json()["id"]
    
    # Borrow 2 more books for the reader (now has 3 total)
    borrow2_response = client.post(
        "/borrows/borrow",
        json={
            "book_id": book2_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert borrow2_response.status_code == 200
    
    borrow3_response = client.post(
        "/borrows/borrow",
        json={
            "book_id": book3_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert borrow3_response.status_code == 200
    
    # Try to borrow a 4th book - should fail
    book4_response = client.post(
        "/books/",
        json={
            "title": "Test Book 4 for Borrowing",
            "author": "Test Author 4",
            "year": 2023,
            "isbn": "978-4444444444",
            "copies": 5
        },
        headers=headers
    )
    assert book4_response.status_code == 200
    book4_id = book4_response.json()["id"]
    
    borrow4_response = client.post(
        "/borrows/borrow",
        json={
            "book_id": book4_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert borrow4_response.status_code == 400
    assert "cannot borrow more than 3 books" in borrow4_response.json()["detail"]
    
    # Test 4: Return functionality
    return_response = client.post(
        "/borrows/return",
        json={
            "book_id": book_id,
            "reader_id": reader_id
        },
        headers=headers
    )
    assert return_response.status_code == 200
    
    # Check that book copies increased again
    book_check_response = client.get(f"/books/{book_id}")
    assert book_check_response.status_code == 200
    final_book = book_check_response.json()
    assert final_book["copies"] == 10  # Should be back to 10 after return