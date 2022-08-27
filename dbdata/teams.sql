CREATE TABLE TEAMS(
    ID serial primary key,
    NAME VARCHAR(100) NOT NULL,
    CONFERENCE VARCHAR(100) NOT NULL,
    DIVISION VARCHAR(100) NOT NULL,
    EST_YEAR INT NOT NULL);


COPY teams(NAME, CONFERENCE, DIVISION, EST_YEAR) FROM '/docker-entrypoint-initdb.d/teams.csv' CSV HEADER;