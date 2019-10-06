FROM mongo:4.2
COPY ./tools/*.js /docker-entrypoint-initdb.d/
