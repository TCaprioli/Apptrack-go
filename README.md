# AppTrack

**AppTrack** is a REST API that manages CRUD operations for job applications, with built-in user authentication. All job application endpoints are secured by an authentication middleware to ensure authorized access.

## Features

- **CRUD Operations**: Manage job applications through the API (Create, Read, Update, Delete).
- **User Authentication**: Authentication for users using Paseto.
- **Protected Routes**: Job application endpoints are secured by an auth middleware.
  
## Tech Stack

- **Backend**: Go
- **Database**: PostgreSQL
- **Authentication**: JWT

## Endpoints

### Authentication

- **POST** `/uers/login` - User login
- **POST** `/users/register` - User registration
- **POST** `/users/me` - User verification

### Job Applications

- **GET** `/applications` - Retrieve all job applications (requires authentication)
- **POST** `/applications` - Create a new job application (requires authentication)
- **GET** `/applications/:id` - Get a specific job application (requires authentication)
- **PUT** `/applications/:id` - Update an existing job application (requires authentication)
- **DELETE** `/applications/:id` - Delete a job application (requires authentication)

## Planned Features

- **Session Management**: Future updates will introduce session management to enhance user experience and security by maintaining active user sessions across requests.
- **Password Recovery**: Integrate go-mail to send secure tokens for resetting forgotten passwords.
- **Third-party Authorization**: Enable users to log in using their Google and LinkedIn accounts.
