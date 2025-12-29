BEGIN;

CREATE TYPE module_type AS ENUM ('multiple_choice', 'matching_type');

CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    subject_id UUID NOT NULL,
    grade_id UUID NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    type module_type NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    FOREIGN KEY (grade_id) REFERENCES grades(id)
);

CREATE TABLE IF NOT EXISTS questions (
    id UUID PRIMARY KEY,
    module_id UUID NOT NULL,
    content VARCHAR(255) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules(id)
);

CREATE TABLE IF NOT EXISTS question_choices (
    id UUID PRIMARY KEY,
    question_id UUID NOT NULL,
    content VARCHAR(255) NOT NULL,
    is_correct_answer BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

COMMIT;