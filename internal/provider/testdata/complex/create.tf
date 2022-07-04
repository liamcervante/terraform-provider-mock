provider "fakelocal" {}

resource "complex_resource" "test" {
  provider = fakelocal
  bool = true
}
