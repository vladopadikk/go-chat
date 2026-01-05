# Go Chat Backend

Real-time backend for private and group chats built with Go.  
The project demonstrates clean architecture, WebSocket real-time messaging, and RESTful API design.

---

##  Features

* JWT authentication (access / refresh tokens)
* Private and group chats
* Real-time messaging via WebSocket
* Message history with pagination
* Clean architecture (handler / service / repository)
* Protection against unauthorized access to chats

---

##  Architecture Overview

The project follows a layered architecture:

```
HTTP / WS
    ↓
Handlers
    ↓
Services
    ↓
Repositories
    ↓
Database
```

**Key principles:**
- WebSocket is implemented as a thin transport layer and reuses the same services as the HTTP API
- All business logic lives in the service layer
- Repositories handle database operations only
- Each layer communicates through interfaces

---

##  Transactions & Data Consistency

Critical business operations are executed within database transactions to guarantee data consistency.

Examples:
- Sending a message:
  - verifies chat membership
  - saves the message
  - commits the transaction only if all checks succeed
- Chat creation:
  - creates a chat record
  - adds chat members atomically

If any step fails, the transaction is rolled back to prevent partial data writes.

---

##  Tech Stack

* **Go 1.21+** - main language
* **Gin** - HTTP web framework
* **Gorilla WebSocket** - real-time communication
* **PostgreSQL** - database
* **JWT** - token-based authentication
* **bcrypt** - password hashing
* **pgx** - PostgreSQL driver
* **goose** - database migrations

---

##  Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── auth/                 # Authentication & JWT
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── middleware.go
│   │   └── tokens.go
│   ├── chat/                 # Chats management
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── messages/             # Messages logic (HTTP)
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── user/                 # User management
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go
│   ├── ws/                   # WebSocket hub & clients
│   │   ├── hub.go
│   │   ├── client.go
│   │   ├── handler.go
│   │   └── message.go
│   ├── database/             # DB connection & helpers
│   │   ├── connection.go
│   │   └── executor.go
│   └── config/               # Configuration
│       └── config.go
├── migrations/               # SQL migrations
│   ├── 001_create_users.sql
│   ├── 002_create_chats.sql
│   └── 003_create_messages.sql
├── .env                      # Environment variables
├── go.mod
├── go.sum
└── README.md
```

---

##  Running the Project

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- goose (for migrations)

### 1. Clone the repository

```bash
git clone https://github.com/vladopadikk/go-chat.git
cd go-chat
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment

Create `.env` file in the root:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_chat

JWT_SECRET=your-secret-key-change-in-production
```

### 4. Create database and apply migrations

```bash
createdb go_chat

goose -dir migrations postgres "user=postgres password=password dbname=go_chat sslmode=disable" up
```

### 5. Run the server

```bash
go run cmd/server/main.go
```

The server starts on `http://localhost:8080`

---

##  REST API Endpoints

### Public Endpoints

#### Register User
```http
POST /api/register
Content-Type: application/json

{
  "username": "ivan",
  "email": "ivan@example.com",
  "password": "password123"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "username": "ivan",
  "email": "ivan@example.com"
}
```

**Errors:**
- `400 Bad Request` - invalid JSON
- `409 Conflict` - email already exists

---

#### Login
```http
POST /api/login
Content-Type: application/json

{
  "email": "ivan@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors:**
- `400 Bad Request` - invalid JSON
- `404 Not Found` - user not found
- `401 Unauthorized` - invalid password

---

### Protected Endpoints

All endpoints below require the `Authorization` header:

```
Authorization: Bearer <access_token>
```

---

#### Create Private Chat
```http
POST /api/chats/private
Content-Type: application/json
Authorization: Bearer <token>

{
  "user_id": 2
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "type": "private",
  "created_at": "2026-01-05T10:30:00Z"
}
```

**Description:**  
Creates a private chat between the authenticated user and the specified user.  
If chat already exists, returns the existing one.

**Errors:**
- `400 Bad Request` - invalid JSON
- `401 Unauthorized` - missing or invalid token
- `500 Internal Server Error` - database error

---

#### Create Group Chat
```http
POST /api/chats/group
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Project Team",
  "participants": [2, 3, 4]
}
```

**Response:** `200 OK`
```json
{
  "id": 2,
  "type": "group",
  "created_at": "2026-01-05T10:35:00Z"
}
```

**Description:**  
Creates a group chat with the specified name and participants.  
The authenticated user is automatically added as a member.

**Errors:**
- `400 Bad Request` - invalid JSON
- `401 Unauthorized` - missing or invalid token
- `500 Internal Server Error` - database error

---

#### Get User's Chats
```http
GET /api/chats
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "chats": [
    {
      "id": 1,
      "type": "private",
      "created_at": "2026-01-05T10:30:00Z"
    },
    {
      "id": 2,
      "type": "group",
      "created_at": "2026-01-05T10:35:00Z"
    }
  ]
}
```

**Description:**  
Returns all chats where the authenticated user is a member.

**Errors:**
- `401 Unauthorized` - missing or invalid token
- `500 Internal Server Error` - database error

---

#### Send Message (HTTP)
```http
POST /api/messages/send
Content-Type: application/json
Authorization: Bearer <token>

