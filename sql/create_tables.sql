
drop table if exists answers;
drop table if exists users;

-- Enable the uuid-ossp extension for generating UUIDs (if needed)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create table for storing user information with client-generated UUIDs
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE
);

-- Create table for storing answers to questions with client-generated UUIDs
CREATE TABLE answers (
    answer_id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    question0 TEXT,
    question1 TEXT,
    question2 TEXT,
    question3 TEXT,
    question4 TEXT,
    question5 TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Optional: Insert example data into the users table with client-generated UUIDs
-- Example UUIDs provided for demonstration
INSERT INTO users (user_id, username) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'user1'),
('550e8400-e29b-41d4-a716-446655440001', 'user2');

-- Optional: Insert example data into the answers table with client-generated UUIDs
-- Use the same UUIDs used for users in the example above
INSERT INTO answers (user_id, question0, question1, question2, question3, question4, question5) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'Answer1-Q1', 'Answer1-Q2', 'Answer1-Q3', 'Answer1-Q4', 'Answer1-Q5', 'Answer1-Q6'),
('550e8400-e29b-41d4-a716-446655440001', 'Answer2-Q1', 'Answer2-Q2', 'Answer2-Q3', 'Answer2-Q4', 'Answer2-Q5', 'Answer2-Q6');
