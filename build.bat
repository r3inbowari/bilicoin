@echo off

@REM author r3inbowari
@REM date 2021/10/1
@REM 编译前请确保已安装 git

set packageName=cmd
set appName=bilicoin
set buildVersion=v1.0.10
set major=1
set minor=0
set patch=10
set Mode=REL

for /f "delims=" %%i in ('go version') do (set goVersion=%%i)
for /f "delims=" %%i in ('git show -s --format^=%%H') do (set gitHash=%%i)
for /f "delims=" %%i in ('git show -s --format^=%%cd') do (set buildTime=%%i)

echo ===================================================
echo                  Go build running
echo ===================================================
echo %goVersion%
echo build hash %gitHash%
echo build time %buildTime%
echo build tag %buildVersion%
echo ===================================================


if not exist build_upx (
    md build_upx
)

cd %packageName%

set GOOS=windows
set GOARCH=amd64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
@REM go build -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%.exe

set GOOS=windows
set GOARCH=arm64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
@REM go build -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%.exe

set GOOS=linux
set GOARCH=amd64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion% ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%


set GOOS=linux
set GOARCH=arm64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion% ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion% ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion% ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

set GOOS=linux
set GOARCH=mipsle
set GOMIPS=softfloat
go build -ldflags "-X 'main.Major=%major%' -X 'main.Minor=%minor%'-X 'main.Patch=%patch%' -X 'main.ReleaseVersion=%buildVersion%' -X 'main.Mode=%Mode%' -X 'main.goVersion=%goVersion%' -X 'main.GitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
copy ..\build\%appName%_%GOOS%_%GOARCH%_%buildVersion% ..\build_upx\
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%
echo ===================================================

cd ../build
certutil -hashfile bilicoin_windows_amd64_%buildVersion%.exe MD5
certutil -hashfile bilicoin_windows_arm64_%buildVersion%.exe MD5
certutil -hashfile bilicoin_linux_amd64_%buildVersion% MD5
certutil -hashfile bilicoin_linux_arm64_%buildVersion% MD5
certutil -hashfile bilicoin_darwin_amd64_%buildVersion% MD5
certutil -hashfile bilicoin_darwin_arm64_%buildVersion% MD5
certutil -hashfile bilicoin_linux_mipsle_%buildVersion% MD5
echo ===================================================

@REM echo %upxArgs%

cd ..\\build_upx
..\upx.exe %upxArgs%

pause