{
  "chat_id": 1,
  "content": "Hello, world!"
}
```

**Response:** `200 OK`
```json
{
  "id": 10,
  "chat_id": 1,
  "sender_id": 1,
  "content": "Hello, world!",
  "created_at": "2026-01-05T10:40:00Z"
}
```

**Description:**  
Sends a message to the specified chat.  
User must be a member of the chat.

**Errors:**
- `400 Bad Request` - invalid JSON
- `401 Unauthorized` - missing or invalid token
- `403 Forbidden` - user is not a member of the chat
- `500 Internal Server Error` - database error

---

#### Get Messages
```http
GET /api/messages/get?chat_id=1&limit=50&offset=0
Authorization: Bearer <token>
```

**Query Parameters:**
- `chat_id` (required) - ID of the chat
- `limit` (optional) - number of messages to return (default: 50, max: 100)
- `offset` (optional) - pagination offset (default: 0)

**Response:** `200 OK`
```json
{
  "messages": [
    {
      "id": 10,
      "chat_id": 1,
      "sender_id": 1,
      "content": "Hello, world!",
      "created_at": "2026-01-05T10:40:00Z"
    },
    {
      "id": 9,
      "chat_id": 1,
      "sender_id": 2,
      "content": "Hi there!",
      "created_at": "2026-01-05T10:38:00Z"
    }
  ]
}
```

**Description:**  
Returns messages from the specified chat in reverse chronological order.  
User must be a member of the chat.

**Errors:**
- `400 Bad Request` - missing or invalid chat_id
- `401 Unauthorized` - missing or invalid token
- `403 Forbidden` - user is not a member of the chat
- `500 Internal Server Error` - database error

---

##  WebSocket API

### Connect to WebSocket

```
ws://localhost:8080/api/ws
```

**Headers:**
```
Authorization: Bearer <access_token>
```

**Description:**  
Establishes a WebSocket connection for real-time messaging.  
User is automatically subscribed to all their chats.

**Connection Requirements:**
- Valid JWT token
- User must have at least one chat 
---

### Send Message (Client → Server)

```json
{
  "type": "send_message",
  "payload": {
    "chat_id": 1,
    "content": "Hello from WebSocket!"
  }
}
```

**Description:**  
Sends a message to the specified chat.  
Message is saved to database and broadcasted to all chat members.

**Validations:**
- User must be a member of the chat
- Content cannot be empty
- chat_id must be valid

---

### Receive Messages (Server → Client)

#### New Message
```json
{
  "type": "new_message",
  "payload": {
    "id": 10,
    "chat_id": 1,
    "sender_id": 2,
    "content": "Hello from WebSocket!",
    "created_at": "2026-01-05T10:45:00Z"
  }
}
```

**Description:**  
Broadcasted to all members of the chat when a new message is sent.  
Includes the sender's ID to distinguish own messages from others.

---

#### Error
```json
{
  "type": "error",
  "payload": {
    "message": "you are not a member of this chat"
  }
}
```

**Possible Error Messages:**
- `"invalid message format"` - JSON parsing error
- `"unknown message type"` - unsupported message type
- `"invalid payload"` - payload doesn't match expected format
- `"you are not a member of this chat"` - access denied
- `"message content cannot be empty"` - validation error

---

##  Security

- **Password hashing:** bcrypt with cost factor 10
- **JWT tokens:**
  - Access token: 15 minutes lifetime
  - Refresh token: 7 days lifetime
- **Authorization:** All protected endpoints verify JWT token
- **Access control:** Users can only access chats they are members of
- **SQL injection protection:** Parameterized queries throughout
- **WebSocket auth:** JWT verification on connection upgrade

---

##  Database Schema

### users
```sql
id            SERIAL PRIMARY KEY
email         VARCHAR(255) UNIQUE NOT NULL
username      VARCHAR(100) NOT NULL
password_hash TEXT NOT NULL
created_at    TIMESTAMP NOT NULL DEFAULT NOW()
```

### chats
```sql
id         BIGSERIAL PRIMARY KEY
type       VARCHAR(20) NOT NULL  -- 'private' or 'group'
name       VARCHAR(255)          -- for group chats only
created_at TIMESTAMP NOT NULL DEFAULT NOW()

CONSTRAINT type_check CHECK (type IN ('private', 'group'))
```

### chat_members
```sql
chat_id   BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE
user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
joined_at TIMESTAMP NOT NULL DEFAULT NOW()

PRIMARY KEY (chat_id, user_id)
```

### messages
```sql
id         BIGSERIAL PRIMARY KEY
chat_id    BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE
sender_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
content    TEXT NOT NULL
created_at TIMESTAMP NOT NULL DEFAULT NOW()

INDEX idx_message_chat_id_created_at ON (chat_id, created_at)
```

---

##  Testing

### Manual Testing with Postman

1. **Register two users** via `POST /api/register`
2. **Login both users** via `POST /api/login` to get tokens
3. **Create a private chat** between them via `POST /api/chats/private`
4. **Send messages** via `POST /api/messages/send`
5. **Get message history** via `GET /api/messages/get`
6. **Connect to WebSocket** and test real-time messaging

### Testing with curl

```bash
# Register user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@test.com","password":"123456"}'

# Login
TOKEN=$(curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"123456"}' \
  | jq -r '.access_token')

# Get chats
curl http://localhost:8080/api/chats \
  -H "Authorization: Bearer $TOKEN"

# Send message
curl -X POST http://localhost:8080/api/messages/send \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"chat_id":1,"content":"Hello!"}'
```

---

##  Possible Improvements
- [ ] File and image attachments
- [ ] Message search
- [ ] Edit/delete messages
- [ ] Docker & Docker Compose


---

##  Purpose

This project was built as a pet project to demonstrate:
- Backend architecture skills
- Real-time communication using WebSocket
- Clean Go code practices
- RESTful API design
- Database design and migrations
- Authentication and authorization

