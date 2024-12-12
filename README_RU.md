
# **Pipelaner**

**Pipelaner** — это высокопроизводительный и эффективный **Framework и Агент** для создания data pipelines. Основой описания пайплайнов служит концепция **_Configuration As Code_** и язык конфигураций [**Pkl**](https://github.com/apple/pkl) от компании **Apple**.

Pipelaner управляет потоками данных с помощью трех ключевых сущностей: **Generator**, **Transform**, и **Sink**.

---

## 📖 **Содержание**
- [Основные сущности](#основные-сущности)
  - [Generator](#generator)
  - [Transform](#transform)
  - [Sink](#sink)
  - [Базовые параметры](#базовые-параметры)
- [Элементы пайплайнов "из коробки"](#элементы-пайплайнов-из-коробки)
  - [Generators](#generators)
  - [Transforms](#transforms)
  - [Sinks](#sinks)
- [Масштабируемость](#масштабируемость)
  - [Одноузловое развертывание](#одноузловое-развертывание)
  - [Многоузловое развертывание](#многоузловое-развертывание)
- [Поддержка](#поддержка)
- [Лицензия](#лицензия)

---

## 📌 **Основные сущности**

### **Generator**
Компонент, отвечающий за создание или получение исходных данных для потока. Generator может генерировать сообщения, события или извлекать данные из различных источников, таких как файлы, базы данных или API.

- **Пример использования:**  
  Чтение данных из файла или получение событий через вебхуки.

---

### **Transform**
Компонент, который обрабатывает данные в потоке. Transform выполняет операции, такие как фильтрация, агрегация, преобразование структуры или очистка данных, чтобы подготовить их для дальнейшей обработки.

- **Пример использования:**  
  Фильтрация записей с заданными условиями или преобразование формата данных из JSON в CSV.

---

### **Sink**
Конечная точка потока данных. Sink отправляет обработанные данные в целевую систему, например, в базу данных, API или систему очередей сообщений.

- **Пример использования:**  
  Сохранение данных в PostgreSQL или отправка их в Kafka-топик.

---

### **Базовые параметры**
| **Параметр**         | **Тип** | **Описание**                                                                                      |
|-----------------------|---------|---------------------------------------------------------------------------------------------------|
| `name`               | String  | Уникальное название элемента пайплайна.                                                          |
| `threads`            | Int     | Количество потоков для обработки сообщений. По умолчанию равно значению переменной `GOMAXPROC`.  |
| `outputBufferSize`   | Int     | Размер выходного буфера. **Не используется в Sink-компонентах.**                                  |

---

## 📦 **Элементы пайплайнов "из коробки"**

### **Generators**
| **Название**                                                                                 | **Описание**                                                                 |
|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|
| [**cmd**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/cmd)             | Считывает вывод команды, например `"/usr/bin/log" "stream --style ndjson"`.  |
| [**kafka**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/kafka)         | Консьюмер для Apache Kafka, передает по пайплайну значения `Value`.          |
| [**pipelaner**](https://github.com/pipelane/pipelaner/tree/main/sources/generator/pipelaner) | GRPC-сервер, передает значения через [gRPC](https://github.com/pipelane/pipelaner/tree/main/proto/service.proto). |

---

### **Transforms**
| **Название**                                                                                   | **Описание**                                                                |
|------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| [**batch**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/batch)           | Формирует пакеты данных заданного размера.                                  |
| [**chunks**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/chunks)         | Разбивает входящие данные на чанки.                                         |
| [**debounce**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/debounce)     | Устраняет "дребезг" (частые повторы) в данных.                              |
| [**filter**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/filter)         | Фильтрует данные по заданным условиям.                                      |
| [**remap**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/remap)           | Переназначает поля или преобразует структуру данных.                        |
| [**throttling**](https://github.com/pipelane/pipelaner/tree/main/sources/transform/throttling) | Ограничивает скорость обработки данных.                                     |

---

### **Sinks**
| **Название**                                                                                   | **Описание**                                                                |
|------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| [**clickhouse**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/clickhouse)      | Отправляет данные в базу данных ClickHouse.                                 |
| [**console**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/console)            | Выводит данные в консоль.                                                   |
| [**http**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/http)                  | Отправляет данные на указанный HTTP-эндпоинт.                               |
| [**kafka**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/kafka)                | Публикует данные в Apache Kafka.                                            |
| [**pipelaner**](https://github.com/pipelane/pipelaner/tree/main/sources/sink/pipelaner)        | Передает данные через gRPC к другим узлам Pipelaner.                        |

---

## 🌐 **Масштабируемость**

### **Одноузловое развертывание**
Для работы на одном хосте:  
![Одноузловая схема](https://github.com/pipelane/pipelaner/blob/c8e232106e9acf8a1d8682d225e369f282f6523a/images/pipelaner-singlehost.png/?raw=true "Одноузловая схема")

---

### **Многоузловое развертывание**
Для распределенной обработки данных на нескольких хостах:  
![Многоузловая схема](https://github.com/pipelane/pipelaner/blob/c8e232106e9acf8a1d8682d225e369f282f6523a/images/pipelaner-multihost.png/?raw=true "Многоузловая схема")

Для распределенного взаимодействия между узлами можно использовать:
1. **gRPC** — через генераторы и синки с параметром `sourceName: "pipelaner"`.
2. **Apache Kafka** — для чтения/записи данных через топики.

Пример конфигурации с использованием Kafka:
```pkl
new Inputs.Kafka {
    common {
        topics {
            "kafka-topic"
        }         
    }
}

new Sinks.Kafka {
    common {
        topics {
            "kafka-topic"
        }         
    }
}
```

---

## 🤝 **Поддержка**

Если у вас есть вопросы, предложения или вы нашли ошибку, пожалуйста, [создайте Issue](https://github.com/pipelane/pipelaner/issues/new) в репозитории.  
Вы также можете участвовать в обсуждениях проекта в разделе [Discussions](https://github.com/pipelane/pipelaner/discussions).

---

## 📜 **Лицензия**

Этот проект распространяется под лицензией [Apache 2.0](https://github.com/pipelane/pipelaner/blob/main/LICENSE).  
Вы можете свободно использовать, изменять и распространять код при соблюдении условий лицензии.
