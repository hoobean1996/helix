#!/bin/bash
# 改进的vmail测试脚本

# 1. 检查服务是否正常运行
echo "==== 检查vmail服务状态 ===="
kubectl get pods -l app=vmail
kubectl get svc vmail

# 2. 获取pod名称
VMAIL_POD=$(kubectl get pods -l app=vmail -o jsonpath='{.items[0].metadata.name}')
echo "VMAIL Pod: $VMAIL_POD"

# 3. 查看测试前的日志
echo "==== 测试前的日志 ===="
kubectl logs $VMAIL_POD

# 4. 设置端口转发以便从本地访问
echo "==== 设置端口转发 ===="
kubectl port-forward svc/vmail 10025:10025 &
PORT_FORWARD_PID=$!
echo "端口转发PID: $PORT_FORWARD_PID"

# 等待端口转发建立（更长的等待时间）
echo "等待端口转发建立..."
sleep 5

# 检查端口转发是否成功
if ! nc -z localhost 10025 >/dev/null 2>&1; then
  echo "错误: 端口转发未成功建立，无法连接到10025端口"
  kill $PORT_FORWARD_PID 2>/dev/null
  exit 1
else
  echo "端口转发成功建立，可以连接到10025端口"
fi

# 5. 使用telnet/nc测试SMTP连接，更详细的SMTP命令
echo "==== 使用nc测试SMTP连接 ===="
(
  sleep 1
  echo "EHLO example.com"
  sleep 1
  echo "MAIL FROM:<test@example.com>"
  sleep 1
  echo "RCPT TO:<recipient@example.com>"
  sleep 1
  echo "DATA"
  sleep 1
  echo "Subject: Test Email"
  echo ""
  echo "This is a test email body."
  echo "."
  sleep 1
  echo "QUIT"
) | nc -v localhost 10025

# 等待处理完成
echo "等待服务处理请求..."
sleep 3

# 6. 查看测试后的日志（查看测试期间产生的新日志）
echo "==== 测试后的日志 ===="
kubectl logs $VMAIL_POD --since=1m

# 7. 结束端口转发
echo "==== 关闭端口转发 ===="
kill $PORT_FORWARD_PID
wait $PORT_FORWARD_PID 2>/dev/null
echo "端口转发已关闭"

# 8. 进入容器查看日志（如果容器内有日志文件）
echo "==== 容器内日志检查 ===="
echo "如果需要，请执行: kubectl exec -it $VMAIL_POD -- sh"
echo "然后在容器内查看日志，例如: cat /var/log/vmail.log"