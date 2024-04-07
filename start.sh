rm ./container/solution 1>/dev/null 2>&1
rm ./container/output/* 1>/dev/null 2>&1

containerName="test_cpp"
docker stop $containerName
docker rm $containerName

# 启动之前需要保证所有数据文件都存在 `container` 中
docker run -itd \
    --name $containerName \
    --net=none \
    --ulimit core=0 \
    -v /sys/fs/cgroup/memory/testMem:/sys/fs/cgroup/memory \
    -v ./container:/root \
    judger-cpp:0.1

docker exec -it $containerName bash -c 'cd /root && ./sandbox -type=compile'
if [ ! -e "./container/solution" ]; then
    echo "compile error, please cat compile_info"
    exit 1
fi

docker exec -it $containerName bash -c 'cd /root && ./sandbox -type=run'
echo result code = $?
