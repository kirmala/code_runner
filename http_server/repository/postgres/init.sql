CREATE TABLE tasks (
    task_id UUID,
    constraint pk_tasks PRIMARY KEY (task_id),
    task_code TEXT,
    task_translator VARCHAR(150) NOT NULL,
    task_result TEXT,
    task_status VARCHAR(150) NOT NULL
    -- user_id UUID NOT NULL,
    -- FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE users(
    user_id UUID,
    constraint pk_users PRIMARY KEY (user_id),
    user_login VARCHAR(150) NOT NULL,
    constraint uq_user_login UNIQUE (user_login),
    user_password VARCHAR(150) NOT NULL
);

CREATE INDEX idx_user_login ON users (user_login);
