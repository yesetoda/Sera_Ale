CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('applicant', 'company'))
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    applicant_id UUID REFERENCES users(id),
    job_id UUID REFERENCES jobs(id),
    resume_link TEXT NOT NULL,
    cover_letter VARCHAR(200),
    status VARCHAR(20) NOT NULL CHECK (status IN ('Applied', 'Reviewed', 'Interview', 'Rejected', 'Hired')),
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(applicant_id, job_id)
); 