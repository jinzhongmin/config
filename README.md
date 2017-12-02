# config
用于解析网络设备的配置文件，如华为的交换机路由器等，初衷是为了解析华为GPON的配置文件编写的，并未测试其他网络设备的配置文件解析情况
## 说明
func clean(text string) string 为了整理GPON配置文件从crt中复制出来随意断行的问题，解析其他配置文件请酌情修改
## 使用
有Config、View、Line三个结构体

Config表示全部配置文件

View代表某个视图，如某个端口配置

Line代表某行命令