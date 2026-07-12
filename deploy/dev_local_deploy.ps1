$VERSION = "v1.0.8"

pushd ..

# 先删容器，避免名称冲突
docker rm -f user_mgr
docker rmi -f user_mgr:$VERSION
docker build -t user_mgr:$VERSION .
docker run -d --name user_mgr --network local_deps_install_dev_net -p 30880:8080 user_mgr:$VERSION
docker ps
docker logs user_mgr -f

popd