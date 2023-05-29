CREATE TABLE "order" (
                       id SERIAL PRIMARY KEY,
                       user_id INT NOT NULL,
                       status VARCHAR(50) NOT NULL,
                       special_requests TEXT,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW(),
                       FOREIGN KEY (user_id) REFERENCES "user"(id)
);
