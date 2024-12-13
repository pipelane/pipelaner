
# **Console**

The **Console** sink component outputs messages to the console, supporting customizable log formats.

---

## **Config Definition**

```pkl
class Console extends Sink {
  fixed sourceName = "console"
  logFormat: LoggerConfig.LogFormat = "plain"
}
```

---

## **Attributes**

| **Attribute** | **Type**                 | **Description**                                     | **Default Value** |
|---------------|--------------------------|-----------------------------------------------------|--------------------|
| `logFormat`   | `LoggerConfig.LogFormat` | The format of the console logs (`plain` or `json`). | `"plain"`          |

---

## **I/O Types**

- **Input Type:** `any`
- **Output:** Outputs messages to the console in the specified format.

---

## **Pkl Configuration Example**

### **Basic Console Sink**
```pkl
new Sinks.Console {
  name = "example-console"
  logFormat = "plain"
}
```

### **Console Sink with JSON Format**
```pkl
new Sinks.Console {
  name = "example-console-json"
  logFormat = "json"
}
```

---

## **Description**

The **Console** sink outputs incoming messages to the console. It supports two log formats:
1. **Plain Text (`plain`)**: Logs messages as plain text.
2. **JSON (`json`)**: Logs messages in a structured JSON format.

---

### **Use Cases**

1. **Debugging**
    - Output pipeline data to the console for debugging purposes.
2. **Local Development**
    - Easily visualize pipeline data during development.

---

## **Notes**

- Ensure that the `logFormat` attribute matches your desired log format.
- Use the `plain` format for human-readable logs and `json` for structured logs suitable for further processing.

---
