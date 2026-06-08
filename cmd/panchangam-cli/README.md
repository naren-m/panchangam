# Panchangam CLI

The CLI exposes only commands that are implemented today.

## Build

```bash
go build -o panchangam-cli ./cmd/panchangam-cli
```

Run without building:

```bash
go run ./cmd/panchangam-cli [command]
```

## Commands

| Command | Purpose |
|---------|---------|
| `get` | Fetch basic Panchangam data from the service |
| `tithi` | Calculate local Tithi details |
| `sun` | Calculate sunrise, sunset, solar noon, and day length |
| `locations` | List predefined city locations |
| `validate` | Check service connectivity |
| `benchmark` | Run simple service performance checks |
| `version` | Print CLI and API version details |
| `health` | Print local CLI health details |

## Examples

```bash
./panchangam-cli tithi -l mumbai
./panchangam-cli sun -l london --detailed
./panchangam-cli get -l tokyo -o json
./panchangam-cli locations
./panchangam-cli validate
```

Use `./panchangam-cli [command] --help` for command-specific flags.
