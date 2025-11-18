FROM migrate/migrate:latest

COPY ./deployments/migrations /migrations

ENTRYPOINT ["migrate", "-path", "/migrations", "-database"]
CMD ["postgres://panchangam:panchangam123@postgres:5432/panchangam?sslmode=disable", "up"]
