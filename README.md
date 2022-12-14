** This has been replaced by [hashicorp/terraform-provider-tfcoremock](github.com/hashicorp/terraform-provider-tfcoremock), please use that repository/provider **

# terraform-provider-mock

A Terraform provider that can model any Terraform configuration, will write the state to disk, and can be used to reason about state changes and test difficult to reproduce edge cases in terraform.

## Static Resources

This provider contains two resources:

- `mock_simple_resource`
- `mock_complex_resource`

Between them, they can be used to model almost any situation you would encounter with resources in any other provider.

### Examples

#### mock_simple_resource

```hcl
resource "mock_simple_resource" "simple_resource" {
  bool = true
  integer = 0
  float = 1.0
  string = "my string"
}
```

#### mock_complex_resource

```hcl

resource "mock_complex_resource" "complex_resource" {
  integer = 0
  
  object = {
    integer = 0
    object = {
      string = "nested_object"
    }
  }
  
  map = {
    "key1": {
      integer = 1
    },
    "key2": {
      integer = 2
    }
  }
  
  nested_list {
    integer = 0
  }
  nested_list {
    integer = 1
  }
}

```

## Dynamic Resources

If, for any reason, the static resources do not provide enough cover for a particular edge case, then you can define dynamic resources.

The provider looks for a file called `dynamic_resources.json` in the same directory as the terraform definitions, and will add resources based on the definitions in this file.

The provider automatically installs an `id` attribute into every dynamic resource.

### Example

The following example defines a resource called `mock_dynamic_resource` with a single `integer` attribute called `my_attribute`.

#### dynamic_resources.json

```json
{
  "mock_dynamic_resource": {
    "attributes": {
      "my_attribute": {
        "type": "integer",
        "optional": true
      }
    }
  }
}
```

#### main.tf

```hcl
provider "mock" {}

resource "mock_dynamic_resource" {
  my_attribute = 0
}
```

## Data Sources

Every resource defined as part of the mock provider (including the simple, complex, and all the dynamic resources) are duplicated and provided as data sources. 
So, any resource you can build you can also model as a data source and retrieve for reference as a `data` block.

## Terraform State

In order to model and test for things like changes in terraform state vs. the reality, this provider writes all the resources to disk.
By default, resources are written into the `terraform.resource` directory and data sources are read from the `terraform.data` directory.

Every resource and data source has an ID attribute.
This can be specified as part of the resource, and must be specified as part of the data source (as it is required to look up the data source).

With this in mind, you can populate `terraform.data` with anything you need, and then running `terraform apply` would write all the created resources in the resource directory.
To test drift between reality and terraform state, you could then just go and edit the resources in the resource directory.
