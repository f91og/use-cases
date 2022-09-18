# Step
## 1. 查看APIServer是否开启了MutatingAdmissionWebhook和ValidatingAdmissionWebhook

### 获取apiserver pod名字
apiserver_pod_name=`kubectl get --no-headers=true po -n kube-system | grep kube-apiserver | awk '{ print $1 }'`
### 查看api server的启动参数plugin
`kubectl get po $apiserver_pod_name -n kube-system -o yaml | grep plugin`
如果输出如下，说明已经开启

- --enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook
否则，需要修改启动参数，请不然直接修改Pod的参数，这样修改不会成功，请修改配置文件
/etc/kubernetes/manifests/kube-apiserver.yaml，加上相应的插件参数后保存，APIServer的Pod会监控该文件的变化，然后重新启动

## 2. 实现webhook的逻辑代码，build出image
## 3. 把image deploy到k8s集群里
# Links
- [从0到1开发K8S_Webhook最佳实践 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/404764407)
- [Helm | Template Function List](https://helm.sh/docs/chart_template_guide/function_list/)
- [charts/apiregistration.yaml at master · helm/charts (github.com)](https://github.com/helm/charts/blob/master/stable/stash/templates/apiregistration.yaml)
- [k8tz/charts at master · k8tz/k8tz (github.com)](https://github.com/k8tz/k8tz/tree/master/charts/k8tz)
- [管理集群中的 TLS 认证 | Kubernetes](https://kubernetes.io/zh-cn/docs/tasks/tls/managing-tls-in-a-cluster/)
- [Helm Hooks 的使用-阳明的博客|Kubernetes|Istio|Prometheus|Python|Golang|云原生 (qikqiak.com)](https://www.qikqiak.com/post/helm-hooks-usage/)
- [Move Your Certs to Helm. Using Helm templates to generate… | by Hagai Barel | Nuvo Tech | Medium](https://medium.com/nuvo-group-tech/move-your-certs-to-helm-4f5f61338aca)
- [Helm | chart开发提示和技巧](https://helm.sh/zh/docs/howto/charts_tips_and_tricks/)
- [Manage Kubernetes Admission Webhook's certificates with cert-manager CA Injector and Vault PKI](https://medium.com/trendyol-tech/manage-kubernetes-admission-webhooks-certificates-with-cert-manager-ca-injector-and-vault-pki-281b065e1044)
- [Creating X.509 TLS certificate in Kubernetes](https://www.digihunch.com/2021/08/creating-tls-certificate-kubernetes/) 👈 在k8s中生成证书的3种方式