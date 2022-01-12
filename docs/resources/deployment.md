---
page_title: "deployment Resource - terraform-provider-camunda"
subcategory: ""
description: |-
  The process definition resource allows you to configure a Camunda deployment.
---

# Resource `camunda_deployment`

The order resource allows you to configure a Camunda process definition.

## Example Usage

```terraform
resource "camunda_deployment" "pd1" {
  key       = "PD_KEY_1"
  resources = [
    {
      name    = "bpmn1.bpmn"
      content = file("files/bpmn1.bpmn")
    }
  ]
}
```

## Argument Reference

- `key` - (Required) Key to identify a deployment.
- `resources` - (Required) A list of deployable files with a name and content.

## Attribute Reference

- `id` - id of a deployment.
