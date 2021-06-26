-- insert users with 'pass' password
INSERT INTO users(user_id, email, name, password, bio, image, created_at, updated_at)
VALUES (1, 'user1@gmail.com', 'user1', '$2a$10$1l4d3affxTcLVeWpB2z2hOIMzgOYLCoCnx4V5IO.sbL3c0zsnIdU2', 'user1 bio', 'user1 image', now(), now()),
       (2, 'user2@gmail.com', 'user2', '$2a$10$1l4d3affxTcLVeWpB2z2hOIMzgOYLCoCnx4V5IO.sbL3c0zsnIdU2', 'user2 bio', 'user2 image', now(), now());