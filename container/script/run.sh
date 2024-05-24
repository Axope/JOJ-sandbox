#!/bin/bash

if [ $# -ne 4 ]; then
    echo "Usage: $0 <memory_limit_in_bytes> <source_solution_file> <test_case> <lang>"
    exit 1
fi

memory_limit="$1"
source_solution_file="$2"
test_case="$3"
lang="$4"

if [ ! -f "$source_solution_file" ]; then
    echo "Source solution file '$source_solution_file' not found."
    exit 1
fi

echo "$memory_limit" > /sys/fs/cgroup/memory/memory.limit_in_bytes

pid=$$
echo $pid > /sys/fs/cgroup/memory/cgroup.procs

output="./output"

cd /root
start_time=$(date +%s%N)
# ./$source_solution_file < ./data/$test_case.in 1>$output/$test_case.out 2>$output/$test_case.err
# return_code=$?
case "$lang" in
  0)
    ./$source_solution_file < ./data/$test_case.in 1>$output/$test_case.out 2>$output/$test_case.err
    return_code=$?
    ;;
#   1)
#     java -cp $(dirname $source_solution_file) $(basename $source_solution_file .java) < ./data/$test_case.in 1>$output/$test_case.out 2>$output/$test_case.err
#     return_code=$?
#     ;;
  2)
    python3 $source_solution_file < ./data/$test_case.in > $output/$test_case.out
    return_code=$?
    ;;
#   3)
#     go run $source_solution_file < ./data/$test_case.in 1>$output/$test_case.out 2>$output/$test_case.err
#     return_code=$?
#     ;;

esac
end_time=$(date +%s%N)

if [ $return_code -ne 0 ]; then
    echo "segment fault"
    exit 3
fi

runtime=$((end_time - start_time))
echo $runtime > $output/real_runtime_$test_case