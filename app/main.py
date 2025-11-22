from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse
import os
from .api import auth, books, readers, borrows

app = FastAPI(title="Library Management API", version="1.0.0")

# Include API routers
app.include_router(auth.router, prefix="/auth", tags=["authentication"])
app.include_router(books.router, prefix="/books", tags=["books"])
app.include_router(readers.router, prefix="/readers", tags=["readers"])
app.include_router(borrows.router, prefix="/borrows", tags=["borrows"])

# Mount the templates directory to serve static files
app.mount("/static", StaticFiles(directory="templates"), name="static")

@app.get("/")
def read_root():
    return {"message": "Welcome to Library Management API"}

@app.get("/dashboard", response_class=HTMLResponse)
def get_dashboard():
    with open("templates/index.html", "r", encoding="utf-8") as file:
        content = file.read()
    return HTMLResponse(content=content)