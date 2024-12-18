
# **Pipelaner**

**Pipelaner** is a high-performance and efficient **Framework and Agent** for creating data pipelines. The core of pipeline descriptions is based on the **_Configuration As Code_** concept and the [**Pkl**](https://github.com/apple/pkl) configuration language by **Apple**.

Pipelaner manages data streams through three key entities: **Generator**, **Transform** and **Sink**.

---

## 📖 **Contents**
- [Core Entities](#core-entities)
  - [Generator](#generator)
  - [Transform](#transform)
  - [Sink](#sink)
  - [Basic Parameters](#basic-parameters)
- [Built-in Pipeline Elements](#built-in-pipeline-elements)
  - [Generators](#generators)
  - [Transforms](#transforms)
  - [Sinks](#sinks)
- [Scalability](#scalability)
  - [Single-Node Deployment](#single-node-deployment)
  - [Multi-Node Deployment](#multi-node-deployment)
- [Examples](#examples)
- [Support](#support)
- [License](#license)

---

## 📌 **Core Entities**

### **Generator**
The component responsible for creating or retrieving source data for the pipeline. Generators can produce messages, events, or retrieve data from various sources such as files, databases, or APIs.

- **Example use case:**  
  Reading data from a file or receiving events via webhooks.

---

### **Transform**
The component that processes data within the pipeline. Transforms perform operations such as filtering, aggregation, data transformation, or cleaning to prepare it for further processing.

- **Example use case:**  
  Filtering records based on specific conditions or converting data format from JSON to CSV.

---

### **Sink**
The final destination for the data stream. Sinks send processed data to a target system, such as a database, API, or message queue.

- **Example use case:**  
  Saving data to PostgreSQL or sending it to a Kafka topic.

---

### **Basic Parameters**
| **Parameter**         | **Type** | **Description**                                                                                  |
|-----------------------|---------|---------------------------------------------------------------------------------------------------|
| `name`               | String  | Unique name of the pipeline element.                                                             |
| `threads`            | Int     | Number of threads for processing messages. Defaults to the value of `GOMAXPROC`.                |
| `outputBufferSize`   | Int     | Size of the output buffer. **Not applicable to Sink components.**                                |

---

## 📦 **Built-in Pipeline Elements**

### **Generators**
| **Name**                                                                                 | **Description**                                                             |
|------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| [**cmd**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/cmd)         | Reads the output of a command, e.g., `"/usr/bin/log" "stream --style ndjson"`. |
| [**kafka**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/kafka)     | Apache Kafka consumer that streams `Value` into the pipeline.              |
| [**pipelaner**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/pipelaner) | GRPC server that streams values via [gRPC](https://github.com/pipelane/pipelaner/tree/main/proto/service.proto). |

---

### **Transforms**
| **Name**                                                                                   | **Description**                                                            |
|--------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| [**batch**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/batch)       | Forms batches of data with a specified size.                               |
| [**chunks**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/chunks)     | Splits incoming data into chunks.                                          |
| [**debounce**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/debounce) | Eliminates "bounce" (frequent repeats) in data.                            |
| [**filter**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/filter)     | Filters data based on specified conditions.                                |
| [**remap**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/remap)       | Reassigns fields or transforms the data structure.                         |
| [**throttling**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/throttling) | Limits data processing rate.                                              |

---

### **Sinks**
| **Name**                                                                                   | **Description**                                                            |
|--------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| [**clickhouse**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/clickhouse)  | Sends data to a ClickHouse database.                                       |
| [**console**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/console)        | Outputs data to the console.                                               |
| [**http**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/http)              | Sends data to a specified HTTP endpoint.                                   |
| [**kafka**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/kafka)            | Publishes data to Apache Kafka.                                            |
| [**pipelaner**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/pipelaner)    | Streams data via [gRPC](https://github.com/pipelane/pipelaner/tree/main/proto/service.proto) to other Pipelaner nodes.                            |

---

## 🌐 **Scalability**

### **Single-Node Deployment**
For operation on a single host:  
![Single Node](https://github.com/pipelane/pipelaner/blob/c8e232106e9acf8a1d8682d225e369f282f6523a/images/pipelaner-singlehost.png/?raw=true "Single Node Deployment")

---

### **Multi-Node Deployment**
For distributed data processing across multiple hosts:  
![Multi-Node](https://github.com/pipelane/pipelaner/blob/c8e232106e9acf8a1d8682d225e369f282f6523a/images/pipelaner-multihost.png/?raw=true "Multi-Node Deployment")

For distributed interaction between nodes, you can use:
1. **gRPC** — via generators and sinks with the parameter `sourceName: "pipelaner"`.
2. **Apache Kafka** — for reading/writing data via topics.

Example configuration using Kafka:
```pkl
new Inputs.Kafka {
    ...
    common {
        ...
        topics {
            "kafka-topic"
        }         
    }
}

new Sinks.Kafka {
    ...
    common {
        ...
        topics {
            "kafka-topic"
        }         
    }
}
```

---
## 🚀 **Examples**

| **Examples**                                                                  | **Description**                                         |
|-------------------------------------------------------------------------------|---------------------------------------------------------|
| [**Basic Pipeline**](https://github.com/pipelane/pipelaner/tree/main/examples/basic)   | A simple example illustrating the creation of a basic pipeline with prebuilt components.                    |
| [**Custom Components**](https://github.com/pipelane/pipelaner/tree/main/examples/custom) | An advanced example showing how to create and integrate custom Generators, Transforms, and Sinks. |
---

### **Overview**

1. **🌟 Basic Pipeline**  
   Learn the fundamentals of creating a pipeline with minimal configuration using ready-to-use components.

2. **🛠 Custom Components**  
   Extend **Pipelaner**’s functionality by developing your own Generators, Transforms, and Sinks.
---

Each example includes clear configuration files and explanations to help you get started quickly.

💡 **Tip**: Use these examples as templates to customize and build your own pipelines efficiently.


## 🤝 **Support**

If you have questions, suggestions, or encounter any issues, please [create an Issue](https://github.com/pipelane/pipelaner/issues/new) in the repository.  
You can also participate in discussions in the [Discussions](https://github.com/pipelane/pipelaner/discussions) section.

---

## 📜 **License**

This project is licensed under the [Apache 2.0](https://github.com/pipelane/pipelaner/blob/main/LICENSE) license.  
You are free to use, modify, and distribute the code under the terms of the license.
