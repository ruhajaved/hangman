CREATE TABLE game_sessions (
    id          SERIAL  PRIMARY KEY,
    word        TEXT    NOT NULL,
    difficulty  INT     NOT NULL,
    theme       TEXT    NOT NULL
);

INSERT INTO game_sessions (word, difficulty, theme) VALUES ('rabbit', 2, 'animals');
INSERT INTO game_sessions (word, difficulty, theme) VALUES ('mississippi', 3, 'rivers'); 
INSERT INTO game_sessions (word, difficulty, theme) VALUES ('computer', 2, 'technology');
INSERT INTO game_sessions (word, difficulty, theme) VALUES ('mclaren', 2, 'car brands');
