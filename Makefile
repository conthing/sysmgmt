.PHONY: build test clean

# make 构建目标ARCH即是宿主机ARCH的程序，当前目录名即为程序名，对应./cmd/xxx目录
# make APP=xxx xxx为程序名，对应./cmd/xxx目录
# make ARCH=arm 构建arm程序，对应OS为linux
# make ARCH=amd64 构建amd64程序，对应OS为宿主机OS

USECGO=1
USEZMQ=0

#EXTRAFILE_amd64windows=driver/lib/amd64windows/WTY.dll driver/lib/amd64windows/WT_H264.dll driver/lib/amd64windows/TH_PLATEID.dll
#EXTRAFILE_amd64linux=driver/lib/amd64linux/libwty.so driver/lib/amd64linux/libthplateid.so
#EXTRAFILE_armlinux=driver/lib/armlinux/libwty.so driver/lib/armlinux/libthplateid.so

# 读取版本值
VERSION=$(shell sed -n '1p' ./VERSION)

# 获取构建时间
BUILD_TIME:=$(shell date "+%F %T %Z")

# 工作目录
WORKDIR=$(shell pwd)

# 默认程序名称 = 工作目录名称
ifeq (,$(APP))
APP=$(shell basename $(WORKDIR))
endif

# 宿主机ARCH和OS
HOSTARCH=$(shell go env GOARCH)
HOSTOS=$(shell go env GOOS)

# 目标ARCH
ifeq (,$(ARCH))
ARCH=$(HOSTARCH)
endif

# 目标OS，如果目标ARCH为arm，目标OS一定是linux
ifeq (arm,$(ARCH))
OS=linux
else
OS=$(HOSTOS)
endif

# 注入版本号和构建时间
GOFLAGS=-ldflags "-X github.com/conthing/utils/common.Version=$(VERSION) -X 'github.com/conthing/utils/common.BuildTime=${BUILD_TIME}'"

# GO命令环境变量设置
GOPREFIX:=

ifneq ($(ARCH),$(HOSTARCH))
GOPREFIX += GOARCH=$(ARCH) 
endif

ifneq ($(OS),$(HOSTOS))
GOPREFIX += GOOS=$(OS) 
endif

ifneq (0,$(USECGO))
GOPREFIX += CGO_ENABLED=1 
## 如果目标ARCH为arm，宿主机ARCH又不是arm，指定交叉编译器
ifeq (arm,$(ARCH))
ifneq (arm,$(HOSTARCH))
ifeq (windows,$(HOSTOS))
GOPREFIX += CC=arm-linux-gnueabihf-gcc CGO_LDFLAGS="-g -O2 -LC:\msys64\opt\gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf\arm-linux-gnueabihf\lib -Wl,-rpath-link C:\msys64\opt\gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf\arm-linux-gnueabihf\lib" 
else
ifeq (1,$(USEZMQ))
GOPREFIX += PKG_CONFIG_PATH=/usr/local/zeromq-4.2.2/arm/lib/pkgconfig CC=arm-linux-gnueabihf-gcc CGO_CFLAGS="-g -O2 -I/usr/local/zeromq-4.2.2/arm/include" CGO_LDFLAGS="-g -O2 -L/usr/local/zeromq-4.2.2/arm/lib -L/usr/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/arm-linux-gnueabihf/lib -Wl,-rpath-link /usr/local/zeromq-4.2.2/arm/lib -Wl,-rpath-link /usr/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/arm-linux-gnueabihf/lib" 
else
GOPREFIX += CC=arm-linux-gnueabihf-gcc CGO_LDFLAGS="-g -O2 -L/usr/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/arm-linux-gnueabihf/lib -Wl,-rpath-link /usr/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/arm-linux-gnueabihf/lib" 
endif
endif
endif
endif
endif

# 目标文件
ifeq (windows,$(OS))
OUTPUT=output/$(ARCH)$(OS)/$(APP).exe
else
OUTPUT=output/$(ARCH)$(OS)/$(APP)
endif



# 每次都重新构建，避免两个环境下遗留的文件影响
.PHONY: $(OUTPUT)

default: $(OUTPUT) postscript

$(OUTPUT):
	@echo build $(OUTPUT)...
	$(GOPREFIX) go build $(GOFLAGS) -o $@ ./cmd/$(APP)
	cp -rf ./cmd/$(APP)/. ./output/$(ARCH)$(OS)
	rm -f ./output/$(ARCH)$(OS)/*.go

ifeq (,$(EXTRAFILE_$(ARCH)$(OS)))
postscript:
else
postscript:
	cp $(EXTRAFILE_$(ARCH)$(OS)) ./output/$(ARCH)$(OS)
endif

test:
	$(GOARM) test ./... -cover

all: $(OUTPUT) test

clean:
	rm -rf output
