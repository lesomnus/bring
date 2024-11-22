# bring

Bring things.

## Usage

From [things.yaml](things.yaml),
```shell
$ bring things.yaml
> [1/3|✓] inventory/logo/github.svg • 426.669507ms
> [2/3|✓] inventory/logo/microsoft.svg • 509.403428ms
> [3/3|!] inventory/money/no-such-thing • 830.91µs
>         bring: request get: Get "https://localhost/not-exists": dial tcp 127.0.0.1:443: connect: connection refused

$ rm inventory/logo/github.svg

# Note that the download of `microsoft.svg` was skipped
# because the digest of the local file matches.
$ bring things.yaml
> [1/3|✓] inventory/logo/github.svg • 354.087526ms
> [2/3|=] inventory/logo/microsoft.svg • 1.620357ms
> [3/3|!] inventory/money/no-such-thing • 701.24µs
>         bring: request get: Get "https://localhost/not-exists": dial tcp 127.0.0.1:443: connect: connection refused
```
