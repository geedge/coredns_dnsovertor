一、插件说明：
插件代码见./ag_201_coredns/plugin/dnsovertor
代码暂未使用git进行版本管理，上传的是初版代码，能实现基本功能，metrics和setup_test功能未经测试。
本插件使Coredns通过Tor的socks5代理端口向xxxx.onion域名的暗网隐藏服务发送DOH请求。
本插件应与Coredns原代码配合使用

插件功能：向暗网DNS服务转发DNS请求，向客户转发DNS响应；处理请求和响应时输出日志
注意事项：暗网DNS地址xxxx.onion提供DNS隐藏服务。

二、使用方法：
1.在plugin.cfg末尾添加dnsovertor:dnsovertor。
2.将dnsovertor文件夹放入Coredns源文件的plugin文件夹。
以上两步已集成到ag_201_coredns文件夹中。
3.通过源代码编译Coredns。


三、代码说明：
dnsovertor.go和setup.go是dnsovertor插件的两个核心go文件，这两个文件就可以实现基本功能。
dohclient.py,dohclient.go使用DoH查询域名alibaba.com.
（https://help.aliyun.com/document_detail/171664.html）

四、匿名解析的测试例子（VPS已失效，流程类似）：
匿名解析网关：8.210.118.152，在443端口提供了DOH服务（域名证书为www.runn77.tk），在53端口提供DNS服务
对www.runn77.tk发起doh请求：python dohclient.py
对8.210.118.152发起dns请求：dig @8.210.118.152 baidu.com 或 nslookup baidu.com 8.210.118.152



