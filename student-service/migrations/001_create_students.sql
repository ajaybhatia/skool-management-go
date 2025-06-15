CREATE TABLE IF NOT EXISTS
    students (
        id SERIAL PRIMARY KEY,
        roll_number VARCHAR(20) NOT NULL,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE,
        phone VARCHAR(20),
        date_of_birth DATE,
        address TEXT,
        school_id INTEGER NOT NULL,
        enrollment_date DATE DEFAULT CURRENT_DATE,
        status VARCHAR(50) DEFAULT 'active',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (roll_number, school_id)
    );


CREATE INDEX IF NOT EXISTS idx_students_school_id ON students (school_id);


CREATE INDEX IF NOT EXISTS idx_students_email ON students (email);


CREATE INDEX IF NOT EXISTS idx_students_status ON students (status);


CREATE INDEX IF NOT EXISTS idx_students_roll_number ON students (roll_number);


CREATE INDEX IF NOT EXISTS idx_students_school_roll ON students (school_id, roll_number);
