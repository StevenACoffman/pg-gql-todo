Code-Owner: Districts
### How to add migrations
1. Create a new file sql file in this directory.
   Migrations are run sorted  lexicographically.
   If the last migration was `037_something.up.sql` you need to add `038_new_thing.up.sql`
2. WARNING: DO NOT create and deploy migrations that alter a table schema and
   update data in a table in the same migration set.  Altering the table schema
   will cause tables to be locked while the entire migration transaction is
   pending, causing errors when users try to read from the table.
3. The sql_export_incremental pipeline to BigQuery requires ALL tables to have primary keys and last_updated fields!!!
4. Pair your `.up.sql` with a corresponding `.down.sql` that undoes the change.
5. DON'T FORGET TO: Update ../sqldb/helpers.go version uint to the latest version.
6. DON'T FORGET TO: Run `make sqlgen` from the services/progress-reports directory

### Update Your Local Dev with the latest version
1. In a terminal session, make sure you are in `$HOME/khan/webapp/services/progress-reports/migrations` and run
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

## Migrate the Production Database
NOTE: Running migrations in prod is super-dangerous! Measure twice, cut once!

1. Deploy your code to production.
2. Make sure to not run migrations that alter tables other than adding a nullable field (therefore locking them) during busy hours
3. Only run one update at a time.
4. Make sure to checkout and pull master with your deployed changes.
5. If there are changes to the report kubernetes jobs (ie the jobs in /services/progress-reports/cmd)
   that are affected by this schema change, pause them before the migration if neccesary.
   https://console.cloud.google.com/kubernetes/workload/overview?project=khan-academy&pageState=(%22savedViews%22:(%22i%22:%222413a70a88334f099ac4646ee089291c%22,%22c%22:%5B%5D,%22n%22:%5B%22report-rollup%22%5D))
6. in one terminal run
Old cloud_sql_proxy version (pre 2.0)
```
cloud_sql_proxy -instances=khan-academy:us-central1:khan-production=tcp:5433
```
New cloud-sql-proxy version (post 2.0)
```
cloud-sql-proxy khan-academy:us-central1:khan-production --port 5433
```
7. Then, in another terminal make sure you are in `$HOME/khan/webapp/services/progress-reports/migrations` and run:

READ steps 1-6 BEFORE proceeding and READ steps 8-11 AFTER updating!!!!

```
export DBPASS="$(keeper --config ${HOME}/.keeper-config.json get TMUHFQL1w24En53n0ugs-Q --format=json | jq -r .password)"
migrate -verbose -path . -database 'postgres://postgres:'"${DBPASS}"'@localhost:5433/reports?sslmode=disable' up 1
```

8. Repeat the above migrate command for an many updates as needed.
9. Update any report kubernetes jobs that need to be updated.
10. Resume any report kubernetes jobs that you paused in step 4.
11. Quit the `cloud-sql-proxy` connection you created in the other terminal!!

## More information below here:

## Database Environments
A schema contains a group of tables.
A database contains a group of schemas.

In prod, we use a cloudsql PostgreSQL instance.
The `reports` DB is for prod only. It's presence should be a gatekeeper to prevent destructive teardowns.

Everywhere else, we assume a local PostgreSQL running on port 5432 with two existing databases.
The `khan_test` DB is for integration tests.
The `khan_dev` DB is for local development.

We can use the `postgres` user in all three environments, but in prod we want to vary the user to better
differentiate the access pattern metrics.

The schemas for all three should be kept the same.

We do this with migrations.

## Migrations

This directory holds the migrations for our reports database (as well as the
runner that will run those migrations).
## Setting up staging environment.

To setup staging environment:
1. Ask in Slack in the `#it` channel for `Cloud SQL Admin` permission for the `khan-academy` project (or find someone with this
   permission)
2. Create a clone of the production database by going to
   https://console.cloud.google.com/sql/instances?project=khan-academy
   and clicking "Create clone" on the `khan-production` database. Call your
   clone `khan-staging`. This will take several minutes. (TODO(steve): automate
   this process with a job that uses `gcloud sql instances clone`)
3. Setup your machine to connect to the staging DB by:
    1. Running cloudsql-proxy using:
      Old Version:
      `cloud_sql_proxy -instances=khan-academy:us-central1:khan-staging=tcp:5433`
      New Version:
      `cloud-sql-proxy khan-academy:us-central1:khan-staging --port 5433`
    2. Change your postgres url to point to port localhost:5433
