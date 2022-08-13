CREATE TABLE JOBS(
    ID serial primary key,
    NAME TEXT
);
CREATE TABLE PHONES(
    ID serial primary key,
    NUMBER VARCHAR(20)
);

CREATE TABLE USERS(
    ID serial primary key,
    NAME VARCHAR(100) NOT NULL,
    SURNAME VARCHAR(100) NOT NULL,
    JOB_ID INT,
    PHONE_ID INT,
    CONSTRAINT FK_JOB
      FOREIGN KEY(JOB_ID) 
	  REFERENCES JOBS(ID),
    CONSTRAINT FK_PHONE
      FOREIGN KEY(PHONE_ID) 
	  REFERENCES PHONES(ID)
);

INSERT INTO PHONES (NUMBER) VALUES('9998887711');