CREATE TABLE categories (
    id VARCHAR(50) PRIMARY KEY CHECK (id ~ '^[a-z0-9_-]+$'),   -- 'f1', 'motogp', 'wec' (always lowercase)
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE circuits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    location VARCHAR(200) NOT NULL,
    country VARCHAR(100) NOT NULL,
    length_km NUMERIC(6,3) NOT NULL CHECK (length_km > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,                  -- e.g. 'VER', 'BAG'
    nationality VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE races (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id VARCHAR(50) NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    circuit_id UUID NOT NULL REFERENCES circuits(id) ON DELETE RESTRICT,
    winner_id UUID NULL REFERENCES drivers(id) ON DELETE SET NULL,
    name VARCHAR(200) NOT NULL,
    season INT NOT NULL CHECK (season > 0),
    round INT NOT NULL CHECK (round > 0),
    race_date DATE NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(payload) = 'object'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, season, round)
);

CREATE INDEX idx_races_category_season ON races(category_id, season);
CREATE INDEX idx_races_race_date ON races(race_date DESC, id ASC);   -- matches the list ORDER BY
CREATE INDEX idx_races_circuit_id ON races(circuit_id);
CREATE INDEX idx_races_winner_id ON races(winner_id);
CREATE INDEX idx_races_payload_gin ON races USING GIN (payload);     -- reserved for future payload filters; no current endpoint uses it

-- ---------------------------------------------------------------------------
-- Teams, engines, chassis, cars
-- ---------------------------------------------------------------------------
CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL UNIQUE,
    short_name VARCHAR(50) NOT NULL,
    country VARCHAR(100) NOT NULL,
    founded_year INT NULL CHECK (founded_year BETWEEN 1800 AND 2200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE engines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manufacturer VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,                              -- model designation
    configuration VARCHAR(50) NOT NULL,                      -- e.g. 'V6', 'V4', 'I4'
    displacement_cc INT NOT NULL CHECK (displacement_cc > 0),
    cylinders SMALLINT NOT NULL CHECK (cylinders > 0),
    aspiration VARCHAR(30) NOT NULL
        CHECK (aspiration IN ('NATURALLY_ASPIRATED', 'TURBOCHARGED', 'HYBRID_TURBO')),
    power_kw NUMERIC(6,1) NULL CHECK (power_kw > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (manufacturer, name)
);

CREATE TABLE chassis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    manufacturer VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    material VARCHAR(100) NOT NULL,                          -- e.g. 'Carbon fibre monocoque'
    wheelbase_mm INT NULL CHECK (wheelbase_mm > 0),
    weight_kg NUMERIC(6,1) NULL CHECK (weight_kg > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (manufacturer, name)
);

CREATE TABLE cars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id VARCHAR(50) NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    engine_id UUID NOT NULL REFERENCES engines(id) ON DELETE RESTRICT,
    chassis_id UUID NOT NULL REFERENCES chassis(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,                              -- e.g. 'RB22'
    season INT NOT NULL CHECK (season > 0),
    parts JSONB NOT NULL CHECK (jsonb_typeof(parts) = 'object'),   -- contract in 3.3
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (team_id, category_id, season, name)
);

CREATE INDEX idx_cars_category_season ON cars(category_id, season);
CREATE INDEX idx_cars_team_id ON cars(team_id);
CREATE INDEX idx_cars_engine_id ON cars(engine_id);
CREATE INDEX idx_cars_chassis_id ON cars(chassis_id);
CREATE INDEX idx_cars_list_order ON cars(season DESC, name ASC, id ASC);   -- matches the list ORDER BY

-- ---------------------------------------------------------------------------
-- Articles: layer 1 = metadata linked to an entity, layer 2 = per-language texts
-- ---------------------------------------------------------------------------
CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(20) NOT NULL
        CHECK (entity_type IN ('category', 'circuit', 'driver', 'race', 'team', 'engine', 'chassis', 'car')),
    entity_id VARCHAR(50) NOT NULL,                          -- category slug, or canonical lowercase UUID text
    kind VARCHAR(50) NOT NULL DEFAULT 'general'
        CHECK (kind ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),           -- free classification, e.g. 'overview', 'report', 'history', 'technical'
    published_at TIMESTAMPTZ NULL,                           -- NULL = draft (never served by the API)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_articles_entity_id CHECK (
        (entity_type = 'category' AND entity_id ~ '^[a-z0-9_-]+$')
        OR (entity_type <> 'category'
            AND entity_id ~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'))
);

CREATE INDEX idx_articles_entity ON articles(entity_type, entity_id);
CREATE INDEX idx_articles_published ON articles(published_at DESC, id ASC);

CREATE TABLE article_texts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    language VARCHAR(35) NOT NULL
        CHECK (language ~ '^[a-z]{2,3}(-[a-z0-9]{2,8})*$'),  -- lowercase BCP 47 tag: 'en', 'pt', 'pt-br'
    title VARCHAR(300) NOT NULL,
    summary TEXT NULL,
    body_md TEXT NOT NULL CHECK (length(body_md) <= 200000), -- Markdown source
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (article_id, language)
);
