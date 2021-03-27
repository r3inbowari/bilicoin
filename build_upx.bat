@echo off
echo Go build running using UPX

set packageName="cmd"
set appName="bilicoin"
set buildVersion="v1.0.4"

cd %packageName%

set GOOS=windows
set GOARCH=amd64
go build -o ../build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
cd ..
upx ./build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe

cd %packageName%
set GOOS=linux
set GOARCH=amd64
go build -o ../build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
cd ..
upx ./build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%

cd %packageName%
set GOOS=linux
set GOARCH=arm64
go build -o ../build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
cd ..
upx ./build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%

cd %packageName%
set GOOS=darwin
set GOARCH=amd64
go build -o ../build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
cd ..
upx ./build_upx/%appName%_%GOOS%_%GOARCH%_%buildVersion%

pause