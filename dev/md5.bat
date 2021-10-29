cd ../build
set version=v1.0.7

certutil -hashfile bilicoin_windows_amd64_%version%.exe MD5
certutil -hashfile bilicoin_windows_arm64_%version%.exe MD5
certutil -hashfile bilicoin_linux_amd64_%version% MD5
certutil -hashfile bilicoin_linux_arm64_%version% MD5
certutil -hashfile bilicoin_darwin_amd64_%version% MD5
certutil -hashfile bilicoin_darwin_arm64_%version% MD5
certutil -hashfile bilicoin_linux_mipsle_%version% MD5