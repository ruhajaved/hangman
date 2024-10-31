CREATE TABLE game_sessions (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    guesses_left INT NOT NULL,
    correct_guesses BOOLEAN[] NOT NULL DEFAULT '{}'
);