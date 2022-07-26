provider "fakelocal" {}

resource "fakelocal_simple_resource" "test" {
  id = "my_id"
  integer = 0
}
