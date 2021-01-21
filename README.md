# Prometheus Moodle Exporter

Exports Moodle Application Metrics:

- moodle_live_users
- moodle_expected_upcoming_partipicants
- moodle_bytes_assign_submission (Labels: course ID)
- moodle_bytes_backup (Labels: course ID)
- moodle_bytes_backup_auto (Labels: course ID)
- moodle_bytes_total

Exposes metrics on port 2345 by default.

Currently configured by standard [PostgreSQL environment variables:](https://www.postgresql.org/docs/current/libpq-envars.html)
- PGHOST
- PGUSER
- PGPASSWORD
- PGDATABASE



## Development

```go build```
```export PGHOST=.. PGUSER=.. PGPASSWORD=.. PGDATABASE=..```
```./moodle_exporter```
