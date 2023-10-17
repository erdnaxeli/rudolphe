CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    name VARCHAR
);

CREATE TABLE scores (
    id INTEGER PRIMARY KEY,
    year INTEGER,
    score INTEGER,
    user_id VARCHAR REFERENCES users(id),
    UNIQUE (year, user_id)
);

CREATE TABLE days (
    day INTEGER,
    parts INTEGER,
    score_id INTEGER REFERENCES scores(id),
    PRIMARY KEY (day, score_id)
);