4. Ensure that you can connect to the staging datastore by:
    1. Opening a Cloud Console session in Google Cloud Platform (it's one of
       the icons in the top-right)
    2. Running `gcloud sql connect khan-staging --user=postgres` and providing
       the password from
       https://keepersecurity.com/vault/#detail/TMUHFQL1w24En53n0ugs-Q. (You'll
       need access to the `Production DB` Keeper group)
    3. Running `\connect <insert_database_name>`
6. Run your migration against staging by running `migrate -verbose -path .  -database postgres://postgres@localhost:5433/khan_dev?sslmode=disable goto <V>`, where `V` is the number like 6, if you want to test a
   specific new migration version number in isolation)
7. Validate:
    1. That your migration ran successfully (and in a reasonable amount of
       time) (and didn't spike the CPU load on the DB machine too much in the
       process - CPU load is visible in the "Overview "tab mentioned above)
    2. That your migration had the desired effect (by running SQL commands in
       your Cloud Console session)
8. Exit the cloudsql-proxy you started above!
8. Once you're done, cleanup the staging database by going to
    https://console.cloud.google.com/sql/instances?project=khan-academy, and
    selecting the "Delete" option for the `khan-staging` database.

## Setting up local connections to staging environment

You can setup staging connection by using local cloudSQL proxy.  The
[official documentation](https://cloud.google.com/sql/docs/postgres/connect-docker)
have futher info.

To setup you will require the service account credentials at
https://keepersecurity.com/vault/#detail/6zo9PKxlTSLAJKwuf7jntg

To run this, you will need to run a docker container:

    docker run -d \
    -v `pwd`/service-account.json:/config \
    -p <local_port (e.g. 5433)>:5432 \
    gcr.io/cloudsql-docker/gce-proxy:1.12 /cloud_sql_proxy \
    -instances=khan-academy:us-central1:<staging DB name>=tcp:0.0.0.0:5432 -credential_file=/config

Then you can connect to remote staging database with with localhost:5433

If you want to run it locally, you can install it with this:
```
go install github.com/GoogleCloudPlatform/cloudsql-proxy/cmd/cloud_sql_proxy
```

## Rollbacks
Every golang-migrate migration file that ends with `.up.sql` has a defined rollback behavior in a file that ends in `.down.sql`. ***You should test
that the rollback works!*** This is most easily done by running
`migrate -verbose -path .  -database postgres://postgres@localhost:5432/khan_dev?sslmode=disable down 1`, and:
1. Seeing that the downgrade succeeds
2. Manually validating that the downgrade successfully reverted the database
back to a pre-migration state

If your rollback isn't trivially simple and ***absolutely*** safe, then you
should follow the above steps to test-run it against realistic data in a
staging environment!


# Aborting/reverting migrations
Trouble-causing migrations can be aborted. Note that migrations are
performed transactionally (even if there are multiple migrations being run!),
so partial migrations should be impossible.

If necessary, already-completed migrations can be reverted. You will have
already tested this process during the development process, by following the
steps in the Rollbacks section above!
(note running golang-migrate with `down 1` will only revert the most recent migration. You may need to run
this multiple times if you need to revert several migrations!)

If you're ever in doubt about the current state of migrations, the
`schema_migrations` table in the database should always contain a single version
indicating the current migration.

#### Stop Rollup and Manually deploying Rollup, then restarting it:
WARNING:  two rollup jobs should NEVER run concurrently it will lead to loss or over counting of learning time.
Make sure if you want to run the job manually, you stop the cron job until it is fully finished then  you can re-start the job.
You can also just stop the cron job, do the deploy and then restart the cron job without ever manually running it.


See https://khanacademy.atlassian.net/wiki/spaces/DIST/pages/1925480456/Report+Rollup+Update+Procedure on how to deploy updated rollup jobs.

#### Update [current migration version in helpers.go](https://github.com/Khan/webapp/blob/master/services/progress-reports/sqldb/helpers.go#L46 "https://github.com/Khan/webapp/blob/master/services/progress-reports/sqldb/helpers.go#L46")

Edit that constant in that file and make it match your new number. You have to do another deploy. ![:disappointed:](https://pf-emoji-service--cdn.us-east-1.prod.public.atl-paas.net/standard/a51a7674-8d5d-4495-a2d2-a67c090f5c3b/64x64/1f61e.png)

## Setup golang-migrate

You need to have [golang-migrate CLI installed](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
 (e.g. `brew install golang-migrate`)

The `golang-migrate` provides a CLI can be used to manually perform the database migrations.

A complete tutorial [is available here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md).


## Permissions you need access to to migrate db in production
You need access to [this secret](https://keepersecurity.com/vault/#detail/TMUHFQL1w24En53n0ugs-Q) and the `cloudsql.instances.connect`, `cloudsql.instances.get` and `cloudsql.instances.login` permissions:
