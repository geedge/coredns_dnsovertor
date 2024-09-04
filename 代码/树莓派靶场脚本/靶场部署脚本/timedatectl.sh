#! /bin/bash
#Tor靶场需要各节点时间同步，最好使用UTC时区。
#由于匿名网关使用openwrt系统，可能需要再次手动设置时区。
source config.sh

success=0
# 循环遍历 IP 列表，设置timedatectl
for IP in "${IP_LIST3[@]}"; do
  echo "在远程IP上进行timedatectl: $IP"
  # 运行远程命令，使用标准时间UTC时区。
  #sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "timedatectl set-timezone Asia/Shanghai"
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "timedatectl set-timezone UTC"
 
  # 检查传输是否成功
  if [ $? -eq 0 ]; then
    let success=success+1
    echo "成功"
  else
    echo "失败"
  fi
  
  echo ""
done

echo "成功调整$success台机器的时钟,期望是31台"
