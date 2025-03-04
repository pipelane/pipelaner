
# **Clickhouse**

The **Clickhouse** sink component writes incoming messages to a Clickhouse database table.

---

## **Config Definition**

```pkl
class ChCredentials {
   address: String(!isEmpty)
   user: String(!isEmpty)
   password: String(!isEmpty)
   database: String(!isEmpty)
}

class Clickhouse extends Sink {
   fixed sourceName = "clickhouse"
   credentials: Common.ChCredentials
   tableName: String(!isEmpty)
   asyncInsert: String = "1"
   waitForAsyncInsert: String = "1"
   maxPartitionsPerInsertBlock: Int = 1000
}
```

---

## **Attributes**

| **Attribute**                        | **Type**             | **Description**                                                              | **Default Value** |
|---------------------------------------|----------------------|------------------------------------------------------------------------------|--------------------|
| `credentials`                         | `Common.ChCredentials` | Reusable credentials for Clickhouse connection.                             | **Required**      |
| `tableName`                           | `String`            | The target table to write to within the database.                            | **Required**      |
| `asyncInsert`                         | `String`            | Enables asynchronous insertion (1 for enabled, 0 for disabled).              | `"1"`             |
| `waitForAsyncInsert`                  | `String`            | Specifies whether to wait for asynchronous insert completion (1 for yes, 0 for no). | `"1"`             |
| `maxPartitionsPerInsertBlock`         | `Int`               | The maximum number of partitions per insert block.                          | `1000`            |

---

## **I/O Types**

- **Input Type:** `map[string]any`, where `key` is a column name, and `value` is the column value.
- **Output:** Writes messages to the specified Clickhouse database and table.

---

## **Pkl Configuration Example**

### **Basic Clickhouse Sink**
```pkl
new Clickhouse {
  name = "example-clickhouse"
  credentials = new Common.ChCredentials {
    address = "http://127.0.0.1:8123"
    user = "default"
    password = "password"
    database = "example_db"
  }
  tableName = "example_table"
}
```

### **Clickhouse Sink with Async Insert**
```pkl
new Clickhouse {
  name = "example-clickhouse-async"
  credentials = new Common.ChCredentials {
    address = "http://127.0.0.1:8123"
    user = "default"
    password = "password"
    database = "example_db"
  }
  tableName = "example_table"
  asyncInsert = "1"
  waitForAsyncInsert = "0"
  maxPartitionsPerInsertBlock = 500
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
  credentials = new Common.ChCredentials {
    address = "http://127.0.0.1:8123"
    user = "default"
    password = "password"
    database = "example_db"
  }
  tableName = "example_table"
  asyncInsert = "1"
  waitForAsyncInsert = "0"
  maxPartitionsPerInsertBlock = 500
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

- Ensure that the `credentials` attribute is correctly configured with valid Clickhouse connection parameters.
- Asynchronous insertion improves performance but requires `waitForAsyncInsert` to be configured appropriately based on use case.
- Use **Transform Chunks** for optimal performance by aggregating messages into large batches.
- Adjust `maxPartitionsPerInsertBlock` to optimize insert performance for large datasets.

---

If you have additional questions or need further clarification, feel free to reach out!
