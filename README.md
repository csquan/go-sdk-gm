# go-sdk-gm 
首先验证了网络上所有的go-gm-sdk，发现99%以上都不能用，或许是我不能理解用法精髓吧（事实上国内99%的开源项目都不能用，内卷严重，鄙视），
终于找到一个能用的sdk，不过那项目关闭了，在此基础上二次开发，支持rest api，功能一一测试，
目前确定能用的功能有：创建通道，加入通道安装实例化调用查询等，随项目上传一个testApi.sh脚本，curl测试都在里面，
目前testAPI整体脚本没有测试过，我是一个个curl拿出来测试的，测试通过，计划下一步把整体脚本测试调整下，这样以后测试一键脚本ok。

使用方法：

项目位置没有要求，不一定gopath下，任意位置即可。
然后把config.yaml中的证书路径替换你自己实际的，如果不想替换，本工程也提供了一个本地了channle2文件夹，这个下面是各种证书，将config.yaml路径改为本地当前目录即可。
然后go run main.go 启动服务在端口4000，
参照testAPI.sh中的curl命令使用。

注意：注册login后返回token，把这个token复制下来后续命令中替换token变量。

由于该版本不太完善，我也会持续更新功能，如有使用问题也请反馈。

QQ：314588595
微信：marschengshuangquan
