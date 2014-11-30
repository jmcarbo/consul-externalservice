rm -rf /tmp/consul
consul agent -server -bootstrap -data-dir /tmp/consul &
CONSUL_PID=$!
sleep 5
bin/consul-externalservice_darwin_amd64 start --node host1 &
CONSUL_ES1_PID=$!
bin/consul-externalservice_darwin_amd64 start --node host2 &
CONSUL_ES2_PID=$!
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/host1/aservice \
    -d '{ "Address":"localhost2", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"running" }'
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/host2/aservice \
    -d '{ "Address":"localhost1", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"running" }'
sleep 5
dig @localhost -p 8600 aservice.service.consul
curl -XPUT http://localhost:8500/v1/kv/ExternalServices/host2/aservice \
    -d '{ "Address":"localhost1", "Port":80, "Interval":"1s", "Command":"ping -c 1 localhost", "TargetState":"stopped" }'
sleep 5
dig @localhost -p 8600 aservice.service.consul

kill -TERM "$CONSUL_PID"
kill -TERM "$CONSUL_ES1_PID"
kill -TERM "$CONSUL_ES2_PID"
