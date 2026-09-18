# PulsePoll - Live Polling Tool

A simple full-stack live polling application built for the HCL GUVI Developer Internship task.

## Required stack

- Frontend: React
- Backend: Go + Gin
- Database: MongoDB
- Realtime: Redis
- WebSocket: Gorilla WebSocket

## Flow

Create poll -> Share link -> Audience votes -> Live results

## Features

- Signup/login with JWT authentication
- Password hashing with bcrypt
- Protected poll creation
- Backend validation
- MongoDB stores users, polls and votes
- Redis stores live vote counters
- Redis Pub/Sub sends vote events
- WebSocket pushes live results to connected browsers
- One vote per browser voter ID per poll
- Copy share link
- Responsive UI

## Project structure

```text
live-polling-app/
├── backend/
│   ├── config/
│   ├── controllers/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── utils/
│   ├── main.go
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   ├── services/
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── styles.css
│   ├── package.json
│   ├── vite.config.js
│   └── .env.example
└── README.md
```

## Local setup

### 1. MongoDB

Use MongoDB locally or MongoDB Atlas.

The application uses database name `livepoll`.

### 2. Redis

Use local Redis or Redis Cloud.

Default local URL:

```text
redis://localhost:6379
```

### 3. Backend

Open a terminal:

```bash
cd backend
```

Copy `.env.example` to `.env` and set your values.

Then:

```bash
go mod tidy
go run .
```

Backend runs on:

```text
http://localhost:8080
```

Health check:

```text
http://localhost:8080/health
```

### 4. Frontend

Open another terminal:

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on the Vite URL shown in the terminal, normally:

```text
http://localhost:5173
```

## Realtime design

When a user votes:

1. Go validates the poll and option.
2. MongoDB saves the vote.
3. Redis `HINCRBY` updates the vote counter.
4. Go reads the current Redis counters.
5. Go publishes the updated results to a Redis Pub/Sub channel.
6. Connected WebSocket clients receive the event.
7. React updates the results without refreshing.

Redis key:

```text
poll:{pollId}:counts
```

Redis Pub/Sub channel:

```text
poll:{pollId}:live
```

This makes Redis an actual part of the realtime system rather than an unused dependency.

## API endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/api/auth/signup` | Create account |
| POST | `/api/auth/login` | Login |
| POST | `/api/polls` | Create poll, requires JWT |
| GET | `/api/polls/my` | Get creator's polls, requires JWT |
| GET | `/api/polls/:id` | Get public poll |
| GET | `/api/polls/:id/results` | Get current results |
| POST | `/api/polls/:id/vote` | Cast vote |
| GET | `/api/polls/:id/live` | WebSocket realtime channel |

## MongoDB collections

### users

Stores account information. Passwords are stored as bcrypt hashes.

### polls

Stores question, options, creator, creation time and closed/open status.

### votes

Stores poll ID, option ID, voter ID and timestamp.

A unique compound index on `(pollId, voterId)` prevents one browser voter from voting twice in the same poll.

## Security decisions

- Passwords are hashed with bcrypt.
- Poll creation requires a JWT.
- Backend validates question length and option count/content.
- Backend checks that a submitted option belongs to the poll.
- MongoDB unique indexes prevent duplicate email and duplicate votes.
- Secrets are kept in environment variables.
- Client input is never trusted without server validation.

## Deployment plan

Recommended simple deployment split:

- Frontend: a React-friendly static hosting provider
- Backend: a Go-friendly web service host
- MongoDB: MongoDB Atlas
- Redis: Redis Cloud or another managed Redis provider

Set the same environment variables in the deployed services.

For the frontend:

```text
VITE_API_URL=https://YOUR-BACKEND-DOMAIN/api
```

For the backend:

```text
MONGO_URI=...
MONGO_DB=livepoll
REDIS_URL=...
JWT_SECRET=...
FRONTEND_URL=https://YOUR-FRONTEND-DOMAIN
PORT=8080
```

Before submission, test the public link from an incognito window and from a second device/browser.

## Demo checklist

1. Open the public application.
2. Sign up/login.
3. Create a poll.
4. Copy the poll link.
5. Open the link in another browser/incognito window.
6. Vote from the second browser.
7. Keep the first browser on the results page.
8. Show the result changing without refreshing.
9. Show the GitHub repository.
10. Record the required 3-5 minute video.

## Important interview topics to understand

Be able to explain:

- Why React is used for the UI.
- Why Go/Gin is the API layer.
- Why MongoDB stores polls and votes.
- Why Redis is used for counters and Pub/Sub.
- What WebSocket does.
- How JWT authentication works.
- Why bcrypt is used.
- How backend validation protects the database.
- How duplicate voting is prevented.
- What happens from clicking Vote to seeing a live result.
