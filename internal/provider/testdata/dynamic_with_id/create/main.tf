provider "fakelocal" {}

resource "fakelocal_dynamic_resource" "test" {
  id = "my_id"
  integer = 0
}
