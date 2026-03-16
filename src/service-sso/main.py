from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pika
import json
import uuid
import os

app = FastAPI(title="Mock SSO Service")

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://rmq_admin:secretpassword@localhost:5672/")

class RegisterRequest(BaseModel):
    username: str
    display_name: str
    password: str

@app.post("/api/v1/auth/register")
async def register_user(req: RegisterRequest):
    user_id = str(uuid.uuid4())

    event_data = {
        "user_id": user_id,
        "username": req.username,
        "display_name": req.display_name
    }

    try:
        parameters = pika.URLParameters(RABBITMQ_URL)
        connection = pika.BlockingConnection(parameters)
        channel = connection.channel()

        channel.exchange_declare(exchange='sso.events', exchange_type='topic', durable=True)

        channel.basic_publish(
            exchange='sso.events',
            routing_key='user.created',
            body=json.dumps(event_data),
            properties=pika.BasicProperties(
                delivery_mode=2,
                content_type='application/json'
            )
        )
        connection.close()
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to publish to RabbitMQ: {str(e)}")

    return {
        "message": "User registered successfully",
        "user_id": user_id,
        "username": req.username
    }