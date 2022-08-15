CREATE TABLE ROLES(
    ID serial primary key,
    NAME VARCHAR(100) NOT NULL
);

CREATE TABLE USERS(
    ID serial primary key,
    USERNAME VARCHAR(100) NOT NULL,
    PASSWORD VARCHAR(100) NOT NULL,
    ROLEID INT,
    FAV_PLAYERID INT,
    FAB_TEAMID INT,
    CONSTRAINT FK_ROLE
      FOREIGN KEY(ROLEID) 
	  REFERENCES ROLES(ID)
);

INSERT INTO ROLES(NAME) VALUES ('Admin');
INSERT INTO ROLES(NAME) VALUES ('Regular');
INSERT INTO ROLES(NAME) VALUES ('Premium');

INSERT INTO USERS VALUES (DEFAULT, 'Jack', 'JackPass', 1, null, null);
INSERT INTO USERS VALUES (DEFAULT, 'Steve', 'StevePass', 2, null, null);
INSERT INTO USERS VALUES (DEFAULT, 'Mike', 'MikePass', 3, null, null);