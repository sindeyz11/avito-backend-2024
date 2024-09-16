CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employee
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

DO
$$
    BEGIN
        BEGIN
            CREATE TYPE organization_type AS ENUM (
                'IE',
                'LLC',
                'JSC'
                );
        EXCEPTION
            WHEN duplicate_object THEN
                NULL;
        END;
    END
$$;

CREATE TABLE IF NOT EXISTS organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);


DO
$$
    BEGIN
        BEGIN
            CREATE TYPE tender_status AS ENUM (
                'Created',
                'Published',
                'Closed'
                );
        EXCEPTION
            WHEN duplicate_object THEN
                NULL;
        END;
    END
$$;


DO
$$
    BEGIN
        BEGIN
            CREATE TYPE tender_service_type AS ENUM (
                'Construction',
                'Delivery',
                'Manufacture'
                );
        EXCEPTION
            WHEN duplicate_object THEN
                NULL;
        END;
    END
$$;

--- Версионирование тендеров сделал не самым оптимальным способом, но в предложениях исправил это
--- Также стоило добавить updated_at, но решил не добавлять, т.к. нигде не отдается
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


DO
$$
    BEGIN
        BEGIN
            CREATE TYPE author_type AS ENUM ('Organization', 'User');
        EXCEPTION
            WHEN duplicate_object THEN
                NULL;
        END;
    END
$$;


DO
$$
    BEGIN
        BEGIN
            CREATE TYPE bid_status AS ENUM (
                'Created',
                'Published',
                'Canceled'
                );
        EXCEPTION
            WHEN duplicate_object THEN
                NULL;
        END;
    END
$$;


CREATE TABLE IF NOT EXISTS bid
(
    bid_id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4() NOT NULL,
    name           VARCHAR(100) NOT NULL,
    description    VARCHAR(500),
    status         bid_status   NOT NULL DEFAULT 'Created',
    tender_id      UUID         NOT NULL,
    tender_version INT          NOT NULL,
    author_type    author_type  NOT NULL,
    author_id      UUID         NOT NULL,
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
    author_type    author_type                     NOT NULL,
    author_id      UUID                            NOT NULL,
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

---- Вставка данных
INSERT INTO organization (id, name, description, type, created_at, updated_at)
VALUES ('550e8400-e29b-41d4-a716-446655440020', 'Organization 1', 'Description 1', 'LLC', '2024-09-14 14:04:34.528768',
        '2024-09-14 14:04:34.528768'),
       ('550e8400-e29b-41d4-a716-446655440021', 'Organization 2', 'Description 2', 'IE', '2024-09-14 14:04:34.528768',
        '2024-09-14 14:04:34.528768'),
       ('550e8400-e29b-41d4-a716-446655440022', 'Organization 3', 'Description 3', 'JSC', '2024-09-14 14:04:34.528768',
        '2024-09-14 14:04:34.528768'),
       ('550e8400-e29b-41d4-a716-446655440023', 'Organization 4', 'Description 4', 'LLC', '2024-09-14 14:04:34.528768',
        '2024-09-14 14:04:34.528768');


INSERT INTO employee (id, username, first_name, last_name, created_at, updated_at)
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'user1', 'First1', 'Last1', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440002', 'user2', 'First2', 'Last2', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440003', 'user3', 'First3', 'Last3', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440004', 'user4', 'First4', 'Last4', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440005', 'user5', 'First5', 'Last5', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440006', 'user6', 'First6', 'Last6', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440007', 'user7', 'First7', 'Last7', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440008', 'user8', 'First8', 'Last8', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440009', 'user9', 'First9', 'Last9', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000a', 'user10', 'First10', 'Last10', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000b', 'user11', 'First11', 'Last11', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000c', 'user12', 'First12', 'Last12', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000d', 'user13', 'First13', 'Last13', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000e', 'user14', 'First14', 'Last14', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544000f', 'user15', 'First15', 'Last15', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440010', 'user16', 'First16', 'Last16', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440011', 'user17', 'First17', 'Last17', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440012', 'user18', 'First18', 'Last18', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440013', 'user19', 'First19', 'Last19', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440014', 'user20', 'First20', 'Last20', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440015', 'user21', 'First21', 'Last21', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440016', 'user22', 'First22', 'Last22', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440017', 'user23', 'First23', 'Last23', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440018', 'user24', 'First24', 'Last24', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-446655440019', 'user25', 'First25', 'Last25', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544001a', 'user26', 'First26', 'Last26', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544001b', 'user27', 'First27', 'Last27', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544001c', 'user28', 'First28', 'Last28', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544001d', 'user29', 'First29', 'Last29', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226'),
       ('550e8400-e29b-41d4-a716-44665544001e', 'user30', 'First30', 'Last30', '2024-09-14 14:04:34.509226',
        '2024-09-14 14:04:34.509226');


INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES ('550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440001'),
       ('550e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440032', '550e8400-e29b-41d4-a716-446655440020',
        '550e8400-e29b-41d4-a716-446655440003'),
       ('550e8400-e29b-41d4-a716-446655440033', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440004'),
       ('550e8400-e29b-41d4-a716-446655440034', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440005'),
       ('550e8400-e29b-41d4-a716-446655440035', '550e8400-e29b-41d4-a716-446655440021',
        '550e8400-e29b-41d4-a716-446655440006'),
       ('550e8400-e29b-41d4-a716-446655440036', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440007'),
       ('550e8400-e29b-41d4-a716-446655440037', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440008'),
       ('550e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440022',
        '550e8400-e29b-41d4-a716-446655440009'),
       ('550e8400-e29b-41d4-a716-446655440039', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000a'),
       ('550e8400-e29b-41d4-a716-44665544003a', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000b'),
       ('550e8400-e29b-41d4-a716-44665544003b', '550e8400-e29b-41d4-a716-446655440023',
        '550e8400-e29b-41d4-a716-44665544000c');
