CREATE TABLE IF NOT EXISTS Category (
                                        id BIGSERIAL NOT NULL,
                                        guid UUID NOT NULL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);