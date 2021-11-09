# Copyright 2021 The Kubernetes Authors.
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


all: sample-credential-provider

sample-credential-provider: ./cmd/main.go
	 GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOPROXY=$(GOPROXY) go build \
		-ldflags="-w -s -X 'main.version=$(VERSION)'" \
		-o=sample-credential-provider \
		cmd/main.go

.PHONY: clean

clean:
	rm -rf sample-credential-provider

