Job Platform REST API

Overview
This project is a RESTful API built for a job platform, allowing companies to post job listings and applicants to browse and apply for jobs. It features authentication, role-based access, job searching with pagination, resume uploads, and application status tracking. The backend is developed using Golang with the Gin framework, PostgreSQL (hosted on Neon), JWT for authentication, bcrypt for password hashing, Cloudinary for file storage, and Swagger for API documentation.
Live Demo
The API is hosted on Render and accessible at:

Backend URL: https://sera-ale-4.onrender.com
Swagger Documentation: https://sera-ale-4.onrender.com/swagger/index.html

Technologies Used

Golang (v1.21): High-performance language for scalable REST APIs.
Gin Framework: Fast and lightweight HTTP framework for routing and middleware.
PostgreSQL (Neon): Reliable relational database, hosted on Neon for managed cloud deployment.
JWT: Secure authentication and authorization using JSON Web Tokens.
bcrypt: Password hashing for secure user credential storage.
Cloudinary: File storage for resume uploads.
Swagger: API documentation for easy testing and integration.
GORM: ORM for streamlined database interactions.
GoMock: Unit testing with mocked database interactions.

Prerequisites
To run locally, install:

Golang (v1.21+): Install Golang
PostgreSQL: Local or Neon instance (Neon Signup)
Cloudinary Account: Cloudinary Signup
Git: For cloning the repository

Setup Instructions

Clone the Repository
git clone https://github.com/yourusername/job-platform-api.git
cd job-platform-api


Install Dependencies
go mod tidy


Set Up the Database

Create a Neon PostgreSQL database and get the connection string.
Run the schema:psql -h <neon-host> -U <username> -d <database> -f schema.sql


Insert roles:INSERT INTO roles (name) VALUES ('applicant'), ('company');




Configure Environment VariablesCreate a .env file:
DATABASE_URL=postgresql://<username>:<password>@<neon-host>/<database>?sslmode=require
JWT_SECRET=your_jwt_secret_key_here
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
PORT=8080


Run the Application Locally
go run main.go

Access at http://localhost:8080.


Running the Application

Hosted Version: Use the live API at https://sera-ale-4.onrender.com.
Local Version: Follow the setup instructions above.

API Endpoints
Key endpoints (see Swagger for full details):

POST /signup: Register a user.
Request: { "name": "John Doe", "email": "john@example.com", "password": "Password123!", "role": "applicant" }
Response: { "success": true, "message": "User created", "object": { "id": "uuid", "name": "John Doe", "email": "john@example.com", "role": "applicant" } }


POST /login: Authenticate a user.
Request: { "email": "john@example.com", "password": "Password123!" }
Response: { "success": true, "message": "Login successful", "object": { "token": "jwt_token" } }


GET /jobs: Browse jobs with filters.
Query: ?title=engineer&page=1&size=10



Access Swagger Documentation

Live: https://sera-ale-4.onrender.com/swagger/index.html
Local: http://localhost:8080/swagger/index.html

Deployment
Hosted on Render:

Database: Neon PostgreSQL.
Configuration: Environment variables set in Render dashboard.
Process: Linked to GitHub for automatic deployments.

Tags
#Golang #RESTAPI #PostgreSQL #JWT #Swagger #Cloudinary #UnitTesting #GinFramework #Neon #Render
