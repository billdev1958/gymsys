FROM postgres:16.3

ENV POSTGRES_USER=ax01
ENV POSTGRES_PASSWORD=1a2s3d4f
ENV POSTGRES_DB=gym_sys

COPY gymSys.sql /docker-entrypoint-initdb.d/

EXPOSE 5432:5432

CMD ["postgres"]
