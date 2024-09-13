CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
    );

CREATE TYPE tender_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
    );

CREATE TABLE IF NOT EXISTS tender
(
    tender_id       UUID                         DEFAULT uuid_generate_v4() NOT NULL,
    name            VARCHAR(100)        NOT NULL,
    description     VARCHAR(500),
    service_type    tender_service_type NOT NULL,
    status          tender_status       NOT NULL DEFAULT 'Created',
    version         INT                 NOT NULL DEFAULT 1,
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    creator_id      UUID                REFERENCES employee (id) ON DELETE SET NULL,
    created_at      TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tender_id, version)
);

CREATE TYPE author_type AS ENUM ('Organization', 'User');

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
    );

CREATE TABLE IF NOT EXISTS bid
(
    bid_id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4() NOT NULL,
    name           VARCHAR(100) NOT NULL,
    description    VARCHAR(500),
    status         bid_status   NOT NULL DEFAULT 'Created',
    tender_id      UUID         NOT NULL,
    tender_version INT          NOT NULL,
    author_type    VARCHAR(50)  NOT NULL,
    author_id      UUID,
    version        INT          NOT NULL DEFAULT 1,
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tender_id, tender_version)
        REFERENCES tender (tender_id, version)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bid_history
(
    bid_id         UUID DEFAULT uuid_generate_v4() NOT NULL,
    name           VARCHAR(100)                    NOT NULL,
    description    VARCHAR(500),
    status         bid_status                      NOT NULL,
    tender_id      UUID                            NOT NULL,
    tender_version INT                             NOT NULL,
    author_type    VARCHAR(50)                     NOT NULL,
    author_id      UUID,
    version        INT                             NOT NULL,
    created_at     TIMESTAMP,
    FOREIGN KEY (tender_id, tender_version)
        REFERENCES tender (tender_id, version)
        ON DELETE CASCADE,
    PRIMARY KEY (bid_id, version),
    UNIQUE (bid_id, version)
);

CREATE TABLE IF NOT EXISTS review
(
    id          UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    bid_id      UUID REFERENCES bid (bid_id) ON DELETE CASCADE,
    description TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);