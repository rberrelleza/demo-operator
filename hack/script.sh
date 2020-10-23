export IMG=ramiro/demo-operator:0.1.0
make docker-build docker-push IMG=$IMG
make install
make deploy IMG=$IMG

kubectl apply -f config/samples/hello_v1beta1_saiyam.yaml

kubectl logs -f -l=control-plane=controller-manager -n=demo-operator-system