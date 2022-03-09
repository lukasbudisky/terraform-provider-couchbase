#! /usr/bin/env bash

set -ux
set -o pipefail

ADDRESS="${1:-http://127.0.0.1}"
PORT="${2:-8091}"
MEM_QUOTA="${3:-512}"
INDEX_QUOTA="${4:-512}"
FTS_QUOTA="${5:-512}"
CBAS_QUOTA="${6:-1024}"
EVENTING_QUOTA="${7:-512}"
USERNAME="${8:-Administrator}"
PASSWORD="${9:-123456}"
RETRY="${10:-60}"

function configure_couchbase() {
    curl -v "${ADDRESS}:${PORT}/pools/default" \
        -d clusterName="cluster" \
        -d memoryQuota="${MEM_QUOTA}" \
        -d indexMemoryQuota="${INDEX_QUOTA}" \
        -d ftsMemoryQuota="${FTS_QUOTA}" \
        -d cbasMemoryQuota="${CBAS_QUOTA}" \
        -d eventingMemoryQuota="${EVENTING_QUOTA}"

    curl -v "${ADDRESS}:${PORT}/node/controller/setupServices" \
        -d "services=kv%2Cn1ql%2Cindex%2Cfts"

    curl -v "${ADDRESS}:${PORT}/settings/web" \
        -d port="${PORT}" \
        -d username="${USERNAME}" \
        -d password="${PASSWORD}"

    curl -v -u "${USERNAME}":"${PASSWORD}" "${ADDRESS}:${PORT}/nodes/self/controller/settings" \
    -d "path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&" \
    -d "index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&" \
    -d "cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&" \
    -d "eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&"

    curl -v -u "${USERNAME}":"${PASSWORD}" "${ADDRESS}:${PORT}/settings/indexes" \
        -d "indexerThreads=4" \
        -d "logLevel=verbose" \
        -d "maxRollbackPoints=10" \
        -d "storageMode=plasma" \
        -d "memorySnapshotInterval=150" \
        -d "stableSnapshotInterval=40000"
}

CYCLE=0
while [ "${CYCLE}" -lt "${RETRY}" ]; do
    CYCLE=$(("${CYCLE}"+1))
    STATUS=$(curl -s -o /dev/null -I -w "%{http_code}" "${ADDRESS}:${PORT}/pools")
    if [[ "${STATUS}" == "200" ]]; then
        echo "couchbase up: ${CYCLE}/${RETRY}"
        configure_couchbase
        exit 0
    fi
    echo "wait for couchbase: ${CYCLE}/${RETRY}"
    sleep 1
done

exit 1
