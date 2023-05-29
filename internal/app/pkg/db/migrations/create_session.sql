CREATE TABLE session (
                         id SERIAL PRIMARY KEY,
                         user_id INT NOT NULL,
                         session_token VARCHAR(255) NOT NULL,
                         expires_at TIMESTAMP NOT NULL,
                         FOREIGN KEY (user_id) REFERENCES "user"(id)
);
