
# **Clickhouse**

The **Clickhouse** sink component writes incoming messages to a Clickhouse database table.

---

## **Config Definition**

```pkl
class Clickhouse extends Sink {
  fixed sourceName = "clickhouse"
  address: String
  user: String
  password: String
  database: String
  tableName: String

  asyncInsert: String = "1"
  waitForAsyncInsert: String = "1"
}
```

---

## **Attributes**

| **Attribute**        | **Type**   | **Description**                                                              | **Default Value** |
|-----------------------|------------|------------------------------------------------------------------------------|--------------------|
| `address`            | `String`   | The address of the Clickhouse server.                                        | **Required**      |
| `user`               | `String`   | The username for authentication.                                             | **Required**      |
| `password`           | `String`   | The password for authentication.                                             | **Required**      |
| `database`           | `String`   | The target database to write to.                                             | **Required**      |
| `tableName`          | `String`   | The target table to write to within the database.                            | **Required**      |
| `asyncInsert`        | `String`   | Enables asynchronous insertion (1 for enabled, 0 for disabled).              | `"1"`             |
| `waitForAsyncInsert` | `String`   | Specifies whether to wait for asynchronous insert completion (1 for yes, 0 for no). | `"1"`             |

---

## **I/O Types**

- **Input Type:** `map[string]any`, where `key` is a name of column, `value` is a value of column.
- **Output:** Writes messages to the specified Clickhouse database and table.

---

## **Pkl Configuration Example**

### **Basic Clickhouse Sink**
```pkl
new Clickhouse {
  name = "example-clickhouse"
  address = "http://127.0.0.1:8123"
  user = "default"
  password = "password"
  database = "example_db"
  tableName = "example_table"
}
```

### **Clickhouse Sink with Async Insert**
```pkl
new Clickhouse {
  name = "example-clickhouse-async"
  address = "http://127.0.0.1:8123"
  user = "default"
  password = "password"
  database = "example_db"
  tableName = "example_table"
  asyncInsert = "1"
  waitForAsyncInsert = "0"
}
```

### **Recommended Configuration with Transform Chunks**
```pkl
new Transforms.Chunks {
  name = "example-chunks"
  maxChunkSize = 1000
  maxIdleTime = 2.s
  inputs {
    ...
  }
}

new Sinks.Clickhouse {
  name = "example-clickhouse-batched"
  address = "http://127.0.0.1:8123"
  user = "default"
  password = "password"
  database = "example_db"
  tableName = "example_table"
  asyncInsert = "1"
  waitForAsyncInsert = "0"
  inputs {
    "example-chunks"
  }
}
```

---

## **Description**

The **Clickhouse** sink allows efficient insertion of pipeline messages into a Clickhouse database table. It supports both synchronous and asynchronous insertion modes.

### **Best Practices**
- **Batching Messages:** Writing large batches of data to Clickhouse is more efficient. It is recommended to use the **Transform Chunks** component before the Clickhouse sink to aggregate messages into larger batches.
- Example: Use a chunk size of 1000 messages and a maximum idle time of 2 seconds.

---

### **Use Cases**

1. **Data Storage**
    - Persist pipeline data to a Clickhouse table for further analysis.
2. **High Performance Insertion**
    - Use asynchronous insert mode to optimize throughput.
3. **Batch Optimization**
    - Combine with the `Chunks` transform to batch messages for efficient database insertion.

---

## **Notes**

- Ensure that the `address`, `user`, `password`, `database`, and `tableName` attributes are correctly configured.
- Asynchronous insertion improves performance but requires `waitForAsyncInsert` to be configured appropriately based on use case.
- Use **Transform Chunks** for optimal performance by aggregating messages into large batches.

---
