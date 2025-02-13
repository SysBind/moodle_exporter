# Prometheus Moodle Exporter

Exports Moodle Application Metrics:

- moodle_live_users
- moodle_expected_upcoming_partipicants
- moodle_bytes_assign_submission (Labels: course ID)
- moodle_bytes_backup (Labels: course ID)
- moodle_bytes_backup_auto (Labels: course ID)
- moodle_bytes_total

Common Labels:
- moodle - Moodle's instance short name

Exposes metrics on port 2345 by default.

Currently configured by standard [PostgreSQL environment variables:](https://www.postgresql.org/docs/current/libpq-envars.html)
- PGHOST
- PGUSER
- PGPASSWORD

Will scan all databases containing moodle installation

_Note_: .pgpass file can be used instead of PGPASSWORD


## Docker Images
Automated build available: ghcr.io/sysbind/moodle_exporter:0.9
(see latest release under _Releases_)

## Troubleshooting
Add the DEBUG=1 enviromnet variable, 
moodle_exporter will sleep 1 Hour on database error before exiting,
Allowing you to exec into the container to test, examine vars, etc.
(Or to attach ephemeral debuggin container)


## Development

### Local
```go build```
```export PGHOST=.. PGUSER=.. PGPASSWORD=..```
```./moodle_exporter```

### Local K8S
- ```kind create cluster --name moodleexporter```
- ```tilt up```
