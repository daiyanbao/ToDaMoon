# 后台运行TDM程序

按照[Go语言小贴士4 － 后台运行](https://zhuanlan.zhihu.com/p/21839884?refer=idada) 的内容，可以非常方便地构建一个ToDaMoon的后台运行方式。

我在zsh中添加了相关设置是
```shell
alias tdm="cd ~/ToDaMoon/ && nohup ./ToDaMoon 1>>tdm.out 2>>tdm.err &"
alias kt="cd ~/ToDaMoon && kill \`cat tdm.pid\`"
alias tto="cd ~/ToDaMoon/ && tailf tdm.out"
alias tte="cd ~/ToDaMoon/ && tailf tdm.err"
alias tdmv="~/ToDaMoon/ToDaMoon v"
```