# **API Documentation**

## **Base URL**
- The application listens on the address specified in the `HTTP_LISTEN_ADDRESS` environment variable.

## **Authentication**
- Authentication is required for most routes except for `/auth`.
- Authentication uses JWT (JSON Web Tokens).
- Routes under `/admin` require additional admin privileges.

---

## **Routes**

### **Public Routes**
#### **Authentication**
- **`POST /api/auth`**
  - **Description**: Authenticates a user and generates a JWT.
  - **Handler**: `authHandler.HandleAuthenticate`
  - **Request Body**: JSON with credentials (e.g., email, password).
    ```json
    {
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: status 203 No Content, "X-Authorization" header with a JWT.
    - Failure: status 401 Unauthorized.
    ```json
    {
      "error": "invalid credentials"
    }
    ```

- **`POST /api/register`**
  - **Description**: Creates a new user.
  - **Handler**: `userHandler.HandlePostUser`.
  - **Request Body**:
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": false
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "firstName": "firstName length should be at least %d characters",
      "lastName": "lastName length should be at least %d characters",
      "email": "email %s is invalid",
      "password": "password length should be at least %d characters"
    }
    ```
    - Failure: 422 Unprocessable Entity. (Email in use)
    ```json
    {
      "email": "email in use",
    }
    ```

---

### **Authenticated Routes (`/api/v1`)**

#### **User Routes**
- **`GET /api/v1/users`**
  - **Description**: Fetches the authenticated user's details.
  - **Handler**: `userHandler.HandleGetMyUser`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": false
    }
    ```

- **`PATCH /api/v1/users/`**
  - **Description**: Updates details of a specific user.
  - **Handler**: `userHandler.HandlePatchMyUser`.
  - **Request Body**: (Each field is optional)
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "updated": "673d37d2a0d5e53e1ceb4df7"
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "firstName": "firstName length should be at least %d characters",
      "lastName": "lastName length should be at least %d characters",
      "email": "email %s is invalid",
      "password": "password length should be at least %d characters"
    }
    ```
    - Failure: 422 Unprocessable Entity. (Email in use)
    ```json
    {
      "email": "email in use",
    }
    ```

---

#### **Hotel Routes**
- **`GET /api/v1/hotels`**
  - **Description**: Fetches a list of all hotels.
  - **Handler**: `hotelHandler.HandleGetHotels`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "673d37d2a0d5e53e1cebade3",
        "name": "Grand Plaza Hotel",
        "location": "New York, NY",
        "rooms": [
          "67156bd2a0d5e53e1ceb3e46",
          "5ea56b6b40d5e53e1ce3e4f7",
          "156b3e42a0d5e3e41ceb43e4",
          "6156b7d2a0d5e53e1ceb3e44"
        ],
        "rating": 5
      }
    ]
    ```

- **`GET /api/v1/hotels/:id`** (:id replaced with an ID)
  - **Description**: Fetches details of a specific hotel by its ID.
  - **Handler**: `hotelHandler.HandleGetHotel`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "id": "673d37d2a0d5e53e1cebade3",
      "name": "Grand Plaza Hotel",
      "location": "New York, NY",
      "rooms": [
        "67156bd2a0d5e53e1ceb3e46",
        "5ea56b6b40d5e53e1ce3e4f7",
        "156b3e42a0d5e3e41ceb43e4",
        "6156b7d2a0d5e53e1ceb3e44"
      ],
      "rating": 5
    }
    ```
    - Failure: 404 Not Found.

---

