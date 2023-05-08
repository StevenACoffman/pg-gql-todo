-- pgcrypto has gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- uuid-ossp has uuid_generate_v4()
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TYPE user_role AS ENUM ('guest', 'member', 'admin');
--
-- CREATE TABLE IF NOT EXISTS users (
--                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--                        name TEXT NOT NULL,
--                        role user_role NOT NULL DEFAULT 'guest',
--                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

create table IF NOT EXISTS todo (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       description TEXT NOT NULL,
                       done bool NOT NULL DEFAULT false,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       last_modified_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--                        user_id UUID NOT NULL,