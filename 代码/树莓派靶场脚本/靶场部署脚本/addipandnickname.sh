#! /bin/bash
#为router_torrc增加必要的配置。 Nickname和Address
source config.sh

length=3
# 循环遍历 IP 列表并执行命令
for IP in "${IP_LIST3[@]}"; do
  echo "在远程IP上进行router_torrc传输和tor启动//tor: $IP"
  # 运行远程命令
  substring=${IP: -$length}
  str1="Router"
  str1+="$substring"
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "echo 'Nickname $str1'>>/home/pi/router_torrc"
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "echo 'Address $IP'>>/home/pi/router_torrc"
 
  # 检查执行是否成功
  if [ $? -eq 0 ]; then
    echo "执行成功"
  else
    echo "执行失败"
  fi
  
  echo ""
done
