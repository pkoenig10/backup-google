# backup-google

[![](https://github.com/pkoenig10/backup-google/actions/workflows/ci.yml/badge.svg)][actions]

A backup service using [Google Cloud Storage](https://cloud.google.com/storage).

Files are selected using the provided [configuration](#configuration) and backed up using a [service account](https://cloud.google.com/iam/docs/service-accounts).

## Configuration

Configuration is provided using command-line flags and a YAML configuration file.

Detailed usage information is available using the `-help` flag.

### Credentials

If a credentials file is not provided, service account credentials will be retrieved using [Google Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/production).

### Configuration file

- `files`

    Glob patterns of files to backup. The [gobwas/glob](https://github.com/gobwas/glob) syntax is used to compile patterns.

#### Example

```yaml
files:
  - foo.txt
  - foo/*/bar.txt
  - foo/**/baz.txt
```

[actions]: https://github.com/pkoenig10/backup-google/actions
