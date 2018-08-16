#! /bash/sh
declare -a dic
dic=()
for((i=1;i<50;i++));
do
  res=$(cd /Users/luobangkui/myfolder/go_test/aboutLearnNSQ;go run about_wrapper.go)
  let dic[$res]+=1
done
for i in $(echo ${!dic[*]})
do
  echo $i:${dic[$i]}
done
