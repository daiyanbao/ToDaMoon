# ToDaMoon
ToDaMoon是我的自己的虚拟币交易系统。

## 使用方式
1. ``git clone https://github.com/aQuaYi/ToDaMoon.git``到你的``$GOPATH``目录。
1. 按照注释，修改`Makefile`文件的第2行的`BINARY`和第4行的`FILE_PATH`设置。
1. 在命令行进入``ToDaMoon/main``目录后，使用``make``命令生成运行程序`ToDaMoon`。
1. 在命令行进入`FILE_PATH`目录，使用`./ToDaMoon`运行程序。


## Tips
1. [Go语言小贴士4 － 后台运行](https://zhuanlan.zhihu.com/p/21839884?refer=idada) 使用这个里面的内容，可以非常方便地构建一个ToDaMoon的后台运行方式。
我在zsh中的相关设置是
```shell
alias tdm="cd ~/ToDaMoon/ && nohup ./ToDaMoon 1>>tdm.out 2>>tdm.err &"
alias killtdm="cd ~/ToDaMoon && kill \`cat tdm.pid\`"
alias tto="cd ~/ToDaMoon/ && tailf tdm.out"
alias tte="cd ~/ToDaMoon/ && tailf tdm.err"
alias tdmv="~/ToDaMoon/ToDaMoon v"
```