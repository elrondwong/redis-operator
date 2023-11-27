#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
SCRIPT_ROOT="${SCRIPT_DIR}/.."
CODEGEN_PKG=$SCRIPT_DIR
source "${CODEGEN_PKG}/kube_codegen.sh"

client-gen -v 0 --go-header-file hack/boilerplate.go.txt --clientset-name versioned --input-base github.com/elrondwong/redis-operator/api --output-package github.com/elrondwong/redis-operator/generated/clientset --input=v1beta2
lister-gen -v 0 --go-header-file hack/boilerplate.go.txt  --input-dirs github.com/elrondwong/redis-operator/api/v1beta2 --output-package github.com/elrondwong/redis-operator/generated/clientset
informer-gen -v 0 --go-header-file hack/boilerplate.go.txt  --input-dirs github.com/elrondwong/redis-operator/api/v1beta2 --output-package github.com/elrondwong/redis-operator/generated/clientset

kube::codegen::gen_client \
    --with-watch \
    --input-pkg-root github.com/elrondwong/redis-operator/api/v1beta2 \
    --output-pkg-root github.com/elrondwong/redis-operator/generated \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../../.." \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
