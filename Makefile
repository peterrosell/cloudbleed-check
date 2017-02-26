version=0.2.0
app_name=cloudbleed-check
release_dir=dist/releases
build:
	CGO_ENABLED=0 go build -ldflags '-w'

build-dists: build-linux build-darwin build-windows

build-linux: build-linux-386 build-linux-amd64

build-darwin: build-darwin-386 build-darwin-amd64

build-windows: build-windows-386 build-windows-amd64

build-linux-386:
	$(call DIST_template,linux,386)

build-linux-amd64:
	$(call DIST_template,linux,amd64)

build-darwin-386:
	$(call DIST_template,darwin,386)

build-darwin-amd64:
	$(call DIST_template,darwin,amd64)

build-windows-386:
	$(call DIST_template,windows,386)

build-windows-amd64:
	$(call DIST_template,windows,amd64)

define DIST_template =
	$(eval dist_os := $(1))
	$(eval dist_arch := $(2))
	$(eval dist_dir := dist/${dist_os}/${app_name}-${version}-${dist_os}-${dist_arch})
	mkdir -p ${dist_dir}
	cp README.md ${dist_dir}
	CGO_ENABLED=0 GOOS=${dist_os} GOARCH=${dist_arch} go build -o ${dist_dir}/${app_name} -ldflags '-w'
	mkdir -p ${release_dir}
	cd ${dist_dir}/.. ; tar -czf ../../${release_dir}/${app_name}-${version}-${dist_os}-${dist_arch}.tar.gz ${app_name}-${version}-${dist_os}-${dist_arch}
endef

docker-build:
	docker build . --tag peterrosell/cloudbleed-check:${version} --tag peterrosell/cloudbleed-check:latest

docker-run:
	cat test-data/matching-test/input-your-sites.txt | docker run -i peterrosell/cloudbleed-check:${version} | rev

docker-push:
	docker push peterrosell/cloudbleed-check:${version}
	docker push peterrosell/cloudbleed-check:latest
