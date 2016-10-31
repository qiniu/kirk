# 安装 
- 安装 [go 1.7](https://golang.org/dl/) 以上版本
- 安装 [glide](https://glide.sh) 包管理工具

# 开发环境配置
同步项目代码：
```bash
$ git clone https://github.com/qiniu/kirk.git 
```

本项目采用 glide 管理包依赖，你可以使用如下命令下载开发中需要用到的包，这些包会被下载到vendor目录下：
```bash
$ glide install
```

如果改动中引入了新的依赖，请使用如下命令更新 glide 配置文件 [glide.yaml](glide.yaml) 和 [glide.lock](glide.yaml) 并将其包含在PR中
```bash
$ glide up
$ glide install
```

# 提交PR
- 在着手影响较大的代码改动，尤其是 breaking change 前，请先建立issue，与本项目组充分沟通并确定修改方案后再开工
- 提交PR前，请确保代码格式检查以及单元测试通过
```bash
$ make style
$ make test
```
- 如果改动涉及 SDK 的公开接口，请更新代码中相应的注释
- 尽量保证每个commit包含一个独立的改动，并在 [CHANGELOG.md](CHANGELOG.md) 添加改动内容

