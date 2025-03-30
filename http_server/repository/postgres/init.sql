CREATE TABLE tasks (
    task_id UUID PRIMARY KEY,
    task_code TEXT,
    task_translator VARCHAR(150) NOT NULL,
    task_result TEXT,
    task_status VARCHAR(150) NOT NULL
    -- user_id UUID NOT NULL,
    -- FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE users(
    user_id UUID PRIMARY KEY,
    user_login VARCHAR(150) NOT NULL,
    user_password VARCHAR(150) NOT NULL
);

CREATE INDEX idx_user_login ON users (user_login);
