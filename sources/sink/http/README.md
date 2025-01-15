
# **HTTP**

The **HTTP** sink component sends incoming messages to an HTTP endpoint using a specified method and optional headers.

---

## **Config Definition**

```pkl
typealias Method = "PATCH" | "POST" | "PUT" | "DELETE" | "GET"

class Http extends Sink {
  fixed sourceName = "http"
  url: String
  method: Method
  headers: Mapping<String, String>
}
```

---

## **Attributes**

| **Attribute** | **Type**                 | **Description**                                                   | **Default Value** |
|---------------|--------------------------|-------------------------------------------------------------------|--------------------|
| `url`         | `String`                 | The URL of the HTTP endpoint.                                     | **Required**      |
| `method`      | `Method`                 | The HTTP method to use (`PATCH`, `POST`, `PUT`, `DELETE`, `GET`). | **Required**      |
| `headers`     | `Mapping<String,String>` | A mapping of HTTP headers to include in the request.              | `{}` (empty map)  |

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
  headers = new Mapping {
     ["Authorization"] {
        "Bearer \(read("env:AUTH_TOKEN"))"
     }
     ["Content-Type"] {
        "application/json"
     }
  }
}
```

### **HTTP Sink with GET Method**
```pkl
new Http {
  name = "example-http-get"
  url = "https://example.com/resource"
  method = "GET"
  headers = new Mapping {
     ["Accept"] {
        "application/json"
     }
  }
}
```

---

## **Description**

The **HTTP** sink sends messages to an HTTP endpoint using the specified method. Optional headers can be used to customize requests for specific APIs or services. This sink is versatile and can be used for various purposes such as sending data, triggering APIs, or managing resources.

---

### **Use Cases**

1. **Data Transmission**
   - Send processed data to an API for further processing or storage.
2. **Triggering Webhooks**
   - Use methods like `POST` or `PATCH` to trigger webhooks with pipeline data.
3. **Resource Management**
   - Use methods like `DELETE` or `PUT` to manage resources on a remote server.
4. **Authenticated Requests**
   - Add authentication headers (e.g., `Authorization`) for secure communication with APIs.

---

## **Notes**

- Ensure the `url` attribute points to a valid HTTP endpoint.
- The `method` attribute must match the desired HTTP operation.
- Use the `headers` attribute to include custom headers like `Authorization` or `Content-Type` as required by your API.

---

If you have additional questions or need further clarification, feel free to reach out!
