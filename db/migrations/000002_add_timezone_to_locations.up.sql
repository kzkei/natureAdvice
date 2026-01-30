-- add tz column
ALTER TABLE locations ADD COLUMN IF NOT EXISTS timezone VARCHAR(50);

-- update existing rows with tz
UPDATE locations SET timezone = 'America/Los_Angeles' WHERE state IN ('CA', 'OR', 'WA');
UPDATE locations SET timezone = 'America/Denver' WHERE state IN ('MT', 'WY', 'CO', 'UT', 'AZ', 'NM');

-- handle edge cases
UPDATE locations SET timezone = 'America/Phoenix' WHERE park_code = 'GCNP'; -- no daylight saving

-- set not null constraint
ALTER TABLE locations ALTER COLUMN timezone SET NOT NULL;