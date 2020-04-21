FROM mysql

# ROOT PASSWORD
ENV MYSQL_ROOT_PASSWORD=v3ryS3c5r3P1ssW4rD

ARG password=m0nz0-Pl5s-pluS
ARG user=monzo

ENV MYSQL_DATABASE=monzo_plus_plus
ENV MYSQL_USER=monzo
ENV MYSQL_PASSWORD=m0nz0-Pl5s-pluS

COPY schema.sql /docker-entrypoint-initdb.d/