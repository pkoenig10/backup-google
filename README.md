# backup-google

[![](https://github.com/pkoenig10/backup-google/actions/workflows/ci.yml/badge.svg)][actions]

A backup service using [Google Cloud Storage](https://cloud.google.com/storage).

Files are selected using the provided [configuration file](#configuration-file) and backed up using a [service account](https://cloud.google.com/iam/docs/service-accounts).

## Configuration

### Credentials

Service account credentials are obtained using [Google Application Default Credentials (ADC)](https://docs.cloud.google.com/docs/authentication/application-default-credentials).

### Environment variables

| Variable | Description | Required? | Default |
|:-|:-|:-:|:-:|
| `BUCKET_NAME` | The Google Cloud Storage bucket name | Yes | - |
| `CONFIG_PATH` | The configation file path. | No | `config.yml` |

### Configuration file

The configuration file is a YAML file with the following properties:

| Property | Object | Description |
|:-|:-:|:-|
| `files` | Array | Glob patterns of files to backup. The [gobwas/glob](https://github.com/gobwas/glob) syntax is used to compile patterns. |

#### Example

```yaml
files:
  - foo.txt
  - foo/*/bar.txt
  - foo/**/baz.txt
```

[actions]: https://github.com/pkoenig10/backup-google/actions
