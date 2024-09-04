#!/bin/bash

# 权威目录服务器启动除了匿名网关外的其他tor进程。
# 使用配置文件 config.sh

# 加载配置文件
source config.sh

# 初始化变量
i=1
success=0
failure=0

# 定义日志文件
log_file="tor_startup.log"

# 记录启动时间
echo "启动时间: $(date)" > $log_file

# 循环遍历 IP 列表并执行远程命令
for IP in "${IP_LIST2[@]}"; do
  echo "在远程IP上tor启动: $IP" | tee -a $log_file
  echo "命令: sshpass -p \"$REMOTE_PASSWORD\" ssh -t \"$REMOTE_USER@$IP\" \"/home/pi/tor -f /home/pi/router_torrc\"" | tee -a $log_file
  
  # 执行远程命令
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "/home/pi/tor -f /home/pi/router_torrc"
  
  # 检查传输是否成功
  if [ $? -eq 0 ]; then
    let success=success+1
    echo "启动tor成功" | tee -a $log_file
  else
    let failure=failure+1
    echo "启动tor失败" | tee -a $log_file
  fi
  
  echo "" | tee -a $log_file
  
  # 增加延迟，避免过快连接
  sleep 2
done

# 总结结果
echo "启动总结" | tee -a $log_file
echo "成功$success台router节点,期望情况29台" | tee -a $log_file
echo "失败$failure台router节点,期望情况0台" | tee -a $log_file