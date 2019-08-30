FROM postgres
ENV POSTGRES_USER hatnotehistory
ENV POSTGRES_PASSWORD hatnotehistory
ENV POSTGRES_DB hatnotehistory
COPY init.sql /docker-entrypoint-initdb.d/10-init.sql
EXPOSE 80
