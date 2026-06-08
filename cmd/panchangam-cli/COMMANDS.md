# Panchangam CLI Command Reference

This file documents only commands that are currently registered.

## Global Flags

- `-s, --server`: service address, default `localhost:8080`
- `-o, --output`: output format, default `table`
- `-t, --timeout`: request timeout, default `10s`
- `-v, --verbose`: verbose output
- `--debug`: debug output

## Commands

### `get`

Fetch basic Panchangam data from the configured service.

```bash
panchangam-cli get -l london
panchangam-cli get --lat 37.7749 --lon -122.4194 --tz America/Los_Angeles
```

### `tithi`

Calculate Tithi for a date and location.

```bash
panchangam-cli tithi -l mumbai
panchangam-cli tithi -d 2024-06-21 --lat 19.0760 --lon 72.8777 --tz Asia/Kolkata
```

### `sun`

Calculate sun timing information.

```bash
panchangam-cli sun -l tokyo
panchangam-cli sun -l london --detailed
```

### `locations`

List predefined locations.

```bash
panchangam-cli locations
```

### `validate`

Check service connectivity.

```bash
panchangam-cli validate
```

### `benchmark`

Run simple service performance checks.

```bash
panchangam-cli benchmark -n 100 -w 10
```

### `version`

Print version details.

```bash
panchangam-cli version
panchangam-cli version -o json
```

### `health`

Print local CLI health details.

```bash
panchangam-cli health
panchangam-cli health -o json
```
