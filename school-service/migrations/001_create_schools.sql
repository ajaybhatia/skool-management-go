CREATE TABLE IF NOT EXISTS
    schools (
        id SERIAL PRIMARY KEY,
        registration_number VARCHAR(50) UNIQUE NOT NULL,
        NAME VARCHAR(255) NOT NULL,
        address TEXT,
        phone VARCHAR(20),
        email VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );


CREATE INDEX IF NOT EXISTS idx_schools_registration_number ON schools (registration_number);


CREATE INDEX IF NOT EXISTS idx_schools_name ON schools (NAME);


CREATE INDEX IF NOT EXISTS idx_schools_email ON schools (email);
