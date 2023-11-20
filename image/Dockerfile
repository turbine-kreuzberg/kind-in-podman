# sudo sysctl -w user.max_user_namespaces=72423
FROM ubuntu:22.04

USER root
RUN apt update && \
    apt dist-upgrade -y && \
    apt install gpg curl lsb-release sudo vim -y

RUN mkdir -p /etc/apt/keyrings && \
    curl -fsSL "https://download.opensuse.org/repositories/devel:kubic:libcontainers:unstable/xUbuntu_$(lsb_release -rs)/Release.key" \
      | gpg --dearmor \
      | tee /etc/apt/keyrings/devel_kubic_libcontainers_unstable.gpg > /dev/null && \
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/devel_kubic_libcontainers_unstable.gpg]\
        https://download.opensuse.org/repositories/devel:kubic:libcontainers:unstable/xUbuntu_$(lsb_release -rs)/ /" \
      | tee /etc/apt/sources.list.d/devel:kubic:libcontainers:unstable.list > /dev/null && \
    apt update && \
    apt install -y fuse-overlayfs podman && \
    apt clean

RUN cp /usr/share/containers/storage.conf /etc/containers/storage.conf && \
    mkdir -p /var/lib/shared && \
    sed -i -e 's|^#mount_program|mount_program|g' \
      -e '/additionalimage.*/a "/var/lib/shared",' \
    #      -e 's|^#ignore_chown_errors.*|ignore_chown_errors = "true"|' \
      -e 's|^mountopt[[:space:]]*=.*$|mountopt = "nodev,fsync=0"|g' /etc/containers/storage.conf && \
    chmod 4755 /usr/bin/newgidmap && \
    chmod 4755 /usr/bin/newuidmap


RUN curl -Lo /usr/bin/kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64 && \
    chmod +x /usr/bin/kind

RUN curl -Lo /usr/bin/kubectl https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

COPY rootfs/ /

#RUN useradd -s /bin/bash -d /home/kind/ -m kind && \
#    chown -R kind:kind /home/kind && \
#    echo kind:10000:5000 > /etc/subuid && \
#    echo kind:10000:5000 > /etc/subgid

#USER kind
#WORKDIR /home/kind

VOLUME /var/lib/containers
#VOLUME /home/kind/.local/share/containers

ENTRYPOINT ["/create-cluster.sh"]
#CMD ["create", "cluster", "--config", "/kind-in-podman.yaml"]

# 65534 - ( 10000 + 5000 - 1)
# 1001020000 - (100000 + 65536 - 1) = 1000854465

# podman run --rm  -ti -v $(pwd)/kubeconfig:/root/.kube/config -p 6443:6443 --cap-add=sys_admin,mknod,net_admin --device=/dev/fuse --security-opt="label=disable" --entrypoint=bash localhost/kind-in-podman:latest
