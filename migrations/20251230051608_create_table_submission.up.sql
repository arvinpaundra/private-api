BEGIN;

CREATE TYPE submission_status AS ENUM ('inprogress', 'submitted', 'canceled');

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY,
    module_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    student_name VARCHAR(255) NOT NULL,
    status submission_status NOT NULL DEFAULT 'inprogress',
    total_questions SMALLINT NOT NULL DEFAULT 0,
    submitted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules(id)
);

CREATE TABLE IF NOT EXISTS submission_answers (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL,
    question_slug VARCHAR(50) NOT NULL,
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (submission_id) REFERENCES submissions(id)
);

COMMIT;