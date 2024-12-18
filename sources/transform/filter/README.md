
# **Filter**

The **Filter** transform component evaluates incoming messages against a custom rule (written in an [expr](https://github.com/expr-lang/expr) rule-based engine) and processes only those that match the rule.

---

## **Class Definition**

```pkl
class Filter extends Transform {
  fixed sourceName = "filter"
  code: String
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                       | **Default Value** |
|---------------|------------|-------------------------------------------------------|--------------------|
| `code`        | `String`   | The filtering rule written in the `expr` engine syntax. | **Required**      |

---

## **I/O Types**

- **Input Type:** `map[string]any`, `string`, `[]byte`
- **Output Type:** `map[string]any`, `string`, `[]byte`

---

## **Pkl Configuration Example**

### **Basic Filter Transform**
```pkl
new Transforms.Filter {
  name = "example-filter"
  code = "Data.count > 5"
}
```

---

## **Description**

The **Filter** transform component uses a custom rule defined in the `expr` engine to evaluate each incoming message. Only messages that satisfy the rule are passed to the next component.

---

### **Example from Unit Test**

The following test demonstrates filtering messages where the `count` field in a JSON object is greater than 5:

```go
{
	name: "test filtering string return 10",
	args: args{
		code: "Data.count > 5",
		val:  []byte("{"count":10}"),
	},
	want: []byte("{"count":10}"),
},
```

---

### **Pkl Configuration Matching with Unit Test**

The equivalent Pkl configuration for the test case above:

```pkl
new Transforms.Filter {
  name = "filter-count-test"
  code = "Data.count > 5"
}
```

---

## **Use Cases**

1. **Message Filtering**
    - Filter out messages that do not meet specific criteria.
2. **Data Validation**
    - Process only messages with valid or desired data.

---

## **Notes**

- The `code` attribute must be a valid [expr](https://github.com/expr-lang/expr] rule).
- Ensure the rule logic matches the expected structure of incoming messages.

---
