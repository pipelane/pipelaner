
# **Chunks**

The **Chunks** transform component splits incoming messages into smaller chunks based on a maximum size or idle time.

---

## **Class Definition**

```pkl
class Chunk extends Transform {
  fixed sourceName = "chunks"
  maxChunkSize: UInt32
  maxIdleTime: Duration
}
```

---

## **Attributes**

| **Attribute**      | **Type**   | **Description**                                                               | **Default Value** |
|---------------------|------------|-------------------------------------------------------------------------------|--------------------|
| `maxChunkSize`     | `UInt32`   | The maximum number of messages to include in each chunk.                     | **Required**      |
| `maxIdleTime`      | `Duration` | The maximum idle time before creating a chunk with fewer than `maxChunkSize`. | **Required**      |

---

## **I/O Types**

- **Input Type:** `any`
- **Output Type:** `chan any`

---

## **Pkl Configuration Example**

### **Basic Chunks Transform**
```pkl
new Chunk {
  name = "example-chunks"
  maxChunkSize = 50
  maxIdleTime = 1.s
}
```

---

## **Description**

The **Chunks** transform component breaks down a continuous stream of messages into smaller, manageable chunks. Chunks are created either when the maximum number of messages (`maxChunkSize`) is reached or when the stream remains idle for a specified duration (`maxIdleTime`).

### **Use Cases**
1. **Stream Segmentation**
    - Split a continuous data stream into smaller, manageable parts for downstream processing.
2. **Rate Control**
    - Use `maxIdleTime` to ensure timely chunk delivery even when the stream has low activity.

---

## **Notes**

- Configure both `maxChunkSize` and `maxIdleTime` to balance performance and latency.
- The choice of `maxChunkSize` and `maxIdleTime` should align with the requirements of downstream components.

---
