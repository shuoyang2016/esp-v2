#!/bin/bash

# Copyright 2019 Google LLC

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Fail on any error.
set -eo pipefail

rm -rf src/go/proto
rm -rf vendor/github.com/envoyproxy/data-plane-api/api
rm -rf vendor/gogoproto
rm -rf vendor/github.com/census-instrumentation/opencensus-proto/gen-go

#TODO(bochun): probably we can programatically generate these.
# HTTP filter common
bazel build //api/envoy/v10/http/common:base_go_proto
mkdir -p src/go/proto/api/envoy/v10/http/common
cp -f bazel-bin/api/envoy/v10/http/common/base_go_proto_/github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/v10/http/common/* src/go/proto/api/envoy/v10/http/common
# HTTP filter service_control
bazel build //api/envoy/v10/http/service_control:config_go_proto
mkdir -p src/go/proto/api/envoy/v10/http/service_control
cp -f bazel-bin/api/envoy/v10/http/service_control/config_go_proto_/github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/v10/http/service_control/* src/go/proto/api/envoy/v10/http/service_control
# HTTP filter path_rewrite
bazel build //api/envoy/v10/http/path_rewrite:config_go_proto
mkdir -p src/go/proto/api/envoy/v10/http/path_rewrite
cp -f bazel-bin/api/envoy/v10/http/path_rewrite/config_go_proto_/github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/v10/http/path_rewrite/* src/go/proto/api/envoy/v10/http/path_rewrite
# HTTP filter backend_auth
bazel build //api/envoy/v10/http/backend_auth:config_go_proto
mkdir -p src/go/proto/api/envoy/v10/http/backend_auth
cp -f bazel-bin/api/envoy/v10/http/backend_auth/config_go_proto_/github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/v10/http/backend_auth/* src/go/proto/api/envoy/v10/http/backend_auth
