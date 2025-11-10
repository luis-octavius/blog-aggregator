-- +goose Up 
CREATE TABLE feeds (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  url TEXT NOT NULL, 
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE(url)
);

-- +goose Down 
DROP TABLE feeds;
