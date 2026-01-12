#!/bin/bash
set -e

# build/tag image
# ${DOCKER} build -t quay.io/clyang82/inventory-api:latest -f Dockerfile .
# ${DOCKER} build -t quay.io/clyang82/inventory-e2e-tests:latest -f Dockerfile-e2e .
# ${DOCKER} build -t quay.io/clyang82/kafka-connect:latest -f Dockerfile.connect .
# 
# ${DOCKER} tag quay.io/clyang82/inventory-api:latest quay.io/clyang82/inventory-api:e2e-test
# ${DOCKER} tag quay.io/clyang82/inventory-e2e-tests:latest quay.io/clyang82/inventory-e2e-tests:e2e-test
# ${DOCKER} tag quay.io/clyang82/kafka-connect:latest quay.io/clyang82/kafka-connect:e2e-test
# 
# ${DOCKER} push quay.io/clyang82/inventory-api:e2e-test
# ${DOCKER} push quay.io/clyang82/inventory-e2e-tests:e2e-test
# ${DOCKER} push quay.io/clyang82/kafka-connect:e2e-test

# checks for the config map first, or creates it if not found
kubectl get configmap inventory-api-psks || kubectl create configmap inventory-api-psks --from-file=config/psks.yaml
[ -f resources.tar.gz ] || tar czf resources.tar.gz -C data/schema/resources .
kubectl get configmap resources-tarball || kubectl create configmap resources-tarball --from-file=resources.tar.gz

kubectl apply -f deploy/kind/inventory/kessel-inventory.yaml
kubectl apply -f deploy/kind/inventory/invdatabase.yaml


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
kubectl get crd httpproxies.projectcontour.io

# checks for the config map first, or creates it if not found
kubectl get configmap spicedb-schema || kubectl create configmap spicedb-schema --from-file=deploy/schema.zed

kubectl apply -f deploy/kind/relations/spicedb-kind-setup/postgres/secret.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/postgres/postgresql.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/postgres/storage.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/bundle.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/spicedb-cr.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/svc-ingress.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/relations-api/secret.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/relations-api/deployment.yaml
kubectl apply -f deploy/kind/relations/spicedb-kind-setup/relations-api/svc.yaml


echo "Waiting for all pods to be ready (1/1)..."
MAX_RETRIES=30


while true; do
  POD_STATUSES=$(kubectl get pods --no-headers)

  NOT_READY=$(echo "$POD_STATUSES" | awk '{print $2}' | grep -v '^1/1$' | wc -l)

  if [ "$NOT_READY" -eq 0 ]; then
    echo "All pods are ready (1/1)."
    echo "Delaying readiness checks to allow Kafka pods to initialize..."
    sleep 30
    break
  fi

  echo "Waiting for pods to be ready... ($NOT_READY pods not ready)"
  kubectl get pods
  sleep 5

done

kubectl apply -f deploy/kind/e2e/e2e-batch.yaml
echo "Setup complete."
rm -rf $TMP_DIR
rm -rf $KIND
