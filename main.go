package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-provider-camunda/camunda"
)

func main() {
	tfsdk.Serve(context.Background(), camunda.New, tfsdk.ServeOpts{
		Name: "camunda",
	})
}
