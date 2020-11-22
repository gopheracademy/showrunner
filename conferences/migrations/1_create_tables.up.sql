CREATE TABLE conference( 
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT NOT NULL
);

INSERT INTO conference (name, slug) VALUES ('GopherCon', 'gc');

CREATE TABLE event(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  slug TEXT NOT NULL UNIQUE,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  location TEXT NOT NULL,
  conference_id SERIAL NOT NULL REFERENCES conference(id)
);

INSERT INTO event (name, slug, start_date, end_date, location, conference_id) VALUES ('GopherCon 2020', 'gc-2020', '2020-11-09 17:00:00+00', '2020-11-13 23:45:00+00', 'Online', 1);

INSERT INTO event (name, slug, start_date, end_date, location, conference_id) VALUES ('GopherCon 2021', 'gc-2021', '2021-11-09 17:00:00+00', '2021-11-13 23:45:00+00', 'Florida', 1);

CREATE TABLE event_slot(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  cost INT NOT NULL,
  capacity INT NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  purchaseable_from TIMESTAMPTZ NOT NULL,
  purchaseable_until TIMESTAMPTZ NOT NULL,
  available_to_public BOOLEAN NOT NULL,
  event_id SERIAL NOT NULL REFERENCES event(id)
);

INSERT INTO event_slot(name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public, event_id) VALUES ('Pre-Conference Workshop: Getting a Jumpstart in Go', 'Description goes here', 400, 30, '2020-11-09 17:00:00+00', '2020-11-09 21:00:00+00', '2020-05-09 17:00:00+00', '2020-11-04 17:00:00+00', true, 1);