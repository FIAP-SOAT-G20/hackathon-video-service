DROP TYPE IF EXISTS video_status;
CREATE TYPE video_status AS ENUM ('CREATED','UPLOADED','PROCESSING','FINISHED','FAILED');

CREATE TABLE IF NOT EXISTS videos
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    hash        VARCHAR(64) NOT NULL,
    link        VARCHAR(255) NOT NULL,
    user_id     INT NOT NULL,
    status     video_status DEFAULT 'CREATED',
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO videos (id, user_id, status, created_at)
VALUES (1, 1, 'CREATED', now()),
       (2, 2, 'UPLOADED', now()),
       (3, 3, 'PROCESSING', '2021-10-01 10:00:00.467'),
       (4, 4, 'FINISHED', '2021-10-01 10:00:00.467'),
       (5, 5, 'FAILED', '2021-10-01 10:00:00.467');


SELECT setval('videos_id_seq', (SELECT MAX(id) FROM videos));
