## General

- Two archs: point-to-point, pub-sub (event-driven)
- Even-driven: producer (sender) of the event (like a new user register) does not have to know who receives the notification. It just put the message on Queue. Highly scalable. 
- Queues are managed by a message broker like RabbitMQ
- RabbitMQ uses AMQP (Advanced Message Queuing Protocol)) protocol. So any library or software that uses the same protocol as client can use it. Respectively, we can swap rabbitMq with any other message borker that uses AMQP.
- There are other protocols also such as MQTT (IOT) and STOMP (WEB)
- Other brokers exists that have their own protocol like Kafka or AWS SQS

## RabbitMQ

- we have exhange, key, and queue
- sender push the message to an exchange with a key
- the exchange on RabbitMq cluster decides in which Queue to put it --> based on the routing key
- How exchange match the incoming request with the routing key? --> Direct (exact match), Topic (pattern matching wildcard), Fanout (all queues), Header (based on reqeust header not the key itself)
- Queues are either: Durable (persist on disk) or Transient (exist only on RAM).



