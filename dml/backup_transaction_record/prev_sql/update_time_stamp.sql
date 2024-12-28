UPDATE "public".t_user_info
SET created_at = spanner.timestamptz_subtract(created_at, '8 hours'),
    updated_at = spanner.timestamptz_subtract(updated_at, '8 hours');

UPDATE "public".t_user_transfer
SET created_at = spanner.timestamptz_subtract(created_at, '8 hours');

UPDATE "public".t_user_transfer
SET transferred_at = spanner.timestamptz_subtract(transferred_at, '8 hours')
WHERE transferred_at IS NOT NULL;

UPDATE "public".t_user_login
SET last_login = spanner.timestamptz_subtract(last_login, '8 hours')
WHERE last_login IS NOT NULL;