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
- if multiple consumer are connected to a queue, the message get consumed by only one them based on a round-robbin algorithm (or any other). 
- if all consumer need to get a message at the same time, it's better to have one queue for each.
- routing keys are made up words separated with dot. Example: user.post.follow
- routing keys can include * (exactly one) or # (zero or more). for example user.*.won does NOT match user.won . It matched user.game.won
- naming convention for queue and key is that the queue name should include the key. Example: key=user.created  queue=user.created.email_notifier
- key name should be: noun.verb




