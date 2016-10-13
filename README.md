# logstash-zconsole

small Go program that Docker container logs from Logstash via zmq transport and out it as http endpoint.
Primary for debugging/development environments based on ELK stack.

1. 
Logstash config
```
input {
    gelf {
                port => 12200
        	host => "0.0.0.0"
        	codec => json
        	type => "docker_logs"
    }
}
output {
    zeromq {
         	topology => "pubsub"
            	address => ["tcp://*:12300"]
            	mode => "server"
            	codec => "json"
    }
}
```

2. 
Configure Docker for Logstash GELF logging, adding this line to docker daemon option
```
# dockerd --log-driver=gelf --log-opt gelf-address=udp://<logstash_ip>:12200 
```
or by adding line to docker run
```
# docker run --log-driver=gelf --log-opt gelf-address=udp://<logstash_ip>:12200 <image> 
```

3. 
Run zconsole
```
# docker run --rm -p 8080:8080 -it -e LOGSTASH_ADDR=tcp://<logstash_ip>:12300 ybalt/logstash-zconsole
```
change LOGSTASH_ADDR to Logstash zmq pub endpoint, IP should be visible from this container.

4.
Connect:
```
# curl 127.0.0.1:8080
```
