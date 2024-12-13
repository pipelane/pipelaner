
# **Batch**

The **Batch** transform component groups incoming messages into batches of a specified size for further processing.

---

## **Config Definition**

```pkl
class Batch extends Transform {
  fixed sourceName = "batch"
  size: UInt32
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                              | **Default Value** |
|---------------|------------|--------------------------------------------------------------|--------------------|
| `size`        | `UInt32`   | The number of messages to include in each batch.             | **Required**      |

---

## **I/O Types**

- **Input Type:** `any`
- **Output Type:** `chan any`

---

## **Pkl Configuration Example**

### **Basic Batch Transform**
```pkl
new Batch {
  name = "example-batch"
  size = 100
}
```

---

## **Description**

The **Batch** transform collects a specified number of incoming messages and processes them as a single unit (batch). This is useful for reducing the overhead of individual message processing or aggregating messages before further transformations or outputs.

### **Use Cases**
1. **Batch Processing**
    - Efficiently handle large volumes of data by processing them in batches.
2. **Data Aggregation**
    - Aggregate data for downstream components that require grouped input.

---

## **Notes**

- Ensure the `size` attribute is configured appropriately for your use case.
- The batch size should be optimized based on the downstream component's processing capabilities.

---
