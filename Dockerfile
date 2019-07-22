FROM mongo:4.0
COPY ./tools/*.js /docker-entrypoint-initdb.d/