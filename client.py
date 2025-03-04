#!/usr/bin/env python3
# 使用Python发送测试邮件到vmail服务

import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
import sys
import argparse

def send_test_email(host, port, sender, recipient, subject, body):
    """发送测试邮件到指定的SMTP服务器"""
    
    # 创建邮件
    msg = MIMEMultipart()
    msg['From'] = sender
    msg['To'] = recipient
    msg['Subject'] = subject
    
    # 添加邮件正文
    msg.attach(MIMEText(body, 'plain'))
    
    try:
        # 连接到SMTP服务器
        print(f"正在连接到 {host}:{port}...")
        server = smtplib.SMTP(host, port, timeout=10)
        server.set_debuglevel(1)  # 启用调试输出
        
        # 发送邮件
        print(f"正在发送邮件: {sender} -> {recipient}")
        server.sendmail(sender, recipient, msg.as_string())
        print("邮件发送成功!")
        
        # 关闭连接
        server.quit()
        return True
        
    except Exception as e:
        print(f"发送邮件时出错: {e}")
        return False

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='发送测试邮件到vmail服务')
    parser.add_argument('--host', default='localhost', help='SMTP服务器主机名')
    parser.add_argument('--port', type=int, default=10025, help='SMTP服务器端口')
    parser.add_argument('--from', dest='sender', default='sender@example.com', help='发件人地址')
    parser.add_argument('--to', dest='recipient', default='recipient@example.com', help='收件人地址')
    parser.add_argument('--subject', default='测试邮件', help='邮件主题')
    parser.add_argument('--body', default='这是一封测试邮件，用于验证vmail服务是否正常工作。', help='邮件正文')
    
    args = parser.parse_args()
    
    success = send_test_email(
        args.host, 
        args.port, 
        args.sender, 
        args.recipient, 
        args.subject, 
        args.body
    )
    
    sys.exit(0 if success else 1)
