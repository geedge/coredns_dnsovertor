#! /bin/bash
# 传输必要的文件。
source config.sh

success=0
# 循环遍历 IP 列表并进行文件传输
for IP in "${IP_LIST3[@]}"; do

  echo "ssh密钥传输到 IP: $IP"
  ssh-copy-id "$REMOTE_USER@$IP"

  echo "传输router_torrc到远程IP: $IP"
  scp "/home/pi/router_torrc" "$REMOTE_USER@$IP:$REMOTE_FOLDER"
  # 运行远程命令
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "mkdir -p /home/pi/config/data"
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "mkdir -p /home/pi/logfile"
  
  echo "传输tor可执行文件到 IP: $IP"
  scp "/usr/local/bin/tor" "$REMOTE_USER@$IP:$REMOTE_FOLDER"

  
  # 检查传输是否成功
  if [ $? -eq 0 ]; then
    let success=success+1
    echo "成功"
  else
    echo "失败"
  fi
  
  echo ""
done

scp "/home/pi/hs_torrc" "$REMOTE_USER@$IP_HS:/root"
echo "成功传输$success台机器,期望是31台"
