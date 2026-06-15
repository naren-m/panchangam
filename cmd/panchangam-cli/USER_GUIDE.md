# Panchangam CLI User Guide

The CLI is for quick local calculations and service checks.

## Install

```bash
go build -o panchangam-cli ./cmd/panchangam-cli
```

## Common Commands

Get today's Tithi for Mumbai:

```bash
./panchangam-cli tithi -l mumbai
```

Get sun times for London:

```bash
./panchangam-cli sun -l london --detailed
```

Fetch service data for Tokyo:

```bash
./panchangam-cli get -l tokyo -o json
```

List location shortcuts:

```bash
./panchangam-cli locations
```

Check the service:

```bash
./panchangam-cli validate
./panchangam-cli health
```

## Locations

Use `-l` with a location code, or pass `--lat`, `--lon`, and `--tz`.

Common location codes:

- `mumbai`
- `london`
- `tokyo`
- `nyc`
- `sydney`

Run `./panchangam-cli locations` for the full list.

## Output

Most commands support:

- `-o table`
- `-o json`
- `-o yaml`

Use JSON when another script needs to read the output.
