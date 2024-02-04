FROM golang

WORKDIR /app

COPY . .

RUN chmod +x ./entrypoint.sh