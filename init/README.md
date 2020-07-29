# init

## initkit

func LoadConfig：在main中使用，默认加载程序所在目录下的config.json。

func DefFlags：在cmd的init中使用，默认加载本库支持的所有flags，从而支持程序的命令行参数设置。

func BindFlags：在cmd中使用，和DefFlags配套。