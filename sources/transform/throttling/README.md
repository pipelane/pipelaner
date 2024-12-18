
# **Throttling**

The **Throttling** transform component controls the rate of message processing by introducing a configurable time interval between consecutive messages.

---

## **Class Definition**

```pkl
class Throttling extends Transform {
   fixed sourceName = "throttling"
   interval: Duration
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                               | **Default Value** |
|---------------|------------|---------------------------------------------------------------|--------------------|
| `interval`    | `Duration` | The interval between consecutive messages.                   | **Required**      |

---

## **I/O Types**

- **Input Type:** `any`
- **Output Type:** `any`

---

## **Pkl Configuration Example**

### **Basic Throttling Transform**
```pkl
new Transforms.Throttling {
  name = "example-throttling"
  interval = 1.s
}
```

---

## **Description**

The **Throttling** transform ensures that messages are processed at a controlled rate by applying a delay (`interval`) between consecutive messages.

---

### **Use Cases**

1. **Rate Limiting**
   - Control the rate at which messages are processed by downstream components.
2. **Load Management**
   - Prevent resource overload by throttling the message flow.

---

## **Notes**

- The `interval` attribute is required and should be configured based on system requirements.
- Use throttling to manage resource usage and maintain system stability.

---
