from pydantic import BaseModel, field_validator, EmailStr
from typing import Optional
import re

class ReaderBase(BaseModel):
    name: str
    email: EmailStr

    @field_validator('name')
    def validate_name(cls, v):
        if not v or v.strip() == "":
            raise ValueError('Name cannot be empty')
        if len(v) > 200:
            raise ValueError('Name must be less than 200 characters')
        # Check if name contains only allowed characters (letters, spaces, hyphens, apostrophes)
        if not re.match(r'^[A-Za-z\s\-\'\.]+$', v.strip()):
            raise ValueError('Name contains invalid characters')
        return v.strip()

    @field_validator('email')
    def validate_email(cls, v):
        # Email validation is handled by EmailStr, but we can add additional checks
        if len(v) > 254:
            raise ValueError('Email must be less than 254 characters')
        return v

class ReaderCreate(ReaderBase):
    pass

class ReaderUpdate(BaseModel):
    name: Optional[str] = None
    email: Optional[EmailStr] = None

    @field_validator('name')
    def validate_name(cls, v):
        if v is not None:
            if not v or v.strip() == "":
                raise ValueError('Name cannot be empty')
            if len(v) > 200:
                raise ValueError('Name must be less than 200 characters')
            # Check if name contains only allowed characters (letters, spaces, hyphens, apostrophes)
            if not re.match(r'^[A-Za-z\s\-\'\.]+$', v.strip()):
                raise ValueError('Name contains invalid characters')
            return v.strip()
        return v

    @field_validator('email')
    def validate_email(cls, v):
        if v is not None:
            if len(v) > 254:
                raise ValueError('Email must be less than 254 characters')
        return v

class Reader(ReaderBase):
    id: int

    class Config:
        from_attributes = True