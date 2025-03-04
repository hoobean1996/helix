#!/bin/bash
# e2e.sh - 测试多个随机临时邮箱的脚本
# 用法: ./e2e.sh [邮箱数量] [主机] [端口]

# 设置默认参数
NUM_EMAILS=${1:-100}  # 默认生成100个邮箱
HOST=${2:-localhost}  # 默认主机
PORT=${3:-2525}       # 默认端口
DOMAIN="vmail.today"  # 邮箱域名

# 创建一个日志文件
LOG_FILE="smtp_test_$(date +%Y%m%d_%H%M%S).log"
echo "开始测试: $(date)" > $LOG_FILE
echo "=============================" >> $LOG_FILE
echo "将测试 $NUM_EMAILS 个随机邮箱地址" >> $LOG_FILE
echo "服务器: $HOST:$PORT" >> $LOG_FILE
echo "=============================" >> $LOG_FILE

# 随机字符串生成函数
random_string() {
  local length=$1
  cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w $length | head -n 1
}

# 计数器
successful=0
failed=0

# 测试每个邮箱
for ((i=1; i<=$NUM_EMAILS; i++)); do
  # 生成随机邮箱
  # 随机选择一种邮箱格式
  format=$((RANDOM % 5))
  case $format in
    0) email=$(random_string 10)@$DOMAIN ;;                        # 简单随机字符串
    1) email=$(random_string 5).$(random_string 5)@$DOMAIN ;;      # 带点的随机字符串
    2) email="user_$(random_string 8)@$DOMAIN" ;;                  # 带前缀的随机字符串
    3) email="test.$(date +%s).$(random_string 4)@$DOMAIN" ;;      # 带时间戳的随机字符串
    4) email="temp_$(printf '%03d' $i)_$(random_string 5)@$DOMAIN" ;; # 带序号的随机字符串
  esac
  
  subject="Test Email #$i - $(date +%s)"
  body="This is test email #$i sent at $(date).\r\nTesting random email: $email"
  
  echo -e "\n[$i/$NUM_EMAILS] 测试邮箱: $email" | tee -a $LOG_FILE
  
  # 生成SMTP命令
  commands="EHLO test.local\r\nMAIL FROM: sender@example.com\r\nRCPT TO: $email\r\nDATA\r\nSubject: $subject\r\nFrom: sender@example.com\r\nTo: $email\r\n\r\n$body\r\n.\r\nQUIT\r\n"
  
  # 使用nc发送邮件并保存结果
  result=$(echo -e "$commands" | nc -w 5 $HOST $PORT 2>&1)
  
  # 记录结果
  echo "$result" >> $LOG_FILE
  
  # 检查结果是否包含成功消息
  if echo "$result" | grep -q "250 2.0.0 Ok: queued"; then
    echo -e "  \033[32m成功\033[0m: 邮件已加入队列" | tee -a $LOG_FILE
    successful=$((successful + 1))
  else
    echo -e "  \033[31m失败\033[0m: 发送失败" | tee -a $LOG_FILE
    failed=$((failed + 1))
  fi
  
  # 短暂暂停，避免过快发送
  sleep 0.5
done

# 总结结果
echo -e "\n=============================" | tee -a $LOG_FILE
echo "测试完成: $(date)" | tee -a $LOG_FILE
echo "总共测试: $NUM_EMAILS 个邮箱" | tee -a $LOG_FILE
echo -e "\033[32m成功\033[0m: $successful" | tee -a $LOG_FILE
echo -e "\033[31m失败\033[0m: $failed" | tee -a $LOG_FILE
echo "详细日志已保存到: $LOG_FILE" | tee -a $LOG_FILE
echo "=============================" | tee -a $LOG_FILE

# 如果有任何失败，返回非零退出码
if [ $failed -gt 0 ]; then
  exit 1
fi

exit 0