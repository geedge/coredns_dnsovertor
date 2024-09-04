#! /bin/bash
# 定义远程服务器的用户名和目标文件夹路径
source config.sh
i=1
success=0
failure=0
# 循环遍历 IP 列表并进行文件传输
for IP in "${IP_LIST2[@]}"; do
  echo "在远程IP上进行tor关闭//tor: $IP"
  # 运行远程命令

  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "killall tor"
 
  # 检查关闭是否成功
  if [ $? -eq 0 ]; then
    let success=success+1
    echo "关闭tor成功"
  else
    let failure=failure+1
    echo "关闭tor失败"
  fi
  
  echo ""
done
echo "成功$success台router节点,期望情况29台"
echo "失败$failure台router节点,期望情况0台"
