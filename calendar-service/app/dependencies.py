from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from dotenv import load_dotenv
import os
from app.models.event import Base

load_dotenv()
DATABASE_URL = os.getenv("DATABASE_CALENDAR_URL")

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
# Base = declarative_base()
Base.metadata.create_all(bind=engine)

# Зависимость для получения сессии БД
def get_db():
    # load_dotenv()
    # engine = create_engine(os.getenv("DATABASE_URL"))
    # SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
    # Base.metadata.create_all(bind=engine)

    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()