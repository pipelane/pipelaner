
# **Kafka**

The **Kafka** input component enables consuming messages from Kafka topics with configurable settings such as partition fetch sizes, offset reset policies, and balancing strategies.

---

## **Class Definition**

```pkl
class Kafka extends Input {
  fixed sourceName = "kafka"
  common: Common.Kafka
  autoCommitEnabled: Boolean = true
  consumerGroupID: String
  autoOffsetReset: AutoOffsetReset = "earliest"
  balancerStrategy: Listing<Strategy> = new Listing<Strategy> {
    "cooperative-sticky"
  }
  maxPartitionFetchBytes: DataSize(validateBuffersSizes) = 1.mib
  fetchMaxBytes: DataSize(validateBuffersSizes) = 50.mib

  hidden validateBuffersSizes = (_) ->
      if (fetchMaxBytes < maxPartitionFetchBytes)
        throw("'fetchMaxBytes' should be more than 'maxPartitionFetchBytes'")
      else true
}
```

---

## **Common.Kafka Definition**

The `common` attribute of the Kafka input component references the **Common.Kafka** class, which defines essential configurations for connecting to Kafka brokers and interacting with topics.

### **Common.Kafka Class**

```pkl
class Kafka {
  saslEnabled: Boolean = false
  saslMechanism: SASLMechanism(isSASLMechanizmSetted)?
  saslUsername: String(isSASLEnabled)?
  saslPassword: String(isSASLEnabled)?
  brokers: Listing<String>
  version: String?
  topics: Listing<String>

  hidden isSASLMechanizmSetted = (_) ->
      if (saslEnabled && saslMechanism == null)
        throw("'saslMechanism' can not be null")
      else true
  hidden isSASLEnabled = (_) ->
      if ((saslEnabled && saslUsername == null && saslPassword == null) ||
          (saslEnabled && saslUsername == "" && saslPassword == ""))
        throw("'saslUsername' and 'saslPassword' can not be empty string or null")
      else true
}
```

### **Common.Kafka Attributes**

| **Attribute**      | **Type**             | **Description**                                                                                       | **Default Value** |
|---------------------|----------------------|-------------------------------------------------------------------------------------------------------|--------------------|
| `saslEnabled`      | `Boolean`           | Enables or disables SASL authentication.                                                             | `false`           |
| `saslMechanism`    | `SASLMechanism`     | SASL mechanism to use (`SCRAM-SHA-512` or `SCRAM-SHA-256`).                                           | `null`            |
| `saslUsername`     | `String`            | SASL authentication username.                                                                        | `null`            |
| `saslPassword`     | `String`            | SASL authentication password.                                                                        | `null`            |
| `brokers`          | `Listing<String>`   | List of Kafka broker addresses.                                                                      | **Required**      |
| `version`          | `String`            | Kafka protocol version (optional).                                                                   | `null`            |
| `topics`           | `Listing<String>`   | List of Kafka topics to subscribe to.                                                                | **Required**      |

### **Validations**

1. **SASL Mechanism Validation**
    - If `saslEnabled` is `true` and `saslMechanism` is not set, the following error is thrown:
      ```
      'saslMechanism' can not be null
      ```

2. **SASL Credentials Validation**
    - If `saslEnabled` is `true`, both `saslUsername` and `saslPassword` must be provided. Otherwise, the following error is thrown:
      ```
      'saslUsername' and 'saslPassword' can not be empty string or null
      ```

---

## **Kafka Input Attributes**

| **Attribute**            | **Type**                 | **Description**                                                                                       | **Default Value**             |
|---------------------------|--------------------------|-------------------------------------------------------------------------------------------------------|--------------------------------|
| `common`                 | `Common.Kafka`          | Reusable Kafka connection settings (e.g., brokers, SASL).                                            | **Required**                  |
| `autoCommitEnabled`      | `Boolean`               | Enables or disables auto-commit for consumer offsets.                                                | `true`                        |
| `consumerGroupID`        | `String`                | Consumer group ID for managing Kafka consumers and partition ownership.                              | **Required**                  |
| `autoOffsetReset`        | `AutoOffsetReset`       | Behavior when there is no initial offset or when the offset is invalid (`earliest` or `latest`).     | `"earliest"`                  |
| `balancerStrategy`       | `Listing<Strategy>`     | Strategy used for partition assignment during rebalancing.                                           | `["cooperative-sticky"]`      |
| `maxPartitionFetchBytes` | `DataSize`              | Maximum data fetched per partition per request.                                                      | `1.mib`                       |
| `fetchMaxBytes`          | `DataSize`              | Maximum data fetched across all partitions per request.                                              | `50.mib`                      |

---

## **Validations**

### **Buffer Size Validation**
- Ensures `fetchMaxBytes` is greater than or equal to `maxPartitionFetchBytes`.
- If the validation fails, an exception is thrown:
  ```
  'fetchMaxBytes' should be more than 'maxPartitionFetchBytes'
  ```

---

## **Pkl Configuration Example**

### **Basic Kafka Input**
```pkl
new Kafka {
  common = new Common.Kafka {
    brokers = ["broker1:9092", "broker2:9092"]
    topics = ["example-topic"]
  }
  consumerGroupID = "example-consumer-group"
}
```

### **Kafka Input with SASL Authentication**
```pkl
new Kafka {
  common = new Common.Kafka {
    saslEnabled = true
    saslMechanism = "SCRAM-SHA-512"
    saslUsername = "example-user"
    saslPassword = "example-password"
    brokers = ["broker1:9092", "broker2:9092"]
    topics = ["secure-topic"]
  }
  consumerGroupID = "example-secure-consumer"
}
```

---

## **Attributes in Detail**

### **Common Kafka Settings (`common`)**
- Defines reusable Kafka connection attributes (e.g., brokers, topics, SASL authentication).

### **Auto Commit Enabled (`autoCommitEnabled`)**
- If `true`, offsets are automatically committed to Kafka.
- If `false`, manual offset commit is required.

### **Consumer Group ID (`consumerGroupID`)**
- Identifies the group of Kafka consumers that share load and maintain offset tracking.

### **Offset Reset Behavior (`autoOffsetReset`)**
- **`earliest`**: Start consuming from the earliest available message.
- **`latest`**: Start consuming from the latest message.

### **Partition Balancing Strategy (`balancerStrategy`)**
- Default: **`cooperative-sticky`** ensures minimal partition movement during rebalancing.
- Custom strategies can be added for advanced partitioning needs.

---

## **Use Cases**

1. **Basic Kafka Consumption**
    - Use the `common` attribute to specify brokers and topics, along with a consumer group ID for basic use cases.

2. **Secure Kafka Consumption**
    - Use `saslEnabled`, `saslMechanism`, `saslUsername`, and `saslPassword` to secure the Kafka connection.

3. **Optimized Data Transfer**
    - Adjust `maxPartitionFetchBytes` and `fetchMaxBytes` to fine-tune data fetching and improve performance.

---

## **Notes**

- Ensure `Common.Kafka` is configured with valid brokers and topics.
- Validate buffer sizes (`fetchMaxBytes` and `maxPartitionFetchBytes`) to avoid runtime errors.

---
