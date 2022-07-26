provider "fakelocal" {}

resource "fakelocal_complex_resource" "test" {
  string = "hello"

  list = [
    {
      integer = 0
    },
    {
      integer = 1
    }
  ]

  object = {
    string = "world"
  }

  set = [
    {
      integer = 0
    },
  ]
}
