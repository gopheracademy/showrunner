BEGIN;

CREATE TABLE job_board(
id SERIAL PRIMARY KEY,
company_name TEXT NOT NULL,
title TEXT NOT NULL,
description TEXT NOT NULL,
link TEXT NOT NULL,
discord TEXT NOT NULL,
rank INT NOT NULL
);

INSERT INTO job_board(
  company_name,
  title,
  description,
  link,
  discord,
  rank
) VALUES (
  'Skynet',
  'T1000 Engineer',
  'Assisting the robot takeover',
  'https://skynet.io/terminator',
  'https://discord.gg/termi-time',
  1
);

INSERT INTO job_board(
  company_name,
  title,
  description,
  link,
  discord,
  rank
) VALUES (
  'Pandora Mining Corp',
  'Avatar Driver',
  'Gaining insight into Pandoras other species',
  'https://pandominer.org/jobs',
  'https://discord.gg/pandora',
  2
);

INSERT INTO job_board(
  company_name,
  title,
  description,
  link,
  discord,
  rank
) VALUES (
  'DC Red Project',
  'Customer Refund Assistant',
  'Assist customers to recieve their refunds',
  'https://syberfunk.gg',
  'https://discord.gg/woops',
  3
);

INSERT INTO job_board(
  company_name,
  title,
  description,
  link,
  discord,
  rank
) VALUES (
  'Slowly',
  'Software Engineer',
  'Slow development styles are a must',
  'https://syberfunk.gg',
  'https://discord.gg/slowly',
  4
);

INSERT INTO job_board(
  company_name,
  title,
  description,
  link,
  discord,
  rank
) VALUES (
  'gopheRPC',
  'Full-time Pancakes walker',
  'Looking for a energetic gopher to support pancakes',
  'https://gopheRPC.gg/pancakes-job',
  'https://discord.gg/gopheRC',
  5
);

COMMIT;