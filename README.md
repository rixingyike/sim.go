# sim.go
一个快速开发小程序后台接口的工具类库


## 如何使用（旧，不带mod）

注意：建议使用相同、带mod的v2.0版本

鉴于go语言目前有了mod，在默认启用了mod的情况下，直接使用该源码可能会比较麻烦。所以建议先将mod关闭：

```
git clone -b v1.0 https://github.com/rixingyike/sim.go.git --depth=1
export GO111MODULE=off
cd ./sim.go
go get ./...
./debug.sh
```


## History
9/29: 修正qiniu api升级引发的编译错误
