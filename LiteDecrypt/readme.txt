跨平台编译：
首先找到是国内文章是利用下载Go的源码包，将对其进行编译安装，完成后Go/src下的make.bash生成跨平台的编译器，这样在每次用时需要指定ＧＯＯＳ及ＧＯＡＲＣＨ来进行编译，如：
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build　／／编译为linux 64位系统下的程序
如果需要编译windows 64 位下的程序则要再写一个命令：
CGO_ENABLED=0 GOOS=windows GOARCH=amd64(或386) go build　／／编译为windows 64位系统下的程序

图形界面：
通过这个main.manifest执行（如果walk安装好了，应该就有这个命令）：
rsrc -manifest main.manifest -o rsrc.syso

上图是rsrc命令的参数列表，比如我为自已的程序加个ico图标：
rsrc -manifest main.manifest –ico icon.ico -o rsrc.syso

不显示dos:
go build -ldflags "-H windowsgui"