DROP TABLE IF EXISTS todotags;
DROP TABLE IF EXISTS highlightcolor;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS tags;

CREATE TABLE todos (
  ID SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  created_on TIMESTAMP NOT NULL DEFAULT NOW(),
  checked_on TIMESTAMP
);

CREATE TABLE tags (
  ID SERIAL PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE highlightcolor (
  ID SERIAL PRIMARY KEY,
  name VARCHAR(20) NOT NULL UNIQUE,
  hex VARCHAR(9) NOT NULL
);

CREATE TABLE todotags (
  todos_id INT NOT NULL,
  tag_id INT NOT NULL,
  highlight_id INT,
  PRIMARY KEY (todos_id, tag_id),
  FOREIGN KEY (todos_id) REFERENCES todos (ID),
  FOREIGN KEY (tag_id) REFERENCES tags (ID),
  FOREIGN KEY (highlight_id) REFERENCES highlightcolor (ID)
);

INSERT INTO 
  highlightcolor (name, hex)
VALUES
  ('Yellow', '#FFF200'),
  ('Red', '#D0312D'),
  ('Green', '#3CB043');

INSERT INTO
  tags (name)
VALUES
  ('GO'),
  ('RaspberryPi');
