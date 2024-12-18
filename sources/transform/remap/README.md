
# **Remap**

The **Remap** transform component uses an [expr](https://github.com/expr-lang/expr)-based engine to transform incoming messages by applying custom remapping logic.

---

## **Class Definition**

```pkl
class Remap extends Transform {
  fixed sourceName = "remap"
  code: String
}
```

---

## **Attributes**

| **Attribute** | **Type**   | **Description**                                     | **Default Value** |
|---------------|------------|-----------------------------------------------------|--------------------|
| `code`        | `String`   | The remapping rule written in the `expr` engine syntax. | **Required**      |

---

## **I/O Types**

- **Input Type:** `map[string]any`, `map[string][]any,string`, `[]byte(value)`, `[]byte`
- **Output Type:** `map[string]any`

---

## **Pkl Configuration Example**

### **Basic Remap Transform**
```pkl
new Transforms.Remap {
  name = "example-remap"
  code = "{ "value_name": Data.name, "value_price": Data.price }"
}
```

---

## **Description**

The **Remap** transform component allows you to transform or modify incoming messages based on custom rules defined in the `expr` language. This is particularly useful for altering the structure or content of messages before they are passed downstream.

---

### **Unit Test Example in Go**

The following test demonstrates remapping a message by extracting specific fields and creating a new structure:

```go
{
	name: "test expr maps return nil",
	args: args{
		val: map[string]any{
			"id":       1,
			"name":     "iPhone 12",
			"price":    999,
			"quantity": 1,
		},
		code: "{ "value_name": Data.name, "value_price": Data.price }",
	},
	want: map[string]any{
		"value_name":  "iPhone 12",
		"value_price": 999,
	},
},
```

---

### **Pkl Configuration Matching Unit Test**

The equivalent Pkl configuration for the test case above:

```pkl
new Transforms.Remap {
  name = "remap-extract-fields"
  code = "{ "value_name": Data.name, "value_price": Data.price }"
}
```

---

## **Use Cases**

1. **Data Transformation**
    - Modify fields in incoming messages, such as scaling values or changing formats.
2. **Data Enrichment**
    - Add new fields or compute derived values based on existing fields.

---

## **Notes**

- The `code` attribute must be a valid [expr](https://github.com/expr-lang/expr) rule.
- Ensure the rule logic aligns with the structure of incoming messages for accurate transformations.

---