#### **Room Routes**
- **`GET /api/v1/hotels/:hid/rooms`** (:hid replaced with an ID)
  - **Description**: Fetches all rooms in a specific hotel.
  - **Handler**: `roomHandler.HandleGetRooms`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "5ea56b6b40d5e53e1ce3e4f7",
        "size": "Large",
        "price": 150.0,
        "hotelID": "673d37d2a0d5e53e1cebade3"
      }
    ]
    ```
    - Failure: 404 Not Found.


- **`GET /api/v1/rooms/:id`** (:id replaced with an ID)
  - **Description**: Fetches details of a specific room by its ID within a hotel.
  - **Handler**: `roomHandler.HandleGetRoomByID`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "id": "5ea56b6b40d5e53e1ce3e4f7",
      "size": "Large",
      "price": 150.0,
      "hotelID": "673d37d2a0d5e53e1cebade3"
    }
    ```
    - Failure: 404 Not Found.

---

#### **Booking Routes**
- **`GET /api/v1/rooms/:id/bookings`** (:id replaced with an ID)
  - **Description**: Fetches all bookings for a specific room.
  - **Handler**: `bookingHandler.HandleGetBookingsByRoom`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "23abdc2aa0d5e53e1ceb32ea",
        "userID": "673d37d2a0d5e53e1ceb4df7",
        "hotelID": "673d37d2a0d5e53e1cebade3",
        "roomID": "5ea56b6b40d5e53e1ce3e4f7",
        "fromDate": "2024-11-17T00:00:00Z",
        "toDate": "2024-11-20T00:00:00Z",
        "CreatedDate": "2024-11-10T10:00:00Z",
        "cancelledAt": "2024-11-15T15:30:00Z",
        "cancelled": true
      }
    ]
    ```
    - Failure: 404 Not Found.

- **`POST /api/v1/rooms/:id/bookings`** (:id replaced with an ID)
  - **Description**: Creates a booking for a specific room.
  - **Handler**: `bookingHandler.HandlePostBooking`.
  - **Request Body**:
    ```json
    {
      "fromDate": "2024-11-17T00:00:00Z",
      "toDate": "2024-11-20T00:00:00Z"
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "67890",
      "userID": "12345",
      "hotelID": "9876543210",
      "roomID": "54321",
      "fromDate": "2024-11-17T00:00:00Z",
      "toDate": "2024-11-20T00:00:00Z",
      "CreatedDate": "2024-11-10T10:00:00Z",
      "cancelledAt": "2024-11-15T15:30:00Z",
      "cancelled": true
    }
    ```
    - Failure: 400 Bad Request. (Invalid pair of dates)
    ```json
    {
      "error": "invalid date"
    }
    ```
    - Failure: 422 Unprocessable Entity. (Dates are busy)
    ```json
    {
      "error": "unavailable date"
    }
    ```

- **`GET /api/v1/hotels/:hid/bookings`** (:hid replaced with an ID)
  - **Description**: Fetches all bookings for a specific hotel.
  - **Handler**: `bookingHandler.HandleGetBookingsByHotel`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "23abdc2aa0d5e53e1ceb32ea",
        "userID": "673d37d2a0d5e53e1ceb4df7",
        "hotelID": "673d37d2a0d5e53e1cebade3",
        "roomID": "5ea56b6b40d5e53e1ce3e4f7",
        "fromDate": "2024-11-17T00:00:00Z",
        "toDate": "2024-11-20T00:00:00Z",
        "CreatedDate": "2024-11-10T10:00:00Z",
        "cancelledAt": "2024-11-15T15:30:00Z",
        "cancelled": true
      }
    ]
    ```
    - Failure: 404 Not Found.

- **`GET /api/v1/bookings`**
  - **Description**: Fetches all bookings across the system.
  - **Handler**: `bookingHandler.HandleGetBookings`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "23abdc2aa0d5e53e1ceb32ea",
        "userID": "673d37d2a0d5e53e1ceb4df7",
        "hotelID": "673d37d2a0d5e53e1cebade3",
        "roomID": "5ea56b6b40d5e53e1ce3e4f7",
        "fromDate": "2024-11-17T00:00:00Z",
        "toDate": "2024-11-20T00:00:00Z",
        "CreatedDate": "2024-11-10T10:00:00Z",
        "cancelledAt": "2024-11-15T15:30:00Z",
        "cancelled": true
      }
    ]
    ```
    - Failure: 404 Not Found.

- **`PATCH /api/v1/bookings/:id`** (:id replaced with an ID)
  - **Description**: Cancels a booking by its ID.
  - **Handler**: `bookingHandler.HandleCancelBooking`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "cancelled": "23abdc2aa0d5e53e1ceb32ea"
    }
    ```
    - Failure: 404 Not Found.
    - Failure: 401 Unauthorized. (Trying to cancel other user booking)

