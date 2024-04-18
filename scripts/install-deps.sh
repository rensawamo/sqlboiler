#!/bin/sh

# for sqlboiler binary
go install github.com/volatiletech/sqlboiler/v4@v4.12.0
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.12.0

# for sql-migrate binary
go install github.com/rubenv/sql-migrate/...@v1.1.2
