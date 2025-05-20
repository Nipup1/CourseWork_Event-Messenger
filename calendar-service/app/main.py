from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.routes import events
from app.scheduler import scheduler
from app.middleware import jwt_middleware
import logging

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Инициализация FastAPI
app = FastAPI(title="Meetings Microservice")

app.middleware("http")(jwt_middleware)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Разрешает запросы с любых доменов
    allow_credentials=True,
    allow_methods=["*"],  # Разрешает все HTTP-методы (GET, POST, etc.)
    allow_headers=["*"],  # Разрешает все заголовки
)

# Подключение маршрутов
app.include_router(events.router)

# Запуск шедулера при старте
@app.on_event("startup")
async def startup_event():
    scheduler.start()
    logger.info("Scheduler started")

# Остановка шедулера при завершении
@app.on_event("shutdown")
async def shutdown_event():
    scheduler.shutdown()
    logger.info("Scheduler stopped")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)