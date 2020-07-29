# service

## configkit

**注意：请勿在init中获取configkit的参数值，那时还未加载。**

**在config.json中即使设置了key=""，则viper中已经set了。但是没有set不代表没有默认值**

config.json中的对象嵌套，对应于字符串的表达是xx.yy。

viper是大小写不敏感的。

viper在cobra使用时，bind用在cmd.Run中，而不是init中