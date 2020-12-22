BEGIN;

CREATE TABLE paper_submission(
  id SERIAL PRIMARY KEY,
  user_id SERIAL NOT NULL REFERENCES users(id),
  conference_id SERIAL NOT NULL REFERENCES conference(id),
  title TEXT NOT NULL,
  elevator_pitch TEXT NOT NULL,
  description TEXT NOT NULL,
  notes TEXT NOT NULL
);


INSERT INTO paper_submission (
  user_id,
  conference_id,
  title,
  elevator_pitch,
  description,
  notes
) VALUES (
  1,
  1,
  'The ABCs: Always Be Coding',
  'Taking the world of warcraft mantra, Always Be Casting, into the world of coding.',
  'This is a test description of a talk that will never happen',
  'Or will it?'
);

INSERT INTO paper_submission (
  user_id,
  conference_id,
  title,
  elevator_pitch,
  description,
  notes
) VALUES (
  2,
  1,
  'Putting the g, in gRPC',
  'Naming is important in coding, it shows intention.',
  'They turned down GoatZilla in favour of Godric.',
  'Shame on them.'
);

INSERT INTO paper_submission (
  user_id,
  conference_id,
  title,
  elevator_pitch,
  description,
  notes
) VALUES (
  2,
  2,
  'Do Cats Make Humans Better Engineers?',
  'Cats claim they increase stress relief by 900%',
  'Can we trust them?',
  'Meow'
);

INSERT INTO paper_submission (
  user_id,
  conference_id,
  title,
  elevator_pitch,
  description,
  notes
) VALUES (
  1,
  2,
  'Init functions: not even once',
  'We all know about the humble init function, but did you know its considered bad practice to use it?',
  'Is it because it reminds british Gophers of the 2000s trend of saying INIT BRUV? discuss.',
  'Init bruv'
);


COMMIT;