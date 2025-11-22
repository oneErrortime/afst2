#!/usr/bin/env python3
"""
Script to initialize the database with all tables and create an admin user
"""
import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from app.database import engine, Base
import app.models.user
import app.models.book
import app.models.reader
import app.models.borrow
from create_admin_user import create_admin_user

def init_db():
    """Initialize the database by creating all tables and an admin user"""
    print("Creating database tables...")
    
    # Create all tables
    Base.metadata.create_all(bind=engine)
    print("Database tables created successfully!")
    
    # Create admin user
    create_admin_user()
    
    print("Database initialization completed!")

if __name__ == "__main__":
    init_db()