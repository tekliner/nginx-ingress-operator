#!/bin/bash
set -e

operator-sdk generate k8s
operator-sdk build 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator:debug
eval $(aws ecr get-login --no-include-email --region us-east-1 --profile staging | sed 's|https://||')
docker push 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator:debug
kubectl delete po -l name=nginx-ingress-operator

