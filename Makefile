.PHONY: build test clean

# SERVICENAME为服务的程序名
# make arm 构建arm程序
# make amd64 构建amd64程序
# make 构建arm程序和amd64程序

SERVICENAME=sysmgmt
USEZEROMQ=0

VERSION=$(shell cat ./VERSION)
BUILD_TIME:=$(shell date "+%F %T")

GOFLAGS=-ldflags "-X github.com/conthing/utils/common.Version=$(VERSION)\
 -X 'github.com/conthing/utils/common.BuildTime=${BUILD_TIME}'"

GOARM=GOARCH=arm GOOS=linux go
GOAMD64=GOARCH=amd64 go
ifneq (0,$(USEZEROMQ))
GOARM=PKG_CONFIG_PATH=/usr/local/zeromq-4.2.2/arm/lib/pkgconfig GOARCH=arm CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -I/usr/local/zeromq-4.2.2/arm/include" CGO_LDFLAGS="-g -O2 -L/usr/local/zeromq-4.2.2/arm/lib -L/usr/arm-linux-gnueabihf/lib -Wl,-rpath-link /usr/local/zeromq-4.2.2/arm/lib -Wl,-rpath-link /usr/arm-linux-gnueabihf/lib" go
GOAMD64=PKG_CONFIG_PATH=/usr/local/zeromq-4.2.2/amd64/lib/pkgconfig GOOS=linux GOARCH=amd64 go
endif

MICROSERVICESARM=output/arm/$(SERVICENAME)
MICROSERVICESAMD64=output/amd64/$(SERVICENAME)
ifeq (windows,$(shell go env GOOS))
MICROSERVICESAMD64=output/amd64/$(SERVICENAME).exe
endif

# 每次都重新构建，避免两个环境下遗留的文件影响
.PHONY: $(MICROSERVICESARM) $(MICROSERVICESAMD64)

all: $(MICROSERVICESARM) $(MICROSERVICESAMD64)
arm: $(MICROSERVICESARM)
amd64: $(MICROSERVICESAMD64)

$(MICROSERVICESARM):
	@echo build $(MICROSERVICESARM)...
	$(GOARM) build $(GOFLAGS) -o $@ ./cmd/$(SERVICENAME)
	cp -rf ./cmd/$(SERVICENAME)/. ./output/arm
	rm -f ./output/arm/*.go
$(MICROSERVICESAMD64):
	@echo build $(MICROSERVICESAMD64)...
	$(GOAMD64) build $(GOFLAGS) -o $@ ./cmd/$(SERVICENAME)
	cp -rf ./cmd/$(SERVICENAME)/. ./output/amd64
	rm -f ./output/amd64/*.go

test_arm:
	$(GOARM) test ./... -cover
test_amd64:
	$(GOAMD64) test ./... -cover

clean:
	rm -rf output
