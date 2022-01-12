terraform {
  required_providers {
    camunda = {
      source = "b-social/camunda"
    }
  }
}

provider "camunda" {
  endpoint = "http://localhost:8080/engine-rest"
  #  username = "demo"
  #  password = "demo"
}

resource "camunda_deployment" "pd1" {
  key       = "PD_KEY_1"
  resources = [
    {
      name    = "bpmn1.bpmn"
      content = file("files/bpmn1.bpmn")
    }
  ]
}

resource "camunda_deployment" "pd2" {
  key       = "PD_KEY_2"
  resources = [
    {
      name    = "bpmn2.bpmn"
      content = file("files/bpmn2.bpmn")
    }
  ]
}
