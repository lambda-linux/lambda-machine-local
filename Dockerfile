FROM lambdalinux/baseimage-lambda:2017.03-004

CMD ["/sbin/my_init"]

COPY [ \
  "./docker-extras/etc-sudoers.d-docker", \
  "/tmp/docker-build/" \
]

RUN \
  # yum
  yum update && \
  yum install bzip2 && \
  yum install python27-PyYAML && \
  yum install python27-mako && \
  yum install sudo && \
  yum install tar && \
  \
  # setup sudo
  usermod -a -G wheel ll-user && \
  cp /tmp/docker-build/etc-sudoers.d-docker /etc/sudoers.d/docker && \
  chmod 440 /etc/sudoers.d/docker && \
  \
  # for python debugging (uncomment when required)
  # yum install gcc48 && \
  # yum install python27-devel && \
  # yum install python27-pip && \
  # su -l ll-user -c "pip-2.7 install --user ipdb==0.8 ipython==5.3.0" && \
  \
  # install github-release
  curl -L https://github.com/aktau/github-release/releases/download/v0.6.2/linux-amd64-github-release.tar.bz2 -o /tmp/docker-build/linux-amd64-github-release.tar.bz2 && \
  pushd /tmp/docker-build && \
  tar xjvf linux-amd64-github-release.tar.bz2 && \
  cp bin/linux/amd64/github-release /usr/bin && \
  popd && \
  \
  # cleanup
  rm -rf /tmp/docker-build && \
  yum clean all && \
  rm -rf /var/cache/yum/* && \
  rm -rf /tmp/* && \
  rm -rf /var/tmp/*

