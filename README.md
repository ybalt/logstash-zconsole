# logstash-zconsole

small Go program that takes logs from Logstash via zmq transport and out it as http endpoint
Primary for debugging/development

Run:
```
# docker run --rm -p 8080:8080 -it -e LOGSTASH_ADDR=tcp://127.0.0.1:12300 ybalt/logstash-zconsole
```
change LOGSTASH_ADDR to Logstash zmq pub endpoint (topics not implemented)

Logstash config
```
output {
		zeromq {
            	topology => "pubsub"
            	address => ["tcp://*:12300"]
            	mode => "server"
            	codec => "json"
            }
}
```
Connect:
```
# curl 127.0.0.1:8080
```
