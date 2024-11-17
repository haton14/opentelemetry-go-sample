/Users/haton14/ghq/github.com/haton14/hono-otlp-example/docker/mysql/init.sql
CREATE TABLE IF NOT EXISTS book (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

INSERT INTO book (title) VALUES ('The Great Gatsby');
INSERT INTO book (title) VALUES ('1984');
INSERT INTO book (title) VALUES ('To Kill a Mockingbird');
