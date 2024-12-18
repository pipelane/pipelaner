
# **Cmd**

The **Cmd** input component executes a specified command and streams its output.

---

## **Output Type**
- `String`

---

## **Pkl Configuration**

```pkl
class Cmd extends Input {
  fixed sourceName = "cmd"
  exec: Listing<String>
}
```

---

## **Parameters**

| **Name** | **Type**        | **Description**       | **Example**                                         |
|----------|-----------------|-----------------------|-----------------------------------------------------|
| `exec`   | `Listing<String>` | Command and arguments | `exec { "/usr/bin/log" "stream --style ndjson" }` |

---

## **Description**

The **Cmd** component runs the specified command with its arguments in a subprocess. The output of the command is then streamed as input for the pipeline.

### Example Configuration
```pkl
new Cmd {
  name = "cmd-example"
  exec = { "/usr/bin/log" "stream --style ndjson" }
}
```

### Use Case
- Capturing system logs and processing them in real-time.
- Integrating external tools or scripts into your pipeline.

---