---

### **Admin Routes (`/api/v1/admin`)**

#### **User Management**
- **`PATCH /api/v1/admin/users/:id`** (:id replaced with an ID)
  - **Description**: Updates details of a specific user.
  - **Handler**: `userHandler.HandlePatchUser`.
  - **Request Body**: (Each field is optional)
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "updated": "673d37d2a0d5e53e1ceb4df7"
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "firstName": "firstName length should be at least %d characters",
      "lastName": "lastName length should be at least %d characters",
      "email": "email %s is invalid",
      "password": "password length should be at least %d characters"
    }
    ```
    - Failure: 422 Unprocessable Entity. (Email in use)
    ```json
    {
      "email": "email in use",
    }
    ```

- **`DELETE /api/v1/admin/users/:id`** (:id replaced with an ID)
  - **Description**: Deletes a user by ID.
  - **Handler**: `userHandler.HandleDeleteUser`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "deleted": "673d37d2a0d5e53e1ceb4df7"
    }
    ```
    - Failure: 404 Not Found.

- **`POST /api/v1/admin/users`**
  - **Description**: Creates a new user.
  - **Handler**: `userHandler.HandlePostUser`.
  - **Request Body**:
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": false
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "firstName": "firstName length should be at least %d characters",
      "lastName": "lastName length should be at least %d characters",
      "email": "email %s is invalid",
      "password": "password length should be at least %d characters"
    }
    ```

- **`POST /api/v1/admin/users/admin`**
  - **Description**: Creates a new admin user.
  - **Handler**: `userHandler.HandlePostAdminUser`.
  - **Request Body**:
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "password": "secertpassword"
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": true
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "firstName": "firstName length should be at least %d characters",
      "lastName": "lastName length should be at least %d characters",
      "email": "email %s is invalid",
      "password": "password length should be at least %d characters"
    }
    ```
    - Failure: 422 Unprocessable Entity. (Email in use)
    ```json
    {
      "email": "email in use",
    }
    ```

- **`GET /api/v1/admin/users/me`**
  - **Description**: Fetches the authenticated admin user's details.
  - **Handler**: `userHandler.HandleGetMyUser`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": true
    }
    ```

- **`GET /api/v1/admin/users`**
  - **Description**: Fetches a list of all users.
  - **Handler**: `userHandler.HandleGetUsers`.
  - **Response**:
    - Success: 200 OK.
    ```json
    [
      {
        "id": "673d37d2a0d5e53e1ceb4df7",
        "firstName": "John",
        "lastName": "Doe",
        "email": "john.doe@example.com",
        "isAdmin": true
      }
    ]
    ```

- **`GET /api/v1/admin/users/:id`** (:id replaced with an ID)
  - **Description**: Fetches details of a specific user by ID.
  - **Handler**: `userHandler.HandleGetUser`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "id": "673d37d2a0d5e53e1ceb4df7",
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "isAdmin": true
    }
    ```
    - Failure: 404 Not Found.

---

