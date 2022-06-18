# goproxy

我们想把MySQL代理到本地，在不用SSH隧道的情况下，我们尽量想使用K8S Port Forward的能力。

但发现K8S无法很好的将 External Service Forward到本地，如下所示：

https://stackoverflow.com/questions/64429094/kubernetes-cant-port-forward-externalname-service


所以写了一个gopxoy

通过以下配置
```toml
[[proxy]]
srcAddr = "172.17.245.230:13306"
dstPort = 9999
```

我们可以将MySQL配置暴露到某个POD的一个端口

然后我们在使用
```bash
kubectl port-forward pod/{你的pod名称}  19999:9999  
```

访问本地 :19999， 那么你就可以访问到远端的mysql