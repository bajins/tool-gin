# tool-gin

> 基于go-gin框架

**此分支使用了`chromedp`**

~~**必须安装`Chrome`**~~

> [CentOS安装Chrome](https://www.bajins.com/System/CentOS.html#chrome)


**当前分支使用了go:embed内嵌资源文件实现打包到一个二进制中，旧的压缩打包方式请访问分支：[zi-pack](https://github.com/woytu/tool-gin/tree/zip-pack)**



## 功能

- 生成激活key
  - [mobaXtermGenerater.js](https://github.com/inused/pages/blob/master/file/tool/js/mobaXtermGenerater.js)
- 获取`xshell`、`xftp`、`xmanager`下载链接
- 格式化NGNIX配置
- 获取Navicat下载地址


## 使用

### 到[releases](https://github.com/woytu/tool-gin/releases)下载解压并运行

```bash
# Windows
# 双击tool-gin-windows.exe根据默认端口8000运行
# 或者在cmd、power shell中
tool-gin-windows.exe


# Linux
nohup ./tool-gin_linux_amd64 -p 5000 >/dev/null 2>&1 &
```
