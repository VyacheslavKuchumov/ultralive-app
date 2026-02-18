FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    docker.io \
    docker-compose \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt ./requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt

COPY bot.py ./bot.py
COPY deploy.sh ./deploy.sh
RUN chmod +x /app/deploy.sh

ENV PYTHONUNBUFFERED=1
CMD ["python3", "/app/bot.py"]
