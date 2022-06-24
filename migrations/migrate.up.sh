migrate \
  -path "${PWD}" \
  -database "postgres://erikrios:erikrios@localhost:5432/moot_db?sslmode=disable" \
  up
