CREATE TABLE "user" (
                      id SERIAL PRIMARY KEY,
                      username VARCHAR(50) UNIQUE NOT NULL,
                      email VARCHAR(100) UNIQUE NOT NULL,
                      password_hash VARCHAR(255) NOT NULL,
                      role VARCHAR(10) NOT NULL CHECK (role IN ('customer', 'chef', 'manager')),
                      created_at TIMESTAMP DEFAULT NOW(),
                      updated_at TIMESTAMP DEFAULT NOW()
);
