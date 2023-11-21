#!/usr/bin/env bash

kubelet_logs () {
  while [[ -z $KIND_CONTAINER ]];
  do
    KIND_CONTAINER=$(podman container ls --format "{{.ID}}")
    sleep 5;
  done

  podman container exec -ti ${KIND_CONTAINER} journalctl -fu kubelet.service
}

/usr/bin/kind create cluster --config /kind-in-podman.yaml -v 8 &

kubelet_logs
