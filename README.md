# config
用于解析网络设备的配置文件，如华为的交换机路由器等，初衷是为了解析华为GPON的配置文件编写的，并未测试其他网络设备的配置文件解析情况
## 说明
func clean(text string) string 为了整理GPON配置文件从crt中复制出来随意断行的问题，解析其他配置文件请酌情修改
## 使用
有Config、View、Line三个结构体

Config表示全部配置文件

View代表某个视图，如某个端口配置

Line代表某行命令

```golang
func (config *Config) Lines(match []string)
//match：某行命令开始的几个字符串

config.Lines([]string{"ont","add"})
//匹配ont add 开头的命令
```

```golang
func (config *Config) ViewsAuto(match []string) []*View
//match：某个视图第一行命令开始的几个字符串
//会匹配从这行开始缩进个数大于第一行的全部命令行，直到遇到下一个缩减个数和第一行一样的行
```
下面这段可以获得两个View
```text
ont-srvprofile gpon profile-id 3 profile-name "HG810"
  ont-port eth 1 
  multicast-forward untag
  port priority-policy eth 1 copy-cos 
  port vlan eth 1 translation 4 user-vlan 4
  commit
 ont-srvprofile gpon profile-id 4 profile-name "HG8010"
  ont-port eth 1 
  multicast-forward untag
  port priority-policy eth 1 copy-cos 
  port vlan eth 1 translation 4 user-vlan 4
  commit
```
