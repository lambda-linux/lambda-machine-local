# Release

Build and run docker image

```
docker build -t lambda-machine-local-release .

docker run --rm -v `pwd`:`pwd` \
  -ti lambda-machine-local-release \
  /sbin/my_init -- /sbin/setuser ll-user /bin/bash -l
```

To create release, first update, commit and push `release/release.yaml`. Then

```
./release --github_security_token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --bin_dir ../bin \
  --pre_release=1

./release --github_security_token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  --bin_dir ../bin
```
