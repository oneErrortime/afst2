from pydantic import BaseModel, field_validator, EmailStr
from typing import Optional
import re

class UserBase(BaseModel):
    email: EmailStr

    @field_validator('email')
    def validate_email(cls, v):
        if len(v) > 254:
            raise ValueError('Email must be less than 254 characters')
        return v

class UserCreate(UserBase):
    password: str

    @field_validator('password')
    def validate_password(cls, v):
        if len(v) > 72:
            raise ValueError('Password must not exceed 72 characters')
        if len(v) < 8:
            raise ValueError('Password must be at least 8 characters long')
        # Check if password contains both letters and numbers
        if not re.search(r'[A-Za-z]', v) or not re.search(r'\d', v):
            raise ValueError('Password must contain both letters and numbers')
        return v

class UserLogin(BaseModel):
    email: EmailStr
    password: str

    @field_validator('password')
    def validate_password(cls, v):
        if len(v) > 72:
            raise ValueError('Password must not exceed 72 characters')
        return v

    @field_validator('email')
    def validate_email(cls, v):
        if len(v) > 254:
            raise ValueError('Email must be less than 254 characters')
        return v

class User(UserBase):
    id: int
    is_active: bool

    class Config:
        from_attributes = True