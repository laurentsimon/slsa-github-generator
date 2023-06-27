#!/bin/bash -eu
#
# Copyright 2023 SLSA Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

mkdir binaries

# Transfer flags and targets to their respective arrays
IFS=' ' read -r -a build_flags <<< "${FLAGS}"
IFS=' ' read -r -a build_targets <<< "${TARGETS}"

# Build with respect to entire arrays of flags and targets
bazel build "${build_flags[@]}" "${build_targets[@]}"

# Use associative array as a set to increase efficency in avoiding double copying the target
declare -A files_set

# Using target string, copy artifact(s) to binaries dir
for curr_target in "${build_targets[@]}"; do
  # Get file(s) generated from build with respect to the target
  bazel_generated=$(bazel cquery --output=starlark --starlark:expr="'\n'.join([f.path for f in target.files.to_list()])" "$curr_target" 2>/dev/null)
  
  # Uses a Starlark expression to pass new line seperated list of file(s) into the set of files
  while read -r file; do
    # Key value is target path, value we do not care about and is set to constant "1"
    files_set["${file}"]="1"
  done <<< "$bazel_generated"
done

# Copy set of unique targets to binaries. Without !, it would give values not keys
# TODO(Issue #2331): switch copy to binaries to a temp dir
for file in "${!files_set[@]}"; do
  # Remove the symbolic link and copy
  cp -L "$file" ./binaries
done
