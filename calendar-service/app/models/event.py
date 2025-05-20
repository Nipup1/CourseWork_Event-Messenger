from sqlalchemy import Column, Integer, String, DateTime, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship

Base = declarative_base()

class Event(Base):
    __tablename__ = "events"
    id = Column(Integer, primary_key=True, index=True)
    chat_id = Column(Integer, nullable=True)
    title = Column(String, index=True)
    description = Column(String, nullable=True)
    event_datetime = Column(DateTime)
    reminder_datetime = Column(DateTime)
    # Связь с таблицей event_participants через модель EventParticipant
    _participants = relationship(
        "EventParticipant",
        back_populates="event",
        cascade="all, delete-orphan"
    )

    @property
    def participants(self):
        return [participant.user_id for participant in self._participants]

    def set_participants(self, user_ids: list[int]):
        """Метод для установки списка user_id"""
        self._participants = [
            EventParticipant(user_id=user_id) for user_id in user_ids
        ]

class EventParticipant(Base):
    __tablename__ = "event_participants"
    event_id = Column(Integer, ForeignKey("events.id"), primary_key=True)
    user_id = Column(Integer, primary_key=True)
    event = relationship("Event", back_populates="_participants")