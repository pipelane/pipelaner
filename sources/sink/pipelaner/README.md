
# **Pipelaner**

The **Pipelaner** sink component sends messages to another Pipelaner node, allowing communication between pipeline instances.

---

## **Config Definition**

```pkl
class Pipelaner extends Sink {
  fixed sourceName = "pipelaner"
  commonConfig: Common.Pipelaner
}
```

---

## **Attributes**

| **Attribute**    | **Type**           | **Description**                                    | **Default Value** |
|-------------------|--------------------|----------------------------------------------------|--------------------|
| `commonConfig`   | `Common.Pipelaner` | Reusable configuration for connecting to a Pipelaner node. | **Required**      |

---

## **I/O Types**

- **Input Type:** `string`, `[]byte`, `json objects`
- **Output:** Sends messages to another Pipelaner node.

---

## **Pkl Configuration Example**

### **Basic Pipelaner Sink**
```pkl
new Pipelaner {
  name = "example-pipelaner"
  commonConfig = new Common.Pipelaner {
    host = "127.0.0.1"
    port = 8080
  }
}
```

### **Pipelaner Sink with TLS**
```pkl
new Sinks.Pipelaner {
  name = "example-pipelaner-secure"
  commonConfig = new Common.Pipelaner {
    host = "127.0.0.1"
    port = 8443
    tls = new TLSConfig {
      certFile = "/path/to/cert.pem"
      keyFile = "/path/to/key.pem"
    }
  }
}
```

---

## **Description**

The **Pipelaner** sink component allows you to send processed pipeline data to another Pipelaner node, enabling distributed processing across multiple instances.

---

### **Use Cases**

1. **Distributed Processing**
    - Send pipeline data to another node for further processing or transformation.
2. **Scalable Pipelines**
    - Use multiple Pipelaner nodes to distribute workload and increase throughput.
3. **Secure Communication**
    - Configure `tls` for encrypted communication between nodes.

---

## **Notes**

- Ensure the `commonConfig` attribute is correctly configured for the target Pipelaner node.
- Use TLS for secure communication in production environments.

---

If you have additional questions or need further clarification, feel free to reach out!
