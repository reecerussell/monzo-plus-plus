FROM mysql

COPY schema.sql /docker-entrypoint-initdb.d/

ADD mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

ENV MYSQL_ROOT_PASSWORD=root_password
ENV MYSQL_DATABASE=monzo_plus_plus
ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=password