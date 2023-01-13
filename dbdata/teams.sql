CREATE TABLE TEAMS(
    ID serial primary key,
    NAME VARCHAR(100) NOT NULL CHECK(length(name)>3),
    CONFERENCE VARCHAR(100) NOT NULL CHECK(conference in ('East', 'West')),
    DIVISION VARCHAR(100) NOT NULL CHECK(division in ('Atlantic', 'Pacific', 'Southeast', 'Central', 'Northwest', 'Southwest')),
    EST_YEAR INT NOT NULL CHECK(est_year > 1940));


COPY teams(NAME, CONFERENCE, DIVISION, EST_YEAR) FROM '/docker-entrypoint-initdb.d/teams.csv' CSV HEADER;