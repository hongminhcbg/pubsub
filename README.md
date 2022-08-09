I. Pubsub handler
    
    m.Ack() mark message process success, it will be removed from queue
    m.Nack() mark message process fail, it will be resended now
    retendtion = [10min:7day], when exceeded retendion, message not consume will be removed

details:

pub/sub is async with latency about 100ms 

core_concepts:

topic: a resource publisher send message

message: combination of data and message attributes 

message_attributes: key-value pair

ack: signal subcriber send to pub/sub after successfully

quatas: https://cloud.google.com/pubsub/quotas

pub: 2GB/s or 200MB/s

sub: 4GB/s or 400MB/s

ACK: 4GB/s or 400MB/s

ordering: if a publisher sends two messages with the same ordering key, the Pub/Sub service delivers the oldest message first. order key should client_id or user_id

TODO: ordering message (https://cloud.google.com/pubsub/docs/ordering) + dead queue (https://cloud.google.com/pubsub/docs/handling-failures)


Notification question:
  
  what happen when noti is expired? 
  context: user make a transaction but after that 5min, the noti will come
  => add message attributes
  
  what happen when third party rate limit exceeded? => use redis to rate limit

  what happen when message queue rate-limit (pubsub 200MB/s, redis out of storage)? 
  => rate limit in receiver api

  how to use many providers? 
  => use round robin to choose

  design for priority notification type?

  db choose? why? what happen when save to db error or send to third party error? save before send or send before save?
