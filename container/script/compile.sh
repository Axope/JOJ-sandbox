#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <memory_limit_in_bytes> <source_file>"
    exit 1
fi

memory_limit="$1"
source_file="$2"

if [ ! -f "$source_file" ]; then
    echo "Source solution file '$source_file' not found."
    exit 1
fi

echo $memory_limit > /sys/fs/cgroup/memory/memory.limit_in_bytes

pid=$$
echo $pid > /sys/fs/cgroup/memory/cgroup.procs

output="./output"

start_time=$(date +%s%N)
g++ "$source_file" -o solution -Wall -O2 -std=c++17 2> $output/compile_info
end_time=$(date +%s%N)
runtime=$((end_time - start_time))
echo $runtime > $output/compile_time