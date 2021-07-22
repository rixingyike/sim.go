# sim.go
一个快速开发小程序后台接口的工具类库

## 如何使用（带mod）

建议使用最新的mod特性，方式如下：

```
git clone -b v2.0 https://github.com/rixingyike/sim.go.git --depth=1
go env -w GOPROXY=https://goproxy.cn,https://gocenter.io,https://goproxy.io,direct
export GO111MODULE=on
cd ./sim.go
go mod download
./debug.sh
```

2021年7月修改

## 如何使用（旧，不带mod）

鉴于go语言目前有了mod，在默认启用了mod的情况下，直接使用该源码可能会比较麻烦。所以建议先将mod关闭：

```
git clone -b v1.0 https://github.com/rixingyike/sim.go.git --depth=1
export GO111MODULE=off
cd ./sim.go
go get ./...
./debug.sh
```

# History
9/29: 修正qiniu api升级引发的编译错误
