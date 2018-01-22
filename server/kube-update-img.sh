#!/bin/bash
# $1 : namespace
# $2 : type
# $3 : name
# $4 : container name
# $5 : image

echo "Update $2/$3 with image: $5"

img=$(kubectl -n $1 get $2/$3 -o jsonpath="{..image}")
echo "Current image: $img"
if [ $img != $5 ]
then
  echo "Different image, use kubectl set image" 
  kubectl -n $1 set image $2/$3 $4=$5 
else
  echo "Same image, try to kill all pods in deploy $2/$3" 
  while IFS= read -r pod; do
    # echo "delete $pod"
    kubectl delete -n $1 $pod
  done < <(kubectl -n $1 get pods -o=name | grep ^pods/$3-)
fi