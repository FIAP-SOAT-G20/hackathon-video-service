DROP TYPE IF EXISTS video_status;
CREATE TYPE video_status AS ENUM ('CREATED','UPLOADED','PROCESSING','FINISHED','FAILED');

CREATE TABLE IF NOT EXISTS videos
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    hash        VARCHAR(64),
    link        VARCHAR(255),
    user_id     INT NOT NULL,
    status     video_status DEFAULT 'CREATED',
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO videos (id, name, description, hash, link, user_id, status, created_at)
VALUES (1, 'Video 1', 'Description 1', 'hash1', 'link1', 1, 'CREATED', now()),
       (2, 'Video 2', 'Description 2', 'hash2', 'link2', 2, 'UPLOADED', now()),
       (3, 'Video 3', 'Description 3', 'hash3', 'link3', 3, 'PROCESSING', '2021-10-01 10:00:00.467'),
       (4, 'Video 4', 'Description 4', 'hash4', 'link4', 4, 'FINISHED', '2021-10-01 10:00:00.467'),
       (5, 'Video 5', 'Description 5', 'hash5', 'link5', 5, 'FAILED', '2021-10-01 10:00:00.467');


SELECT setval('videos_id_seq', (SELECT MAX(id) FROM videos));
