#!/bin/bash

set -euo pipefail

cat <<EOF >DIFF
15c15,19
< name: SLSA builder delegator
---
> # This is a version of the delegator workflow that requires as few permissions
> # as possible. TRWs may use this workflow so that they may request fewer
> # GITHUB_TOKEN permissions from end-users.
> 
> name: SLSA low-permission builder delegator
100c104
<           slsa-workflow-recipient: "delegator_generic_slsa3.yml"
---
>           slsa-workflow-recipient: "delegator_lowperms-generic_slsa3.yml"
131,133d134
<       # TODO(#2076): Use dynamic GITHUB_TOKEN permissions.
<       contents: write # To release assets.
<       packages: write # To publish to GitHub packages.
EOF

actual_diff=$(
    diff .github/workflows/delegator_generic_slsa3.yml .github/workflows/delegator_lowperms-generic_slsa3.yml \
    || true
)
expected_diff=$(cat ./DIFF)

if [[ "$expected_diff" != "$actual_diff" ]]; then
    echo "Unexpected differences between the delegator workflows"
    echo "$actual_diff"
    exit 1
fi