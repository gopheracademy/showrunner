BEGIN;
ALTER TABLE conference ADD current BOOLEAN;
UPDATE conference SET current = false;
UPDATE conference SET current = true WHERE slug='gc-2021';
COMMIT;



