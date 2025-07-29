
# **Create custom generator/transform/sink in External Projects**

This guide provides detailed instructions on how to create custom generators/transforms/sinks.

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
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/Pipelaner.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/Components.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/sink/Sinks.pkl"
import "example.pkl"

pipelines {
  new Components.Pipeline {
    name = "example-pipeline"
    inputs {
      new example.ExampleGenInt {
        name = "example-gen-int"
        count = 10
      }
    }
    transforms {
      new example.ExampleMul {
        name = "example-mul"
        inputs {
          "example-gen-int"
        }
        mul = 2
      }
      new example.ExampleMul {
        threads = 10
        name = "example-mul2"
        inputs {
          "example-mul"
        }
        mul = 5
      }
    }
    sinks {
      new Sinks.Console {
        threads = 10
        name = "console"
        inputs {
          "example-mul2"
        }
      }
    }
  }
}

settings {
  gracefulShutdownDelay = 15.s
  logger {
    logLevel = "info"
  }
  healthCheck {
    enable = false
    host = "127.0.0.1"
    port = 8080
  }
  metrics {
    enable = false
    host = "127.0.0.1"
    port = 8082
  }
}
```

---

## 🛠 **Step 3: Implement Custom Components**

If you need custom components, create an implementation file (e.g., `pkl/custom.pkl`) with the following content:

```pkl
@go.Package {name = "gen/custom"}
module pipelaner.source.example

import "package://pkg.pkl-lang.org/pkl-go/pkl.golang@0.11.0#/go.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/input/Inputs.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/sink/Sinks.pkl"
import "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@1.0.0#/source/transform/Transforms.pkl"

class ExampleGenInt extends Inputs.Input {
  fixed sourceName = "example-generator"
  count: Int
}

class ExampleMul extends Transforms.Transform {
  fixed sourceName = "example-mul"
  mul: Int
}

class ExampleConsole extends Sinks.Sink {
  fixed sourceName = "example-console"
}
```

---

## 🔧 **Step 4: Generate Code**

If custom components were created, generate the required code using the following command:

```shell
pkl-gen-go pkl/custom.pkl
```

---

## 🚀 **Step 5: Implement and Register Components**

To use the custom components in your project, implement and register them in the source of **Pipelaner**. An example implementation can be found in [custom.go](https://github.com/pipelane/pipelaner/tree/main/example/custom/custom.go):

```go
source.RegisterInput("example-generator", &GenInt{})
source.RegisterTransform("example-mul", &TransMul{})
```

---

## 📜 **License**

This project is licensed under the [Apache 2.0](https://github.com/pipelane/pipelaner/blob/main/LICENSE) license.  
You are free to use, modify, and distribute the code under the terms of the license.
