# goboyle Models

The database uses Postgres

### To setup the DB:

1. Login to postgres shell

	```psql```

2. Create database for SSS

	```create database healness;```

3. Install uuid extension for IDs

	```CREATE EXTENSION IF NOT EXISTS "uuid-ossp";```

4. Run `go test` command to make sure everything is setup correctly, and DB tables get created
