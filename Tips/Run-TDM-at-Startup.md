# 开机启动ToDaMoon

>假设
>存放ToDaMoon执行文件的文件夹的绝对路径是/path/to/ToDaMoon
>存放启动脚本的文件夹的绝对路径是/path/to/script

## 1.编写启动脚本launchToDaMoon.sh
```shell
#!/bin/bash
#用来启动ToDaMoon程序，让Ta读取到正确的配置

#切换到ToDaMoon所在的目录
cd /path/to/ToDaMoon

#运行ToDaMoon程序
./ToDaMoon 1>>tdm.out 2>>tdm.err

exit 0
```

## 2.把launchToDaMoon.sh存放在/path/to/script目录中。

## 3.在/etc/rc.local的exit 0前填写以下内容
```shell
#统一设置脚本存放的目录为环境变量
SCRIPT_DIR=/path/to/script
export SCRIPT_DIR

#运行ToDaMoon的启动脚本
sh $SCRIPT_DIR/launchToDaMoon.sh
```
