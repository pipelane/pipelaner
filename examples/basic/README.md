
# **Using Pipelaner in External Projects**

This guide provides detailed instructions on how to integrate **Pkl** with **Pipelaner** in your external projects.

---

## 📂 **Step 1: Configure Dependencies**

1. Create a directory named `pkl` in your project.
2. Inside the `pkl` directory, create a file named `PklProject` with the following content:

```pkl
dependencies {
  ["pipelaner"] {
    uri = "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@x.x.x"
  }
}
```

Replace `x.x.x` with the required version of **Pipelaner**.

---

## ⚙️ **Step 2: Configure Pipelines**

Create a pipeline configuration file (e.g., `config.pkl`) with the following content:

```pkl
amends "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/Pipelaner.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/Components.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/sink/Sinks.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/input/Inputs.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/transform/Transforms.pkl"

pipelines {
  new Components.Pipeline {
    name = "example-pipeline"
    inputs {
      new Inputs.Cmd {
        name = "osx-logs"
        exec {
          "/usr/bin/log"
          "stream --style ndjson"
        }
      }
    }
    transforms {
      new Transforms.Chunk {
        name = "log-buffering"
        threads = 10
        outputBufferSize = 10_000
        inputs {
          "osx-logs"
        }
        maxChunkSize = 1_000
        maxIdleTime = 20.s
      }
    }
    sinks {
      new Sinks.Console {
        threads = 10
        name = "console-log"
        inputs {
          "log-buffering"
        }
      }
    }
  }
}

settings {
  gracefulShutdownDelay = 15.s
  logger {
    logLevel = "info"
    logFormat = "json"
  }
  healthCheck {
    host = "127.0.0.1"
    port = 8080
  }
  metrics {
    host = "127.0.0.1"
    port = 8082
  }
}
```

---

---

## 📜 **License**

This project is licensed under the [Apache 2.0](https://github.com/pipelane/pipelaner/blob/main/LICENSE) license.  
You are free to use, modify, and distribute the code under the terms of the license.
