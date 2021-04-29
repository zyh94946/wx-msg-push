VERSION=`cat VERSION`
COMMIT=`git rev-parse --short HEAD`
CFGPATH='./config.toml'
BUILDDATE=`date "+%Y-%m-%d"`

BUILD_DIR=.
APP_NAME=wx-msg-push

build = CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)$(3) -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go  
upx = upx ${BUILD_DIR}/$(APP_NAME)-$(1)-$(2)$(3)

LINUX = linux-amd64

WINDOWS = windows-amd64-.exe

DARWIN = darwin-amd64

ALL = $(LINUX) \
      $(WINDOWS) \
	  $(DARWIN)

build_linux: $(LINUX:%=build/%)

build_windows: $(WINDOWS:%=build/%)

build_darwin: $(DARWIN:%=build/%)

build_all: $(ALL:%=build/%)

build/%:
	$(call build,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))
	$(call upx,$(firstword $(subst -, , $*)),$(word 2, $(subst -, ,$*)),$(word 3, $(subst -, ,$*)))

build:
	go build -o wx-msg-push -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go

runs:
	go run -ldflags "-X main.v=${VERSION} -X main.c=${COMMIT} -X main.d=${BUILDDATE}" main.go server -c ${CFGPATH}
