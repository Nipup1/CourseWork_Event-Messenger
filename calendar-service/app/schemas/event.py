from __future__ import annotations
from pydantic import BaseModel
from datetime import datetime
from typing import List, Optional

class EventCreate(BaseModel):
    chat_id: Optional[int] = None
    title: str
    description: Optional[str] = None
    event_datetime: datetime
    reminder_datetime: datetime
    participants: List[int]  # Список ID участников

class EventResponse(EventCreate):
    id: int

    class Config:
        from_attributes = True