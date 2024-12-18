
# **Debounce**

The **Debounce** transform component limits the rate of message processing by introducing a delay (debounce interval) between consecutive messages.

---

## **Class Definition**

```pkl
class Debounce extends Transform {
  fixed sourceName = "debounce"
  interval: Duration
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                | **Default Value** |
|---------------|------------|------------------------------------------------|--------------------|
| `interval`    | `Duration` | The debounce interval between consecutive messages. | **Required**      |

---

## **I/O Types**

- **Input Type:** `any`
- **Output Type:** `any`

---

## **Pkl Configuration Example**

### **Basic Debounce Transform**
```pkl
new Transforms.Debounce {
  name = "example-debounce"
  interval = 500.ms
}
```

---

## **Description**

The **Debounce** transform component ensures that messages are processed with a minimum interval between them. This is useful for controlling the rate of message processing, reducing load, or avoiding repeated processing within a short period.

### **Use Cases**
1. **Rate Limiting**
    - Reduce the rate of messages sent to downstream components.
2. **Event Filtering**
    - Prevent rapid, repeated processing of the same or similar events.

---

## **Notes**

- Configure the `interval` attribute based on the required debounce duration for your use case.
- Debouncing helps in optimizing resource usage and preventing overload in downstream components.

---

If you have additional questions or need further clarification, feel free to reach out!
