# Refresh Token Usage Guide

## Overview
The tickr API now supports refresh tokens for enhanced security. Access tokens expire in 15 minutes, while refresh tokens are valid for 7 days.

## Authentication Flow

### 1. Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### 2. Using Access Token
```bash
curl -X GET http://localhost:8080/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Refresh Access Token
When the access token expires (after 15 minutes), use the refresh token:

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. Logout
To invalidate refresh tokens:

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6"
  }'
```

## Security Features

1. **Short-lived Access Tokens**: 15 minutes expiration
2. **Long-lived Refresh Tokens**: 7 days expiration
3. **Token Invalidation**: Logout invalidates all user refresh tokens
4. **Database Storage**: Refresh tokens are stored securely in the database
5. **Automatic Cleanup**: Expired refresh tokens are automatically cleaned up

## Environment Variables

Set the JWT secret in your environment:
```bash
export JWT_SECRET="your-super-secret-jwt-key-here"
```

## Database Migration

The refresh token table will be automatically created when you run the application. The migration file is located at `migrations/005_refresh_tokens_table.sql`.