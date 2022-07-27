provider "fakelocal" {}

resource "fakelocal_dynamic_resource" "test" {
  integer = 0

  nested_list {
    integer = 0
  }
}
