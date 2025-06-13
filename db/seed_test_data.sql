-- db/seed_test_data.sql

-- Seed users
INSERT INTO users (username, email, password_hash, full_name, bio) VALUES
('testuser1', 'testuser1@example.com', 'hashed_password1', 'Test User One', 'Bio of test user one.'),
('testuser2', 'testuser2@example.com', 'hashed_password2', 'Test User Two', 'Bio of test user two.');

-- Seed discussions
INSERT INTO discussions (user_id, title, content) VALUES
(1, 'First Discussion by User1', 'This is the content of the first discussion posted by User1.'),
(2, 'Second Discussion by User2', 'This is the content of the second discussion posted by User2.');

-- Seed comments
INSERT INTO comments (discussion_id, user_id, content) VALUES
(1, 2, 'User2 commenting on User1''s discussion.'),
(2, 1, 'User1 commenting on User2''s discussion.');

-- Seed tags
INSERT INTO tags (name) VALUES
('general'),
('tech'),
('go');

-- Seed discussion_tags
INSERT INTO discussion_tags (discussion_id, tag_id) VALUES
(1, 1), -- First discussion tagged as 'general'
(1, 2), -- First discussion tagged as 'tech'
(2, 3); -- Second discussion tagged as 'go'
