# Simple Backup

A tool aimed to simplify process of backing up files and MySQL databases. Postgres is coming soon.

## Download

Look at [releases](https://github.com/ashep/sbk/releases) section.

## Configuration

Coming soon. Until then feel free to use [config schema](config/config.schema.json)
and [config.sample.yaml](config.sample.yaml) as a reference.

## Changelog

### 23.06.10.2

- Added debug mode.
- Fixed logging.

### 23.06.10.1

- Fixed `mysqldump` freezing.
- Added time rounding in reports.
- Added MySQL dump size to report.

### 23.06.10

- Added `--skip-lock-tables` to `mysqldump` call.
- Fixed order of phases: db first, then files.
- Improved logging.
- Default `log_dir` path set to `/var/log/backup`

### 23.06.09.1

Fix Telegram token length validation.

### 23.06.09

Initial release.

## To Do

- Perform check for necessary tools installed at startup.
- Log `gzip` output to file.

## Authors

- Oleksandr Shepetko
