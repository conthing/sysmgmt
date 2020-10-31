.PHONY: build test clean

# make SERVICENAME=xxx xxx为服务的程序名，对应./cmd/xxx目录
# make arm 构建arm程序
# make amd64 构建amd64程序
# make 构建arm程序和amd64程序

USEZEROMQ=0
USECGO=0

VERSION=$(shell sed -n '1p' ./VERSION)
BUILD_TIME:=$(shell date "+%F %T %Z")

WORKDIR=$(shell pwd)
ifeq (,$(SERVICENAME))
SERVICENAME=$(shell basename $(WORKDIR))
endif
HOSTOS=$(shell go env GOOS)
HOSTARCH=$(shell go env GOARCH)

GOFLAGS=-ldflags "-X github.com/conthing/utils/common.Version=$(VERSION) -X 'github.com/conthing/utils/common.BuildTime=${BUILD_TIME}'"

GOARM=GOARCH=arm GOOS=linux go
GOAMD64=GOARCH=amd64 go
ifneq (0,$(USECGO))
GOARM=GOARCH=arm GOOS=linux CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 CGO_LDFLAGS="-g -O2 -L/usr/arm-linux-gnueabihf/lib"  go
GOAMD64=GOARCH=amd64 CGO_ENABLED=1 go
ifneq (0,$(USEZEROMQ))
GOARM=PKG_CONFIG_PATH=/usr/local/zeromq-4.2.2/arm/lib/pkgconfig GOARCH=arm GOOS=linux CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -I/usr/local/zeromq-4.2.2/arm/include" CGO_LDFLAGS="-g -O2 -L/usr/local/zeromq-4.2.2/arm/lib -L/usr/arm-linux-gnueabihf/lib -Wl,-rpath-link /usr/local/zeromq-4.2.2/arm/lib -Wl,-rpath-link /usr/arm-linux-gnueabihf/lib" go
ifeq (windows,$(HOSTOS))
GOAMD64=go
else
GOAMD64=PKG_CONFIG_PATH=/usr/local/zeromq-4.2.2/amd64/lib/pkgconfig GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -I/usr/local/zeromq-4.2.2/amd64/include" CGO_LDFLAGS="-g -O2 -L/usr/local/zeromq-4.2.2/amd64/lib -Wl,-rpath-link /usr/local/zeromq-4.2.2/amd64/lib" go
endif
endif
endif

MICROSERVICESARM=output/arm/$(SERVICENAME)
MICROSERVICESAMD64=output/amd64/$(SERVICENAME)
ifeq (windows,$(HOSTOS))
MICROSERVICESAMD64=output/amd64/$(SERVICENAME).exe
endif

# 每次都重新构建，避免两个环境下遗留的文件影响
.PHONY: $(MICROSERVICESARM) $(MICROSERVICESAMD64)

# 第一个target为默认target，设置为hostarch对应的target
ifeq (arm,$(HOSTARCH))
arm: $(MICROSERVICESARM)
amd64: $(MICROSERVICESAMD64)
else
amd64: $(MICROSERVICESAMD64)
arm: $(MICROSERVICESARM)
endif

all: $(MICROSERVICESARM) $(MICROSERVICESAMD64)

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
