# Network-detector Client
这个工具用来给Local快速排查网络问题，其中提供了多个检查内容，包括：

|检查内容|说明|
|----|----|
|DNS|尝试解析DNS域名，获取对应的IP，并使用获取到的IP来进行后续流程|
|Ping|使用Ping去尝试跟对应IP的网络连通性|
|MTR|traceroute，查看网络的转发，还有网络质量|
|TCP Connection|尝试进行TCP连接，有时候防火墙就对此进行拦截|

然后对于用户来说，只需要拿到我们编译的文件，并在本地执行就好。
注意的是这个程序需要Root权限，如果没有Root权限的话那很多检查无法正常运行。
对于不同操作系统来说，构建方式和用户执行方式也不一样，如下表所示。

|操作系统|开发构建方式|用户执行方式|
|----|----|----|
|Windows|项目下执行make build-windows|右键以管理员身份运行|
|MacOS|项目下执行make build-mac|shell中执行chmod u+x network-detector && sudo ./network-detecto|
|Linux|项目下执行make build-linux|shell中执行chmod u+x network-detecto && sudo ./network-detecto|

