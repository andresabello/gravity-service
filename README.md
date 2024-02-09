# Migrations
### Make sure migrate CLI is installed in the system.
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
### Create new table migration. Change create_users_table for your table name - create_name_table 
```
migrate create -ext sql -dir internal/database/migrations -seq create_users_table
```
This will generate and up and down file with the given name


### Run generated migration 
```
migrate -database "postgres://developer:secret@postgres:5432/gravity?sslmode=disable" -path internal/database/migrations/  up
```

### Reverse the generated migration 
```
migrate -database "postgres://developer:secret@postgres:5432/gravity?sslmode=disable" -path internal/database/migrations/  down 
```

### Force a specific migration
```
migrate -database "postgres://developer:secret@postgres:5432/gravity?sslmode=disable" -path internal/database/migrations/  force 20231002034516
```

# REST API
## /api/v1/makes
Get all car makes. If you pass a year, then it will only query makes for that specific year.

## /api/v1/makes/{name}/models
Get all car models from a specific make. We need the year once again to only get the models made for that specific year.
