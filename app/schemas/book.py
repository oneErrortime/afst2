from pydantic import BaseModel, field_validator
from typing import Optional
import re

class BookBase(BaseModel):
    title: str
    author: str
    year: Optional[int] = None
    isbn: Optional[str] = None
    copies: Optional[int] = 1
    description: Optional[str] = None

    @field_validator('title')
    def validate_title(cls, v):
        if not v or v.strip() == "":
            raise ValueError('Title cannot be empty')
        if len(v) > 500:
            raise ValueError('Title must be less than 500 characters')
        return v.strip()

    @field_validator('author')
    def validate_author(cls, v):
        if not v or v.strip() == "":
            raise ValueError('Author cannot be empty')
        if len(v) > 200:
            raise ValueError('Author name must be less than 200 characters')
        return v.strip()

    @field_validator('year')
    def validate_year(cls, v):
        if v is not None:
            from datetime import datetime
            current_year = datetime.now().year
            if v < 1000 or v > current_year + 10:  # Allow up to 10 years in the future
                raise ValueError(f'Year must be between 1000 and {current_year + 10}')
        return v

    @field_validator('isbn')
    def validate_isbn(cls, v):
        if v is not None:
            # Remove hyphens and spaces
            clean_isbn = re.sub(r'[-\s]', '', v)
            if len(clean_isbn) not in [10, 13]:
                raise ValueError('ISBN must be either 10 or 13 characters long')
            # Check if all characters are digits (for basic validation)
            if not clean_isbn.isdigit() and not (len(clean_isbn) == 10 and clean_isbn[:-1].isdigit() and clean_isbn[-1] in '0123456789Xx'):
                raise ValueError('Invalid ISBN format')
        return v

    @field_validator('copies')
    def validate_copies(cls, v):
        if v is not None and v < 0:
            raise ValueError('Copies cannot be negative')
        return v

    @field_validator('description')
    def validate_description(cls, v):
        if v is not None and len(v) > 2000:
            raise ValueError('Description must be less than 2000 characters')
        return v

class BookCreate(BookBase):
    pass

class BookUpdate(BaseModel):
    title: Optional[str] = None
    author: Optional[str] = None
    year: Optional[int] = None
    isbn: Optional[str] = None
    copies: Optional[int] = None
    description: Optional[str] = None

class Book(BookBase):
    id: int

    class Config:
        from_attributes = True