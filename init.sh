#!/bin/bash
set -euo pipefail

install -d /tmp/tpm/

sudo swtpm socket --daemon --tpmstate dir=/tmp/tpm/ --tpm2 --ctrl type=tcp,port=2321 --server type=tcp,port=2320 --flags not-need-init --log level=8

# swtpm server socket
socat TCP-LISTEN:2322,fork TCP:127.0.0.1:2320 &
# Control socket
socat TCP-LISTEN:2323,fork TCP:127.0.0.1:2321 &

wait
