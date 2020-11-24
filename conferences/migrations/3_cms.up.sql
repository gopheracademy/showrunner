
ALTER TABLE conference ADD published BOOLEAN;
UPDATE conference SET published = true;
ALTER TABLE conference ALTER published SET NOT NULL;

ALTER TABLE event ADD published BOOLEAN;
UPDATE event SET published = true;
ALTER TABLE event ALTER published SET NOT NULL;

CREATE TABLE site(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);


CREATE TABLE domain (
  id SERIAL PRIMARY KEY,
  fqdn TEXT NOT NULL,
  is_primary BOOLEAN NOT NULL,
  site_id SERIAL NOT NULL REFERENCES site(id)

);

INSERT INTO site(
  name
) VALUES (
  'GopherCon' 
);

INSERT INTO domain (
  fqdn,
  is_primary,
  site_id
) VALUES (
  'www.gophercon.com', 
  true,
  1
);