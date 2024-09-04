#! /bin/bash
# 定义远程服务器的用户名和目标文件夹路径
# 定义要传输的 IP 列表
source config.sh
i=1
success=0
failure=0
# 循环遍历 IP 列表并进行文件传输
for IP in "${IP_LIST2[@]}"; do
  echo "ssh远程执行uname -a，检查能否成功连接第$i台机器: $IP"
  # 运行远程命令
  let i=i+1
  sshpass -p "$REMOTE_PASSWORD" ssh -t "$REMOTE_USER@$IP" "uname -a"

  # 检查传输是否成功
  if [ $? -eq 0 ]; then
    let success=success+1
    echo "连接成功"
  else
    let failure=failure+1
    echo "连接失败"
  fi
  
  echo ""
done
#其他3台手动启动，使用root权限。
echo "正常连接$success台机器,期望情况是29台"
echo "连接失败$failure台机器,期望情况是0台"
