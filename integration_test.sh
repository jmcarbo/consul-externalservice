rm -rf /tmp/consul
consul agent -server -bootstrap -data-dir /tmp/consul &
sleep 5
bin/consul-externalservice_darwin_amd64 start --node host1 &
bin/consul-externalservice_darwin_amd64 start --node host2 &
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/host1/aservice \
    -d '{ "Address":"localhost2", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"running" }'
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/host2/aservice \
    -d '{ "Address":"localhost1", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"running" }'
