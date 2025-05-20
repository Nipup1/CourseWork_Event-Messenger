from fastapi import APIRouter, HTTPException, Depends, Query
from sqlalchemy.orm import Session
from app.models.event import Event, EventParticipant
from app.schemas.event import EventCreate, EventResponse
from app.dependencies import get_db
from app.scheduler import scheduler
from apscheduler.triggers.date import DateTrigger
import os
import websockets
import json
import logging
from datetime import datetime
from typing import List
from dotenv import load_dotenv


logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/events", tags=["events"])

# URL WebSocket и токен админа
load_dotenv()
WEBSOCKET_URL = f"{os.getenv('WS_URL')}{os.getenv('ADMIN_TOKEN')}"

logger.info(WEBSOCKET_URL)
# Функция для отправки напоминания через WebSocket
async def send_reminder(event_id: int, title: str, chat_id: int, event_datetime: datetime, is_event_datetime: bool):
    try:
        async with websockets.connect(WEBSOCKET_URL) as websocket:
            if is_event_datetime:
                str = f"Напоминание: Событие '{title}' началось!"
            else:
                str = f"Напоминание: Событие '{title}' начнётся {event_datetime.strftime("%Y-%m-%d %H:%M:%S")}!"
            
            message = {
                "chat_id": chat_id,
                "content": str
            }
            await websocket.send(json.dumps(message))
            logger.info(f"Reminder sent for event {event_id}: {title} to chat {chat_id}")
    except Exception as e:
        logger.error(f"Failed to send reminder for event {event_id}: {str(e)}")

@router.post("/", response_model=EventResponse)
async def create_event(event: EventCreate, db: Session = Depends(get_db)):
    try:
        # Создаем событие
        db_event = Event(
            chat_id=event.chat_id,
            title=event.title,
            description=event.description,
            event_datetime=event.event_datetime,
            reminder_datetime=event.reminder_datetime
        )
        # Устанавливаем участников
        db_event.set_participants(event.participants)
        db.add(db_event)
        db.commit()
        db.refresh(db_event)
        
        # Планируем напоминание
        scheduler.add_job(
            send_reminder,
            trigger=DateTrigger(run_date=event.reminder_datetime),
            args=[db_event.id, db_event.title, db_event.chat_id, db_event.event_datetime, False]
        )

        scheduler.add_job(
            send_reminder,
            trigger=DateTrigger(run_date=event.event_datetime),
            args=[db_event.id, db_event.title, db_event.chat_id, db_event.event_datetime, True]
        )
        
        logger.info(f"Created event with id {db_event.id} and participants {db_event.participants}")
        return db_event
    except Exception as e:
        logger.error(f"Failed to create event: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/", response_model=List[EventResponse])
async def get_events(
    user_id: int = Query(..., description="ID of the user to filter events"),
    db: Session = Depends(get_db)
):
    # Фильтруем события по user_id
    events = db.query(Event).join(Event._participants).filter(EventParticipant.user_id == user_id).all()
    return events

@router.get("/{event_id}", response_model=EventResponse)
async def get_event(event_id: int, db: Session = Depends(get_db)):
    event = db.query(Event).filter(Event.id == event_id).first()
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")
    return event

@router.put("/{event_id}", response_model=EventResponse)
async def update_event(event_id: int, event: EventCreate, db: Session = Depends(get_db)):
    logger.info("ADWpdkapwdkpadkpa")
    logger.info(event)
    db_event = db.query(Event).filter(Event.id == event_id).first()
    if not db_event:
        raise HTTPException(status_code=404, detail="Event not found")
    
    # Обновляем поля события
    for key, value in event.dict(exclude={"participants"}).items():
        setattr(db_event, key, value)
    
    # Обновляем участников
    db_event.set_participants(event.participants)
    
    db.commit()
    db.refresh(db_event)
    
    # Обновляем напоминание
    scheduler.add_job(
        send_reminder,
        trigger=DateTrigger(run_date=event.reminder_datetime),
        args=[db_event.id, db_event.title, db_event.chat_id]
    )
    
    logger.info(f"Updated event with id {db_event.id} and participants {db_event.participants}")
    return db_event

@router.delete("/{event_id}")
async def delete_event(event_id: int, db: Session = Depends(get_db)):
    event = db.query(Event).filter(Event.id == event_id).first()
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")
    db.delete(event)
    db.commit()
    logger.info(f"Deleted event with id {event_id}")
    return {"message": "Event deleted"}