/* -------------------------------------------------------------
   1️⃣  Create the enum type for role (compatible with all PG versions)
   ------------------------------------------------------------- */
DO $$
BEGIN
    -- If the enum does NOT exist yet, create it
    IF NOT EXISTS (
        SELECT 1
        FROM   pg_type t
        JOIN   pg_namespace n ON n.oid = t.typnamespace
        WHERE  t.typname = 'user_role'
          AND  n.nspname = 'public'   -- change if you use a different schema
    ) THEN
        CREATE TYPE user_role AS ENUM ('user', 'admin');
    END IF;
END $$;

/* -------------------------------------------------------------
   2️⃣  Enable UUID generation (run only once)
   ------------------------------------------------------------- */
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

/* -------------------------------------------------------------
   3️⃣  Create the users table
   ------------------------------------------------------------- */
CREATE TABLE IF NOT EXISTS users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    role       user_role   NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
