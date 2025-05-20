from fastapi import HTTPException, Request, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import JWTError, jwt
from dotenv import load_dotenv
import os

load_dotenv()

# Настройки JWT
SECRET_KEY = os.getenv("SECRET")

# Класс для проверки токена
security = HTTPBearer()

async def jwt_middleware(request: Request, call_next):
    try:
        # Проверяем наличие заголовка Authorization
        auth_header = request.headers.get("Authorization")
        if not auth_header:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Authorization header missing",
                headers={"WWW-Authenticate": "Bearer"},
            )

        # Извлекаем токен из заголовка
        credentials = await security(request)
        token = credentials.credentials

        # Проверяем валидность токена
        try:
            payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
            # Здесь можно добавить дополнительные проверки payload, если нужно
        except JWTError:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid or expired token",
                headers={"WWW-Authenticate": "Bearer"},
            )

        # Если токен валиден, продолжаем выполнение запроса
        response = await call_next(request)
        return response

    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"An error occurred: {str(e)}",
        )