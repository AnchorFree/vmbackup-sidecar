[![GitHub license](https://img.shields.io/github/license/AnchorFree/vmbackup-sidecar.svg)](https://github.com/AnchorFree/vmbackup-sidecar/blob/master/LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/AnchorFree/vmbackup-sidecar)](https://goreportcard.com/report/github.com/AnchorFree/vmbackup-sidecar)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/AnchorFree/vmbackup-sidecar?include_prereleases)

vmbackup-sidecar
----------------

Provides a sidecar backup service for [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics/tree/cluster) *vmstorage* component.


## How it works

**vmbackup-sidecar** provides `/backup/create` API endpoint. After receiving `GET` request, it creates *VictoriaMetrics* `vmstorage` snapshot via [snapshot API](https://github.com/VictoriaMetrics/VictoriaMetrics/wiki/Cluster-VictoriaMetrics#url-format) and syncs it with S3 bucket using [aws s3 sync](https://docs.aws.amazon.com/cli/latest/reference/s3/sync.html). When sync with S3 is completed, all `vmstorage` snapshots are removed to free up space. Depending on the number of `vmstorage` instances (pods) running, snapshots will be stored in S3 bucket in following tree:

```
<s3-bucket-name>
├── <vmstorage-pod-1>
│   ├── data
│   └── indexdb
├── <vmstorage-pod-2>
│   ├── data
│   └── indexdb
└── <vmstorage-pod-3>
    ├── data
    └── indexdb
```


## Docker image

```
$ docker pull anchorfree/vmbackup-sidecar:latest
```


## Configuration

Parameters for communicating with `vmstorage` instance and credentials for accessing S3 bucket are stored in environment variables:

```bash
VMSTORAGE_HOST: localhost      # vmstorage hostname (should be available via localhost as both containers are in the same Pod)
VMSTORAGE_PORT: "8482"         # which port vmstorage is listening for API requests
VM_SNAPSHOT_BUCKET: my-bucket  # S3 bucket name for syncing snapshot
VM_STORAGE_DATA_PATH: "/vmstorage-data"  # Corresponds to --storageDataPath flag in VictoriaMetrics setup
ENVIRONMENT: "dev"            # either "prod" or "dev", affects logs output (structlog for prod, plain for dev)
AWS_ACCESS_KEY_ID: <access_key>
AWS_SECRET_ACCESS_KEY: <secret_key>
```

Command-line options:

```bash
./backup-vm -h
Usage of ./backup-vm:
  -help
    Show usage
  -only-show-errors
    Only errors and warnings are displayed. All other output is suppressed
  -port int
    Port to listen (default 8488)
```

## Deployment

**vmbackup-sidecar** [container](https://hub.docker.com/r/anchorfree/vmbackup-sidecar) is supposed to be running in the same [pod](https://kubernetes.io/docs/concepts/workloads/pods/pod/) as *VictoriaMetrics* **vmstorage** component.

Here is a [helm-chart](https://helm.sh/docs/developing_charts/) which implements such kind of **vmbackup-sidecar** deployment: https://github.com/AnchorFree/helm-charts/blob/master/stable/victoria-metrics/templates/vmstorage-statefulset.yaml