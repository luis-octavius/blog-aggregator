-- +goose Up 
CREATE TABLE posts (
  id SERIAL PRIMARY KEY, 
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  title TEXT NOT NULL, 
  url TEXT NOT NULL,
  description TEXT NOT NULL, 
  published_at TIMESTAMP NOT NULL,
  feed_id INT NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE(url)
);

-- +goose Down 
DROP TABLE posts;
