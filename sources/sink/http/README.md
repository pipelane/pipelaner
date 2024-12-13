
# **HTTP**

The **HTTP** sink component sends incoming messages to an HTTP endpoint using a specified method.

---

## **Config Definition**

```pkl
typealias Method = "PATCH" | "POST" | "PUT" | "DELETE" | "GET"

class Http extends Sink {
  fixed sourceName = "http"
  url: String
  method: Method
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                                   | **Default Value** |
|---------------|------------|-------------------------------------------------------------------|--------------------|
| `url`         | `String`   | The URL of the HTTP endpoint.                                     | **Required**      |
| `method`      | `Method`   | The HTTP method to use (`PATCH`, `POST`, `PUT`, `DELETE`, `GET`). | **Required**      |

---

## **I/O Types**

- **Input Type:** `map[string]any`, `[]byte`, `string`, `json structs`
- **Output:** Sends messages to the specified HTTP endpoint.

---

## **Pkl Configuration Example**

### **Basic HTTP Sink**
```pkl
new Http {
  name = "example-http"
  url = "https://example.com/api"
  method = "POST"
}
```

### **HTTP Sink with GET Method**
```pkl
new Http {
  name = "example-http-get"
  url = "https://example.com/resource"
  method = "GET"
}
```

---

## **Description**

The **HTTP** sink sends messages to an HTTP endpoint using the specified method. It is versatile and can be used for various purposes such as sending data, triggering APIs, or deleting resources.

---

### **Use Cases**

1. **Data Transmission**
    - Send processed data to an API for further processing or storage.
2. **Triggering Webhooks**
    - Use methods like `POST` or `PATCH` to trigger webhooks with pipeline data.
3. **Resource Management**
    - Use methods like `DELETE` or `PUT` to manage resources on a remote server.

---

## **Notes**

- Ensure the `url` attribute points to a valid HTTP endpoint.
- The `method` attribute must match the desired HTTP operation.

---

If you have additional questions or need further clarification, feel free to reach out!
