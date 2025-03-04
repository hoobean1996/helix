一定需要配置relay规则，否则的话postfix转发失败
POSTFIX_POD=$(kubectl get pod -l app=postfix -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POSTFIX_POD -- bash
# 在容器内
# 重置虚拟别名配置
postconf -e "virtual_alias_maps ="
postconf -e "virtual_alias_domains ="

# 设置中继域和中继主机
postconf -e "relay_domains = vmail.today"
postconf -e "relayhost = vmail.default.svc.cluster.local:25"

# 重新加载配置
postfix reload