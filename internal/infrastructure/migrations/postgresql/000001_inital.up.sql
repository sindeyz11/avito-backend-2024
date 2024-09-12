CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--------------------------

CREATE TABLE employee
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-------------------------

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);

----------------------------

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
    status          tender_status         NOT NULL DEFAULT 'Created',
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
    bid_id         UUID                  DEFAULT uuid_generate_v4() NOT NULL,
    name           VARCHAR(100) NOT NULL,
    description    VARCHAR(500),
    status         status_enum  NOT NULL DEFAULT 'Created',
    tender_id      UUID         NOT NULL,
    tender_version INT          NOT NULL,
    author_type    VARCHAR(50)  NOT NULL,
    author_id      UUID,
    version        INT          NOT NULL DEFAULT 1,
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (bid_id, version),

    FOREIGN KEY (tender_id, tender_version)
        REFERENCES tender (tender_id, version)
        ON DELETE CASCADE,

    UNIQUE (bid_id, version)
);