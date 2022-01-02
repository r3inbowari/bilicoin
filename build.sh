#!/bin/bash

# author r3inbowari
# date 2021/10/30
# 编译前请确保已安装 git

packageName=cmd
appName=bilicoin
buildVersion=v1.1.0
major=1
minor=1
patch=0
Mode=REL

goVersion=$(go version)
gitHash=$(git show -s --format=%H)
buildTime=$(git show -s --format=%cd)

echo ===================================================
echo "                 Go build running"
echo ===================================================
echo $goVersion
echo build hash $gitHash
echo build time $buildTime
echo build tag $buildVersion
echo ===================================================

if [ ! -d "build_upx" ]; then
  mkdir build_upx
fi

cd cmd
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

GOOS=windows
GOARCH=amd64
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}.exe
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=windows
GOARCH=arm64
go env -w GOOS=windows
go env -w GOARCH=arm64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}.exe
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=linux
GOARCH=amd64
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=linux
GOARCH=arm64
go env -w GOOS=linux
go env -w GOARCH=arm64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=darwin
GOARCH=amd64
go env -w GOOS=darwin
go env -w GOARCH=amd64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=darwin
GOARCH=arm64
go env -w GOOS=darwin
go env -w GOARCH=arm64
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}

GOOS=linux
GOARCH=mipsle
go env -w GOOS=linux
go env -w GOARCH=mipsle
go env -w GOMIPS=softfloat
go build -ldflags "-X 'main.Major=${major}' -X 'main.Minor=${minor}'-X 'main.Patch=${patch}' -X 'main.ReleaseVersion=${buildVersion}' -X 'main.Mode=${Mode}' -X 'main.goVersion=${goVersion}' -X 'main.GitHash=${gitHash}' -X 'main.buildTime=${buildTime}'" -o ../build/${appName}_${GOOS}_${GOARCH}_${buildVersion}
echo Done ${appName}_${GOOS}_${GOARCH}_${buildVersion}
