import uuid
from fastapi import FastAPI
from pydantic import BaseModel


class GuestAccountRequest(BaseModel):
    username: str
    password: str
    serviceLevel: str
    accountValidityPeriod: int


app = FastAPI()


@app.post("/api/v1.0/am/accounts/guest/accounts", status_code=201)
async def chat_completions(request: GuestAccountRequest):
    return {
        "status": 201,
        "message": "samp.server.success.add",
        "data": {
            "id": uuid.uuid4(),
            "username": request.username,
            "password": request.password,
            "serviceLevel": request.serviceLevel,
            "expireTime": request.accountValidityPeriod,
        }
    }
