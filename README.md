# go-sdk-gm 
首先验证了网络上所有的go-gm-sdk，发现99%以上都不能用，或许是我不能理解用法精髓吧（事实上国内99%的开源项目都不能用，内卷严重，鄙视），
终于找到一个能用的sdk，不过那项目关闭了，在此基础上二次开发，支持rest api，功能一一测试，
目前确定能用的功能有：创建通道，加入通道安装实例化调用查询等，随项目上传一个testApi.sh脚本，curl测试都在里面，
可以一键testAPI.sh脚本测试，大部分功能在测试脚本里，少量没写进去，全部测试用例建议查看main文件。

使用方法：

项目位置没有要求，不一定gopath下，任意位置即可。
然后把config.yaml中的证书路径替换你自己实际的，如果不想替换，本工程也提供了一个本地了channle2文件夹，这个下面是各种证书，将config.yaml路径改为本地当前目录即可。
然后go run main.go 启动服务在端口4000，
参照testAPI.sh中的curl命令使用或者运行testAPI.sh

注意：手动运行中，注册login后返回token，把这个token复制下来后续命令中替换token变量。

QQ：314588595
微信：marschengshuangquan

mysql -u root -p

CREATE TABLE balance(`name`  varchar(20) not null,`sum` INT UNSIGNED, PRIMARY KEY (name)) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO balance(name,sum) VALUES ("A",100 );


查询：

curl -s -X POST http://localhost:4000/channels/mychannel/invokechaincodes/example -H  "authorization:eyJhb
GciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU4ODQxMDIsIm9yZ05hbWUiOiJvcmcxIiwidXNlcm5hbWUiOiJqaW0ifQ.LrFBpEb9ogWu2z8ZnWRsSTedIjwDAMc8ENOuD_BdcFo" -H "conte
nt-type: application/json" -d "{
 \"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
 \"fcn\":\"query\",
  \"args\":[\"a\"]
 }"

move交易：

curl -s -X POST http://localhost:4000/channels/mychannel/invokechaincodes/example -H \
"authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDU4ODQxMDIsIm9yZ05hbWUiOiJvcmcxIiwidXNlcm5hbWUiOiJqaW0ifQ.LrFBpEb9ogWu2z8ZnWRsSTedIjwDAMc8ENOuD_BdcFo" -H "content-type: application/json"  \
-d "{  
 \"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
 \"fcn\":\"move\",
  \"args\":[\"a\",\"b\",\"1\"]
 }"

go操作mysql：
https://www.cnblogs.com/kaichenkai/p/11140555.html
