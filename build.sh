#!/bin/bash
# 编译linux_x86_64
GOOS=linux GOARCH=amd64 go build -o sing-box-sub-helper
# 压缩
tar -czvf sing-box-sub-helper_linux_amd64.tar.gz sing-box-sub-helper
rm -rf sing-box-sub-helper
# 编译linux_arm64
GOOS=linux GOARCH=arm64 go build -o sing-box-sub-helper
# 压缩
tar -czvf sing-box-sub-helper_linux_arm64.tar.gz sing-box-sub-helper
rm -rf sing-box-sub-helper
# 编译linux_arm7
GOOS=linux GOARCH=arm GOARM=7 go build -o sing-box-sub-helper
# 压缩
tar -czvf sing-box-sub-helper_linux_arm7.tar.gz sing-box-sub-helper
rm -rf sing-box-sub-helper
# 编译mac_amd64
GOOS=darwin GOARCH=amd64 go build -o sing-box-sub-helper
# 压缩
tar -czvf sing-box-sub-helper_mac_amd64.tar.gz sing-box-sub-helper
rm -rf sing-box-sub-helper
# 编译mac_arm64
GOOS=darwin GOARCH=arm64 go build -o sing-box-sub-helper
# 压缩
tar -czvf sing-box-sub-helper_mac_arm64.tar.gz sing-box-sub-helper
rm -rf sing-box-sub-helper
# 编译windows_amd64
GOOS=windows GOARCH=amd64 go build -o sing-box-sub-helper.exe
# 压缩
tar -czvf sing-box-sub-helper_windows_amd64.tar.gz sing-box-sub-helper.exe
rm -rf sing-box-sub-helper.exe
