-- https://en.wikipedia.org/wiki/Haversine_formula
CREATE OR REPLACE FUNCTION haversine(Lat1 FLOAT(6), Long1 FLOAT(6), Lat2 FLOAT(6), Long2 FLOAT(6)) RETURNS INT
AS $$ SELECT 2 * 6335 * sqrt(pow(sin((radians(Lat2) - radians(Lat1)) / 2), 2) + cos(radians(Lat1)) * cos(radians(Lat2)) * pow(sin((radians(Long2) - radians(Long1)) / 2), 2)) $$
LANGUAGE SQL;

CREATE TABLE IF NOT EXISTS partners
(
    id      INT PRIMARY KEY,
    lat     FLOAT(6) NOT NULL,
    long    FLOAT(6) NOT NULL,
    radius  INT NOT NULL,
    rating  INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS categories
(
    id          INT NOT NULL,
    partner_id  INT NOT NULL REFERENCES partners(id),
    description VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, partner_id)
);

CREATE TABLE IF NOT EXISTS materials
(
    id              INT NOT NULL,
    partner_id      INT NOT NULL REFERENCES partners(id),
    description     VARCHAR(255) NOT NULL,
    PRIMARY KEY (id, partner_id)
);

INSERT INTO partners (id, lat, long, radius, rating) VALUES (1, 1.3, 1.3, 200, 1);
INSERT INTO partners (id, lat, long, radius, rating) VALUES (2, 1.2, 1.2, 200, 3);
INSERT INTO partners (id, lat, long, radius, rating) VALUES (3, 1.1, 1.1, 200, 1);
INSERT INTO partners (id, lat, long, radius, rating) VALUES (4, 1.4, 1.4, 200, 2);
INSERT INTO partners (id, lat, long, radius, rating) VALUES (5, 3.0, 3.0, 200, 5);
INSERT INTO partners (id, lat, long, radius, rating) VALUES (6, 4.0, 4.0, 200, 5);

INSERT INTO categories (id, partner_id, description) VALUES (1, 1, 'Flooring materials');
INSERT INTO categories (id, partner_id, description) VALUES (1, 2, 'Flooring materials');
INSERT INTO categories (id, partner_id, description) VALUES (1, 3, 'Flooring materials');
INSERT INTO categories (id, partner_id, description) VALUES (1, 4, 'Flooring materials');
INSERT INTO categories (id, partner_id, description) VALUES (1, 5, 'Flooring materials');
INSERT INTO categories (id, partner_id, description) VALUES (1, 6, 'Flooring materials');

INSERT INTO materials (id, partner_id, description) VALUES (1, 1, 'Wood');
INSERT INTO materials (id, partner_id, description) VALUES (1, 2, 'Wood');
INSERT INTO materials (id, partner_id, description) VALUES (1, 3, 'Wood');
INSERT INTO materials (id, partner_id, description) VALUES (1, 4, 'Wood');
INSERT INTO materials (id, partner_id, description) VALUES (1, 5, 'Wood');
INSERT INTO materials (id, partner_id, description) VALUES (1, 6, 'Wood');

INSERT INTO materials (id, partner_id, description) VALUES (2, 1, 'Carpet');
INSERT INTO materials (id, partner_id, description) VALUES (2, 2, 'Carpet');
INSERT INTO materials (id, partner_id, description) VALUES (2, 3, 'Carpet');
INSERT INTO materials (id, partner_id, description) VALUES (2, 4, 'Carpet');
INSERT INTO materials (id, partner_id, description) VALUES (2, 5, 'Carpet');
INSERT INTO materials (id, partner_id, description) VALUES (2, 6, 'Carpet');

INSERT INTO materials (id, partner_id, description) VALUES (3, 1, 'Tile');
INSERT INTO materials (id, partner_id, description) VALUES (3, 2, 'Tile');
INSERT INTO materials (id, partner_id, description) VALUES (3, 3, 'Tile');
INSERT INTO materials (id, partner_id, description) VALUES (3, 5, 'Tile');
INSERT INTO materials (id, partner_id, description) VALUES (3, 6, 'Tile');
