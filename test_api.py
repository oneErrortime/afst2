#!/usr/bin/env python3
"""
Test script to verify API endpoints are working
"""
import requests
import json

BASE_URL = "http://localhost:8000"

def test_endpoints():
    print("Testing API endpoints...")
    
    # Test root endpoint
    print("\n1. Testing root endpoint:")
    try:
        response = requests.get(f"{BASE_URL}/")
        print(f"   Status: {response.status_code}")
        print(f"   Response: {response.json()}")
    except Exception as e:
        print(f"   Error: {e}")
    
    # Test auth endpoints
    print("\n2. Testing auth endpoints:")
    
    # Test GET /auth/login (should return 405 as it's not a valid method)
    try:
        response = requests.get(f"{BASE_URL}/auth/login")
        print(f"   GET /auth/login - Status: {response.status_code}")
    except Exception as e:
        print(f"   Error with GET /auth/login: {e}")
    
    # Test POST /auth/login with invalid data (should return 422 or 401)
    try:
        response = requests.post(f"{BASE_URL}/auth/login", json={"email": "test@test.com", "password": "test"})
        print(f"   POST /auth/login - Status: {response.status_code}")
    except Exception as e:
        print(f"   Error with POST /auth/login: {e}")
    
    # Test POST /auth/register with invalid data
    try:
        response = requests.post(f"{BASE_URL}/auth/register", json={"email": "test@test.com", "password": "test"})
        print(f"   POST /auth/register - Status: {response.status_code}")
    except Exception as e:
        print(f"   Error with POST /auth/register: {e}")
    
    # Test other endpoints
    print("\n3. Testing other endpoints:")
    endpoints = ["/books", "/readers", "/borrows"]
    for endpoint in endpoints:
        try:
            response = requests.get(f"{BASE_URL}{endpoint}")
            print(f"   GET {endpoint} - Status: {response.status_code}")
        except Exception as e:
            print(f"   Error with GET {endpoint}: {e}")

if __name__ == "__main__":
    test_endpoints()