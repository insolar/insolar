FROM postgres:11
MAINTAINER Alexander Alexeev <alexander.alexeev@insolar.io>

COPY ./initdb-conf.sh /docker-entrypoint-initdb.d/initdb-conf.sh
COPY ./postgresql.conf /docker-entrypoint-initdb.d/postgresql.conf
