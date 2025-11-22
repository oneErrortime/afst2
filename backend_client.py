#!/usr/bin/env python3
"""
Backend client to communicate with the Library Management API
"""
import requests
import json

# Get the server IP from the environment or use default
SERVER_IP = "0.0.0.0"  # This will be replaced with the actual IP when you run it
SERVER_PORT = 8000
BASE_URL = f"http://{SERVER_IP}:{SERVER_PORT}"

def test_connection():
    """Test basic connection to the server"""
    try:
        response = requests.get(f"{BASE_URL}/")
        print(f"✓ Connected to server: {response.json()}")
        return True
    except Exception as e:
        print(f"✗ Connection failed: {e}")
        return False

def register_user(email, password):
    """Register a new user"""
    try:
        response = requests.post(
            f"{BASE_URL}/auth/register",
            json={"email": email, "password": password}
        )
        if response.status_code == 200:
            print(f"✓ User registered successfully: {response.json()}")
            return response.json()
        elif response.status_code == 400:
            print(f"✗ User registration failed: {response.json().get('detail', 'User already exists')}")
        else:
            print(f"✗ User registration failed with status {response.status_code}: {response.text}")
        return None
    except Exception as e:
        print(f"✗ Error during registration: {e}")
        return None

def login_user(email, password):
    """Login a user and get access token"""
    try:
        response = requests.post(
            f"{BASE_URL}/auth/login",
            json={"email": email, "password": password}
        )
        if response.status_code == 200:
            data = response.json()
            print(f"✓ Login successful: {data}")
            return data.get("access_token")
        else:
            print(f"✗ Login failed with status {response.status_code}: {response.text}")
            return None
    except Exception as e:
        print(f"✗ Error during login: {e}")
        return None

def get_books(token=None):
    """Get all books from the library"""
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    
    try:
        response = requests.get(f"{BASE_URL}/books", headers=headers)
        if response.status_code == 200:
            books = response.json()
            print(f"✓ Retrieved {len(books)} books")
            return books
        else:
            print(f"✗ Failed to get books with status {response.status_code}: {response.text}")
            return []
    except Exception as e:
        print(f"✗ Error getting books: {e}")
        return []

def main():
    print("=== Library Management API Client ===")
    print(f"Connecting to: {BASE_URL}")
    
    # Test connection
    if not test_connection():
        print("Cannot connect to server. Please make sure the backend is running.")
        return
    
    print("\n=== Available Operations ===")
    print("1. Register a new user")
    print("2. Login with existing user")
    print("3. Get all books")
    print("4. Create temporary admin user (if needed)")
    print("5. Test all endpoints")
    
    choice = input("\nEnter your choice (1-5): ").strip()
    
    if choice == "1":
        email = input("Enter email: ").strip()
        password = input("Enter password: ").strip()
        register_user(email, password)
        
    elif choice == "2":
        email = input("Enter email: ").strip()
        password = input("Enter password: ").strip()
        login_user(email, password)
        
    elif choice == "3":
        token = input("Enter access token (or press Enter to skip): ").strip()
        if not token:
            print("No token provided, trying without authentication...")
        get_books(token if token else None)
        
    elif choice == "4":
        # Create admin user if it doesn't exist
        from create_admin_user import create_admin_user
        create_admin_user()
        print("Admin user creation process completed.")
        
    elif choice == "5":
        print("\n--- Testing all endpoints ---")
        
        # Test root endpoint
        print("\n1. Testing root endpoint...")
        try:
            response = requests.get(f"{BASE_URL}/")
            print(f"   Status: {response.status_code}, Response: {response.json()}")
        except Exception as e:
            print(f"   Error: {e}")
        
        # Test auth endpoints
        print("\n2. Testing auth endpoints...")
        
        # Try to login with admin user
        print("   Attempting to login with admin user...")
        admin_token = login_user("admin@example.com", "admin123")
        
        # Test getting books
        print("\n3. Testing books endpoint...")
        get_books(admin_token)
        
    else:
        print("Invalid choice!")

if __name__ == "__main__":
    main()