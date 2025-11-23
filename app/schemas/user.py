from pydantic import BaseModel, field_validator, EmailStr, Field
from typing import Optional
import re


class UserBase(BaseModel):
    """Базовая схема пользователя"""
    email: EmailStr = Field(..., max_length=254, description="User email address")

    @field_validator('email')
    @classmethod
    def validate_email(cls, v: str) -> str:
        """Валидация email"""
        if len(v) > 254:
            raise ValueError('Email must be less than 254 characters')
        
        # Дополнительная проверка на безопасность
        if '\x00' in v or any(ord(c) < 32 for c in v):
            raise ValueError('Email contains invalid characters')
        
        return v.lower().strip()


class UserCreate(UserBase):
    """Схема создания пользователя"""
    password: str = Field(
        ..., 
        min_length=8, 
        max_length=72,
        description="Password (8-72 chars, must contain letters and numbers)"
    )

    @field_validator('password')
    @classmethod
    def validate_password(cls, v: str) -> str:
        """Валидация пароля"""
        if len(v) > 72:
            raise ValueError('Password must not exceed 72 characters (bcrypt limitation)')
        
        if len(v) < 8:
            raise ValueError('Password must be at least 8 characters long')
        
        # Проверка на буквы и цифры
        if not re.search(r'[A-Za-z]', v):
            raise ValueError('Password must contain letters')
        
        if not re.search(r'\d', v):
            raise ValueError('Password must contain numbers')
        
        # Проверка на null bytes и управляющие символы
        if '\x00' in v or any(ord(c) < 32 for c in v if c not in '\t\n\r'):
            raise ValueError('Password contains invalid characters')
        
        return v


class UserLogin(BaseModel):
    """Схема для логина"""
    email: EmailStr = Field(..., max_length=254)
    password: str = Field(..., max_length=72)

    @field_validator('password')
    @classmethod
    def validate_password(cls, v: str) -> str:
        """Базовая валидация пароля для логина"""
        if len(v) > 72:
            raise ValueError('Password must not exceed 72 characters')
        
        if '\x00' in v:
            raise ValueError('Invalid password format')
        
        return v

    @field_validator('email')
    @classmethod
    def validate_email(cls, v: str) -> str:
        """Валидация email для логина"""
        if len(v) > 254:
            raise ValueError('Email must be less than 254 characters')
        
        return v.lower().strip()


class User(UserBase):
    """Схема пользователя для ответа"""
    id: int = Field(..., gt=0)
    is_active: bool = True

    model_config = {"from_attributes": True}


class Token(BaseModel):
    """Схема токена"""
    access_token: str = Field(..., min_length=1)
    token_type: str = Field(default="bearer")

    @field_validator('access_token')
    @classmethod
    def validate_token(cls, v: str) -> str:
        """Валидация токена"""
        if not v or len(v.strip()) == 0:
            raise ValueError('Token cannot be empty')
        return v
    
    @field_validator('token_type')
    @classmethod
    def validate_token_type(cls, v: str) -> str:
        """Валидация типа токена"""
        if v.lower() != 'bearer':
            raise ValueError('Only bearer token type is supported')
        return v.lower()


class TokenData(BaseModel):
    """Схема данных токена"""
    email: Optional[str] = None