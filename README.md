# Chatter
 
A self-hosted real-time chat web application with end-to-end encryption and file sharing.
 
---
 
## Tech Stack
 
Go | Svelte | PostgreSQL | MinIO
 
---
 
## Features
 
- Real-time WebSocket chat with notifications and read receipts
- AES-GCM encryption for messages and files at rest
- File uploads via MinIO (S3-compatible)
- User authentication
- Add users to start conversations

---
 
## Setup

1. Copy the example env files and fill in your values:
```bash
   cp .env.example .env
   cp web/.env.example web/.env
```

2. Start the stack:
```bash
   docker compose up --build
```

- API: http://localhost:8080
- Frontend: http://localhost:8081
- MinIO console: http://localhost:9003
 
---