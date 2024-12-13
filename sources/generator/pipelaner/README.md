
# **Pipelaner**

The **Pipelaner** input component allows communication between Pipelaner nodes, supporting different connection types, including Unix sockets and HTTP/2.

---

## **Class Definition**

```pkl
typealias ConnectionType = "unix" | "http2"

class Pipelaner extends Input {
  fixed sourceName = "pipelaner"
  commonConfig: Common.Pipelaner(isUnixScoketConnection)?
  connectionType: ConnectionType = "http2"
  unixSocketPath: String(isUnixPathSetted)?

  hidden isUnixPathSetted = (_) ->
      if (connectionType == "unix" && unixSocketPath == null)
        throw("if you use unix socket, please set up 'unixSocketPath'")
      else true

  hidden isUnixScoketConnection = (_) ->
      if (connectionType != "unix" && commonConfig == null)
        throw("if you use http2 socket, please set up 'commonConfig'")
      else true
}
```

---

## **Common.Pipelaner Definition**

The `commonConfig` attribute of the Pipelaner input component references the **Common.Pipelaner** class, which defines basic connection settings for communication between Pipelaner nodes.

### **Common.Pipelaner Class**

```pkl
class TLSConfig {
  certFile: String(isTLSKeyAndCertSetted)
  keyFile: String(isTLSKeyAndCertSetted)
  hidden isTLSKeyAndCertSetted = (_) ->
      if (keyFile == "" || certFile == "")
        throw("key and cert files required with enabled tls")
      else true
}

class Pipelaner {
  host: String = "localhost"
  port: Int
  tls: TLSConfig?
}
```

### **Common.Pipelaner Attributes**

| **Attribute**      | **Type**        | **Description**                                                      | **Default Value** |
|---------------------|-----------------|----------------------------------------------------------------------|--------------------|
| `host`             | `String`       | Hostname of the Pipelaner node.                                      | `"localhost"`     |
| `port`             | `Int`          | Port number of the Pipelaner node.                                   | **Required**      |
| `tls`              | `TLSConfig`    | TLS configuration for secure communication (optional).               | `null`            |

### **TLSConfig Attributes**

| **Attribute**      | **Type**    | **Description**                                             | **Default Value** |
|---------------------|-------------|-------------------------------------------------------------|--------------------|
| `certFile`         | `String`    | Path to the TLS certificate file.                          | **Required**      |
| `keyFile`          | `String`    | Path to the TLS key file.                                  | **Required**      |

### **Validations**

1. **TLS Key and Certificate Validation**
    - If TLS is enabled, both `certFile` and `keyFile` must be provided. Otherwise, the following error is thrown:
      ```
      key and cert files required with enabled tls
      ```

2. **HTTP2 Socket Validation**
    - If `connectionType` is not `"unix"` and `commonConfig` is not set, the following error is thrown:
      ```
      if you use http2 socket, please set up 'commonConfig'
      ```

3. **Unix Socket Path Validation**
    - If `connectionType` is `"unix"` and `unixSocketPath` is not set, the following error is thrown:
      ```
      if you use unix socket, please set up 'unixSocketPath'
      ```

---

## **Pipelaner Input Attributes**

| **Attribute**       | **Type**              | **Description**                                                                                  | **Default Value** |
|----------------------|-----------------------|--------------------------------------------------------------------------------------------------|--------------------|
| `commonConfig`      | `Common.Pipelaner`    | Basic connection settings for communication between Pipelaner nodes (required for HTTP2).        | `null`            |
| `connectionType`    | `ConnectionType`      | Specifies the type of connection (`http2` or `unix`).                                            | `"http2"`         |
| `unixSocketPath`    | `String`             | Path to the Unix socket, required when `connectionType` is set to `"unix"`.                     | `null`            |

---

## **Pkl Configuration Example**

### **Pipelaner Input with HTTP/2**
```pkl
new Pipelaner {
  connectionType = "http2"
  commonConfig = new Common.Pipelaner {
    host = "127.0.0.1"
    port = 8080
  }
}
```

### **Pipelaner Input with Unix Socket**
```pkl
new Pipelaner {
  connectionType = "unix"
  unixSocketPath = "/tmp/pipelaner.sock"
}
```

### **Pipelaner Input with TLS**
```pkl
new Pipelaner {
  connectionType = "http2"
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

## **Attributes in Detail**

### **Common Configurations (`commonConfig`)**
- Defines reusable connection settings for Pipelaner nodes.
- Required for `http2` communication.

### **Connection Type (`connectionType`)**
- **`http2`**: Uses HTTP/2 for communication. Requires `commonConfig`.
- **`unix`**: Uses Unix socket for communication. Requires `unixSocketPath`.

### **Unix Socket Path (`unixSocketPath`)**
- Required only when `connectionType` is set to `"unix"`.
- Defines the path to the Unix socket for communication.

### **TLS Configuration (`tls`)**
- Optional TLS settings to enable secure communication between Pipelaner nodes.
- Requires both `certFile` and `keyFile` to be set.

---

## **Use Cases**

1. **HTTP/2 Communication**
    - Default communication mode for connecting Pipelaner nodes over HTTP/2.

2. **Unix Socket Communication**
    - Use Unix sockets for faster, local inter-node communication by setting `connectionType` to `"unix"` and specifying `unixSocketPath`.

3. **Secure Communication with TLS**
    - Enable TLS by setting the `tls` attribute in the `commonConfig`.

---

## **Notes**

- Ensure `Common.Pipelaner` is properly configured for shared connection settings.
- Always validate `unixSocketPath` when using Unix socket communication.
- Use TLS configuration for secure communication in production environments.

---

