# Copyright 2022 Bryon Baker

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Use the git tag to set the image version etc
VERSION := $(shell git describe --tags --dirty=-modified --always)

clean:
	rm bin/pinger

build:
	go build -o bin/pinger cmd/main.go

package: clean build
	podman build -t pinger:${VERSION} .
	podman tag pinger:${VERSION} pinger:latest
	podman tag pinger:${VERSION} quay.io/bryonbaker/pinger:${VERSION}
	podman tag pinger:${VERSION} quay.io/bryonbaker/pinger:latest
run:
	go run cmd/main.go

test: clean build
	echo "Tsk tsk! \"make test\" is not implemented yet."

all: clean build
