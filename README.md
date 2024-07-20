# Simple Backup

A tool aimed to simplify backing-up files and MySQL databases.

## Download

Look at [releases](https://github.com/ashep/sbk/releases) section.

## Configuration

Coming soon. Until then, feel free to use [config schema](config/config.schema.json)
and [config.sample.yaml](config.sample.yaml) as a reference.

## Changelog

### 0.2.0

- `mysql.sources.filename` config option added.
- Logging fixed.

### 0.1.0

- Debug mode added.
- Logging fixed.

### 0.0.4

- `mysqldump` freezing fixed.
- Time rounding in reports fixed.
- MySQL dump size reporting added.

### 0.0.3

- `--skip-lock-tables` added to `mysqldump` call.
- Order of phases fixed: db first, then files.
- Logging improved.
- Default `log_dir` path set to `/var/log/backup`

### 0.0.2

Telegram token length validation fixed.

### 0.0.1

Initial release.

## To Do

- Perform check for necessary tools installed at startup.
- Log `gzip` output to file.

## Authors

- [Oleksandr Shepetko](https://shepetko.com)
