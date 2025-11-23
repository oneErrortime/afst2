from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.orm import Session
from datetime import timedelta
from typing import Annotated

from .. import models, schemas
from ..database import get_db
from ..auth.jwt_handler import (
    get_password_hash, 
    authenticate_user, 
    create_access_token,
    get_current_active_user
)
from ..config import ACCESS_TOKEN_EXPIRE_MINUTES

router = APIRouter()


@router.post("/register", response_model=schemas.User)
def register_user(
    user: schemas.UserCreate, 
    db: Session = Depends(get_db)
):
    """
    Регистрация нового пользователя
    
    - **email**: валидный email адрес
    - **password**: минимум 8 символов, должен содержать буквы и цифры
    """
    # Проверка существующего пользователя
    existing_user = db.query(models.User).filter(
        models.User.email == user.email
    ).first()
    
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Email already registered"
        )
    
    # Создание нового пользователя
    hashed_password = get_password_hash(user.password)
    db_user = models.User(
        email=user.email, 
        hashed_password=hashed_password
    )
    
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    
    return db_user


@router.post("/login", response_model=schemas.Token)
def login_user(
    form_data: Annotated[OAuth2PasswordRequestForm, Depends()],
    db: Session = Depends(get_db)
):
    """
    Вход пользователя (OAuth2 compatible)
    
    Принимает form data:
    - **username**: email пользователя
    - **password**: пароль пользователя
    
    Возвращает JWT токен
    """
    # OAuth2PasswordRequestForm использует поле 'username', 
    # но мы используем его как email
    user = authenticate_user(db, form_data.username, form_data.password)
    
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Создание токена
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
        data={"sub": user.email}, 
        expires_delta=access_token_expires
    )
    
    return schemas.Token(
        access_token=access_token,
        token_type="bearer"
    )


@router.get("/me", response_model=schemas.User)
def read_users_me(
    current_user: Annotated[models.User, Depends(get_current_active_user)]
):
    """
    Получить информацию о текущем пользователе
    
    Требует валидный JWT токен
    """
    return current_user


@router.post("/logout")
def logout_user(
    current_user: Annotated[models.User, Depends(get_current_active_user)]
):
    """
    Выход пользователя
    
    На самом деле просто подтверждение, что токен валидный.
    Реальный logout происходит на клиенте (удаление токена).
    """
    return {"message": "Successfully logged out"}