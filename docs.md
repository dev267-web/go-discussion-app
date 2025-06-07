# 🗂️ API Documentation: Go Discussion App

---

## 🧑 User & Profile Management

Users must create a profile to initiate discussions.

| Method | Endpoint         | Description                      |
|--------|------------------|----------------------------------|
| POST   | `/auth/register` | Register a new user              |
| POST   | `/auth/login`    | Authenticate user and return token |
| GET    | `/users/:id`     | Get user profile by ID           |
| PUT    | `/users/:id`     | Update user profile              |
| DELETE | `/users/:id`     | Delete user profile              |

- **All protected routes use JWT-based authentication middleware.**
- **DTOs are used to validate user input.**

---

## 💬 Discussion Topic Management

| Method | Endpoint                | Description                                   |
|--------|-------------------------|-----------------------------------------------|
| POST   | `/discussions`          | Create a new discussion (auth & profile required) |
| GET    | `/discussions`          | Get all discussions (optional filters)        |
| GET    | `/discussions/:id`      | Get a single discussion topic                 |
| PUT    | `/discussions/:id`      | Update a discussion topic                     |
| DELETE | `/discussions/:id`      | Delete a discussion topic                     |

### 🏷️ Filtering & Tagging

| Method | Endpoint                        | Description                        |
|--------|---------------------------------|------------------------------------|
| GET    | `/discussions/user/:userId`     | Get all discussions by a user      |
| GET    | `/discussions/tag/:tag`         | Get discussions by a tag           |
| POST   | `/discussions/:id/tags`         | Add tags to a discussion topic     |

### ⏰ Scheduled Discussions

| Method | Endpoint                  | Description                              |
|--------|---------------------------|------------------------------------------|
| POST   | `/discussions/schedule`   | Schedule a discussion for future posting |

---

## 💬 Comments on Discussions

| Method | Endpoint                          | Description                        |
|--------|-----------------------------------|------------------------------------|
| POST   | `/discussions/:id/comments`       | Add a comment to a discussion      |
| GET    | `/discussions/:id/comments`       | Get all comments of a discussion   |

---

## 📩 Subscriptions & Email Notifications

| Method | Endpoint                              | Description                                         |
|--------|---------------------------------------|-----------------------------------------------------|
| POST   | `/discussions/:id/subscribe`          | Subscribe to a discussion via email                 |
| DELETE | `/discussions/:id/unsubscribe`        | Unsubscribe from a discussion                       |
| POST   | `/discussions/:id/notify`             | (Internal) Trigger email notifications to subscribers|

---

## 🧪 Utility / Admin APIs (Optional)

| Method | Endpoint      | Description                                 |
|--------|--------------|---------------------------------------------|
| GET    | `/tags`      | Get all available tags                      |
| GET    | `/health`    | Health check endpoint for monitoring        |
| GET    | `/docs`      | API documentation (Swagger or similar)      |

---