#### **Room Management**
- **`DELETE /api/v1/admin/rooms/id`** (:id replaced with an ID)
  - **Description**: Deletes a room by ID.
  - **Handler**: `roomHandler.HandleDeleteRoom`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "deleted": "5ea56b6b40d5e53e1ce3e4f7"
    }
    ```
    - Failure: 404 Not Found.

- **`POST /api/v1/admin/hotels/:hid/rooms`** (:hid replaced with an ID)
  - **Description**: Creates a new room in a specific hotel.
  - **Handler**: `roomHandler.HandlePostRoom`.
  - **Request Body**:
    ```json
    {
      "size": "Large",
      "price": 150.0
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "5ea56b6b40d5e53e1ce3e4f7",
      "size": "Large",
      "price": 150.0,
      "hotelID": "673d37d2a0d5e53e1cebade3"
    }
    ```
    - Failure: 400 Bad Request.
    ```json
    {
        "error":"invalid parameters"
    }
    ```
    - Failure: 404 Not Found.

---

#### **Hotel Management**
- **`DELETE /api/v1/admin/hotels/:id`** (:id replaced with an ID)
  - **Description**: Deletes a hotel and all its rooms.
  - **Handler**: `hotelHandler.HandleDeleteHotel` and `roomHandler.HandleDeleteRoomsByHotel`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "deleted": "673d37d2a0d5e53e1cebade3"
    }
    ```
    - Failure: 404 Not Found.

- **`POST /api/v1/admin/hotels`**
  - **Description**: Creates a new hotel.
  - **Handler**: `hotelHandler.HandlePostHotel`.
  - **Request Body**:
    ```json
    {
      "name": "Grand Hotel",
      "location": "Paris, France",
      "rating": 5
    }
    ```
  - **Response**:
    - Success: 201 Created.
    ```json
    {
      "id": "673d37d2a0d5e53e1cebade3",
      "name": "Grand Hotel",
      "location": "Paris, France",
      "rooms": [],
      "rating": 5
    }
    ```
    - Failure: 400 Bad Request.
    ```json
    {
      "name": "hotel name should be at least %d characters",
      "location": "hotel location should be at least %d characters",
      "rating": "hotel rating should be greater than %d"
    }
    ```

- **`PATCH /api/v1/admin/hotels/:id`** (:id replaced with an ID)
  - **Description**: Updates hotel details.
  - **Handler**: `hotelHandler.HandlePatchHotel`.
  - **Request Body**: (Each field is optional)
    ```json
    {
      "name": "Grand Hotel",
      "location": "Paris, France",
    }
    ```
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "updated": "673d37d2a0d5e53e1ceb4df7"
    }
    ```
    - Failure: 400 Bad Request. (Each field is optional)
    ```json
    {
      "name": "hotel name should be at least %d characters",
      "location": "hotel location should be at least %d characters",
    }
    ```

---

#### **Booking Management**
- **`DELETE /api/v1/admin/bookings/:id`** (:id replaced with an ID)
  - **Description**: Deletes a booking by ID.
  - **Handler**: `bookingHandler.HandleDeleteBooking`.
  - **Response**:
    - Success: 200 OK.
    ```json
    {
      "deleted": "673d37d2a0d5e53e1cebade3"
    }
    ```
    - Failure: 404 Not Found.

---

## **Middleware**
- **`middleware.JWTAuthentication(uStore)`**:
  - Enforces JWT authentication for routes under `/api/v1`.
  - JWT must be present in "X-Authorization" header.
- **`middleware.AdminMiddleware`**:
  - Enforces admin privileges for routes under `/api/v1/admin`.

---

## **Environment Variables**
Ensure the following environment variables are set in `.env`:
- `MONGO_DB_URI`: MongoDB connection URI (e.g, `mongodb://localhost:27017`)
- `MONGO_DB_NAME`: MongoDB database name (e.g, `hotel-reserv-db`)
- `MONGO_DB_TEST_NAME`: MongoDB test database name (e.g, `hotel-reserv-db-test`)
- `HTTP_LISTEN_ADDRESS`: The address the server listens on (e.g., `:4000`).
### Defaults:
```env
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=hotel-reserv-db
MONGO_DB_TEST_NAME=hotel-reserv-db-test
HTTP_LISTEN_ADDRESS=:4000
```

---

This documentation provides an overview of all available routes, their methods,
access controls, and expected request/response structures.
