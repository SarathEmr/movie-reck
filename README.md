# 🎬 Movie Recommendation App

A containerized movie recommendation application with user authentication, REST APIs, and an LLM-powered recommendation engine.

---

## 🚀 Build and Run the App

Use Docker Compose to build and start the application:

```bash
docker-compose up --build
```

Seeding the `app_user` table:
```
INSERT INTO app_user (id, username, password_hash, name) VALUES   (gen_random_uuid(), 'sarath.n', '$2a$10$VST6mGqsynym6dAgTFROde/PgnOsjymB21/O2r/YVotqpDYD5uq2e', 'Sarath N'); 
```

---

## 🌐 Access the Application

Once running, open your browser and navigate to:

```
http://localhost:3000
```

**Development Environment Credentials:**
- **Username:** `sarath.n`
- **Password:** `abc123`

---

## 🔌 API Endpoints (REST)

- **POST** `/login/`  
  Authenticate a user (passwords are securely hashed and stored in the database).

- **GET** `/movie/recommendation/`  
  Retrieve movie recommendations for the authenticated user.

---

## 🗄️ Technology Stack

### Backend
- Language: Go
- Router: Gin (https://github.com/gin-gonic/gin)
- Authentication: Password hash stored in the database

### Database
- PostgreSQL

### Frontend
- ReactJS

### LLM
- Gemini 2.5 Flash Lite

### Containerization
- Docker & Docker Compose

---

## Screenshots

![alt text](<Screenshot 2026-05-09 165538.png>)

![alt text](<Screenshot 2026-05-09 164952.png>)