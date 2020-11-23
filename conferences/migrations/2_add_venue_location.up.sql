BEGIN;

CREATE TABLE venue(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  address TEXT NOT NULL,
  directions TEXT NOT NULL, 
  google_map_url TEXT NOT NULL,
  capacity int
);

INSERT INTO venue (
  name,
  description,
  address,
  directions,
  google_map_url,
  capacity) VALUES (
    'GopherTown', 
    'The virtual GopherCon meeting space!', 
    'GopherTown, GopherVille, Gophera, 12345', 
    'Do not follow the beavers', 
    'https://gophercon.com', 
    2500
);

INSERT INTO venue (
  name, 
  description, 
  address, 
  directions, 
  google_map_url, 
  capacity) VALUES (
    'Disneyworld', 
    'The Walt Disney World Dolphin Resort', 
    '1500 Epcot Resorts Boulevard, Lake Buena Vista, Florida 32830',
    'Follow the yellow brick road', 
    'https://goo.gl/maps/AwZVpoSXzbzgJkoL8', 
    3000
);

CREATE TABLE location(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL, 
  address TEXT NOT NULL,
  directions TEXT NOT NULL,
  google_map_url TEXT NULL,
  capacity INT NOT NULL,
  venue_id SERIAL NOT NULL REFERENCES venue(id)
);


INSERT INTO location (
  name, 
  description, 
  address, 
  directions, 
  google_map_url, 
  capacity, 
  venue_id
  ) VALUES (
    'The Den', 
    'Suitable for napping', 
    'GopherTown, The Internet, 418', 
    'Left at the crab, right at the penguin, straight down the hole', 
    'https://gophercon.com/gopher-den', 
    450, 
    1
);

INSERT INTO location (
  name, 
  description, 
  address, 
  directions, 
  google_map_url, 
  capacity, 
  venue_id
  ) VALUES (
    'Cinderella Castle', 
    'Where princesses give speeches', 
    'Cinderella Castle, Orlando, Florida 32836, USA', 
    'Down the rabbit hole', 
    'https://goo.gl/maps/G668N9bK14BbryK36', 
    300, 
    2
);


ALTER TABLE conference ADD venue_id SERIAL REFERENCES venue(id);

UPDATE conference SET venue_id = 1 WHERE name = 'GopherCon 2020';

UPDATE conference SET venue_id = 2 WHERE name = 'GopherCon 2021';

ALTER TABLE conference DROP location; 

ALTER TABLE conference ALTER venue_id SET NOT NULL;

ALTER TABLE conference_slot ADD location_id SERIAL REFERENCES location(id);

INSERT INTO conference_slot(
  name, 
  description, 
  cost, 
  capacity, 
  start_date, 
  end_date, 
  purchaseable_from, 
  purchaseable_until, 
  available_to_public, 
  conference_id, 
  location_id
  ) VALUES (
    'Flappy Gopher Round 2', 
    'Will they beat it this time?', 
    50, 
    150, 
    '2021-11-09 17:00:00+00', 
    '2021-11-09 21:00:00+00', 
    '2021-05-09 17:00:00+00', 
    '2021-11-04 17:00:00+00', 
    true, 
    2, 
    2
);

UPDATE conference_slot 
SET location_id = 1 
WHERE name = 'Pre-Conference Workshop: Getting a Jumpstart in Go';

ALTER TABLE conference_slot ALTER location_id SET NOT NULL;

COMMIT;