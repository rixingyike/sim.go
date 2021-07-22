# sim.go
一个快速开发小程序后台接口的工具类库

## 如何使用

鉴于go语言目前有了mod，在默认启用了mod的情况下，直接使用该源码可能会比较麻烦。所以，建议：

```
git clone https://github.com/rixingyike/sim.go.git --depth=1
export GO111MODULE=off
go get ./...
./debug.sh
```


## History
9/29: 修正qiniu api升级引发的编译错误
