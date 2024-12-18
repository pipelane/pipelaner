
# **Kafka**

The **Kafka** sink component sends incoming messages to a Kafka topic with configurable batching and message limits.

---

## **Config Definition**

```pkl
class Kafka extends Sink {
  fixed sourceName = "kafka"
  common: Common.Kafka
  maxRequestSize: DataSize = 10.mib
  lingerMs: Duration = 0.ms
  batchNumMessages: Int = 100_000
}
```

---

## **Attributes**

| **Attribute**       | **Type**      | **Description**                                                                       | **Default Value** |
|----------------------|---------------|---------------------------------------------------------------------------------------|--------------------|
| `common`            | `Common.Kafka`| Reusable Kafka connection configurations.                                             | **Required**      |
| `maxRequestSize`    | `DataSize`    | The maximum size of a Kafka request batch.                                            | `10.mib`          |
| `lingerMs`          | `Duration`    | The maximum amount of time to wait before sending a batch of messages.                | `0.ms`            |
| `batchNumMessages`  | `Int`         | The maximum number of messages in a batch.                                            | `100,000`         |

---

## **I/O Types**

- **Input Type:** `string`, `[]byte`, `json objects`, `chan string`, `chan []byte`, `chan json objects`.
- **Output:** Sends messages to the specified Kafka topic.

---

## **Pkl Configuration Example**

### **Basic Kafka Sink**
```pkl
new Sinks.Kafka {
  name = "example-kafka"
  common = new Common.Kafka {
    brokers { 
      "broker1:9092"
      "broker2:9092" 
    }
    topics { 
      "example-topic" 
    }
  }
}
```

### **Kafka Sink with Custom Batching**
```pkl
new Sinks.Kafka {
  name = "example-kafka-batched"
  common = new Common.Kafka {
    brokers {
      "broker1:9092"
      "broker2:9092"
    }
    topics { 
      "example-topic" 
    }
  }
  maxRequestSize = 5.mib
  lingerMs = 10.ms
  batchNumMessages = 50_000
}
```

---

## **Description**

The **Kafka** sink component allows sending processed pipeline messages to a Kafka topic. It supports advanced configuration options for controlling the size and frequency of message batches.

---

### **Use Cases**

1. **High-Throughput Messaging**
    - Send large volumes of messages to Kafka efficiently by leveraging batching.
2. **Low-Latency Applications**
    - Adjust `lingerMs` to prioritize immediate message delivery.
3. **Optimized Resource Usage**
    - Configure `batchNumMessages` and `maxRequestSize` to reduce the overhead of frequent requests.

---

## **Notes**

- Ensure the `common` attribute is configured with valid Kafka connection settings.
- Use `maxRequestSize` and `batchNumMessages` to fine-tune performance based on your use case.
- Increasing `lingerMs` can improve throughput but may introduce slight delays in message delivery.

---