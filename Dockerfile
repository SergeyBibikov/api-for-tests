FROM postgres:latest
COPY ./dbdata/init.sql /docker-entrypoint-initdb.d/
