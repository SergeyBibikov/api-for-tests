FROM postgres:latest
COPY ./dbdata/ /docker-entrypoint-initdb.d/
