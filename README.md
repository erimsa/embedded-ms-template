version: '2'
services:
  heat-sensor-service:
    image: heat-sensor-service
    command: [./wait-for-it.sh, "elasticsearch:9200", "--",  "/heat-sensor-service"]
    build:
      context: "./heat-sensor-service"
      dockerfile: Dockerfile
    
    ports:
      - 3000:3000
    depends_on:
      - redis
      - elasticsearch
    links:
      - redis
      - elasticsearch
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.11.1
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - esdata1:/usr/share/elasticsearch/data
    environment: 
      - discovery.type=single-node
  redis:
    #image: eliar/redis
    image: redis
    ports:
      - "6379:6379"

volumes:
  esdata1:
    driver: local