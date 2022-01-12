---
page_title: "deployment Resource - terraform-provider-camunda"
subcategory: ""
description: |-
  The deployment resource allows you to configure a Camunda deployment.
---

# Resource `camunda_deployment`

The deployment resource allows you to configure a [Camunda deployment](https://docs.camunda.org/manual/7.16/reference/rest/deployment/).

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
