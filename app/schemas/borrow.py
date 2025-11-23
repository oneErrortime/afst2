from pydantic import BaseModel, field_validator
from typing import Optional
from datetime import datetime

class BorrowBase(BaseModel):
    book_id: int
    reader_id: int

    @field_validator('book_id')
    def validate_book_id(cls, v):
        if v <= 0:
            raise ValueError('Book ID must be a positive integer')
        return v

    @field_validator('reader_id')
    def validate_reader_id(cls, v):
        if v <= 0:
            raise ValueError('Reader ID must be a positive integer')
        return v

class BorrowCreate(BorrowBase):
    pass

class BorrowReturn(BaseModel):
    book_id: int
    reader_id: int

    @field_validator('book_id')
    def validate_book_id(cls, v):
        if v <= 0:
            raise ValueError('Book ID must be a positive integer')
        return v

    @field_validator('reader_id')
    def validate_reader_id(cls, v):
        if v <= 0:
            raise ValueError('Reader ID must be a positive integer')
        return v

class Borrow(BorrowBase):
    id: int
    borrow_date: datetime
    return_date: Optional[datetime] = None

    class Config:
        from_attributes = True