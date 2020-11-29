BEGIN;

CREATE TYPE sponsorship_level AS ENUM ('platinum', 'gold', 'silver', 'bronze');

CREATE TABLE sponsor(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  address TEXT NOT NULL,
  website TEXT NOT NULL,
  conference_id SERIAL NOT NULL REFERENCES conference(id)
  sponsorship_level sponsorship_level NOT NULL
);

CREATE TYPE role AS ENUM ('marketing', 'logistics', 'technical', 'other', 'sole_contact');

CREATE TABLE sponsor_contact_information(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  role role NOT NULL, 
  email TEXT NULL,
  phone TEXT NULL,
  sponsor_id SERIAL NOT NULL REFERENCES sponsor(id)
);

INSERT INTO sponsor (
  name, 
  address, 
  website, 
  conference_id
  sponsorship_level
) VALUES (
  'Google', 
  'Google Plaza, Google Town, G00 G13', 
  'https://www.google.com', 
  1,
  'platinum'
);

INSERT INTO sponsor (
  name, 
  address, 
  website, 
  sponsorship_level
) VALUES (
  '1Password', 
  'The Crypt, Qwerty Town, 123 456', 
  'https://www.1password.com', 
  1,
  'gold'
);

INSERT INTO sponsor (
  name, 
  address, 
  website, 
  sponsorship_level
) VALUES (
  'Sourcegraph', 
  'Universal House, Code Creek, Search Street, 533728', 
  'https://www.sourcegraph.com', 
  1, 
  'silver'
);

INSERT INTO sponsor (
  name, 
  address, 
  website, 
  sponsorship_level
) VALUES (
  'Sonobi', 
  'So Office, Nobi Building, 500081', 
  'https://www.sonobi.com', 
  1, 
  'bronze'
);

INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Gary Gopher', 
  'sole_contact', 
  'gary@gopher.com', 
  '600613',
  1
);

INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Patrik Passy', 
  'technical', 
  'pat@1password.com', 
  '9455',
  2
);

INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Larry Logipass', 
  'logistics', 
  'larry@1password.com', 
  '888888888',
  2
);

INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Sinn Qulack', 
  'technical', 
  'sinn@sourcegraph.com', 
  '666666',
  3
);


INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Sally Nobey', 
  'marketing', 
  'sally@sonobi.com', 
  '87654321',
  4
);

INSERT INTO sponsor_contact_information (
  name, 
  role, 
  email, 
  phone, 
  sponsor_id
) VALUES (
  'Sammy Boney', 
  'other', 
  'sammy@sonobi.com', 
  '5554325',
  4
);

COMMIT;



