### How to add migrations
1. Create a new file sql file in this directory.
   Migrations are run sorted  lexicographically.
   If the last migration was `037_something.up.sql` you need to add `038_new_thing.up.sql`
2. WARNING: DO NOT create and deploy migrations that alter a table schema and
   update data in a table in the same migration set.  Altering the table schema
   will cause tables to be locked while the entire migration transaction is
   pending, causing errors when users try to read from the table.
4. Pair your `.up.sql` with a corresponding `.down.sql` that undoes the change.

### Manually Update Your Local Dev with the latest version
1. In a terminal session, make sure you are in `assets/migrations` and run
```
migrate -verbose -path .  -database "postgres://postgres@localhost:5432/khan_dev?sslmode=disable" up 1
```
see `# Setup golang-migrate` below if you don't have migrate installed.

2. Manually validate that the upgrade successfully altered the database in the expected way
```
psql khan_dev -U postgres
```

3. Test the downgrade by running
```
migrate -verbose -path .  -database "postgres://postgres@localhost:5432/khan_dev?sslmode=disable" down 1
```
Note: Migrations can often have subtle side effects, some of which are only
visible on production-like data! To test your migration against production-like
data you can create and use a staging environment.
See `## Setting up a staging environment` at the end of this file for how to do this.

## More information below here:

## Database Environments
A schema contains a group of tables.
A database contains a group of schemas.

We can use the `postgres` user in all three environments, but in prod we want to vary the user to better
differentiate the access pattern metrics.

The schemas for all three should be kept the same.

We do this with migrations.

## Migrations

This directory holds the migrations for our reports database (as well as the
runner that will run those migrations).

## Setup golang-migrate

You need to have [golang-migrate CLI installed](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
 (e.g. `brew install golang-migrate`)

The `golang-migrate` provides a CLI can be used to manually perform the database migrations.

A complete tutorial [is available here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md).
