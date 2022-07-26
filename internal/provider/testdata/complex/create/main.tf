provider "fakelocal" {}

resource "fakelocal_complex_resource" "test" {
  string = "hello"

  list = [
    {
      integer = 0
    }
  ]

  object = {
    bool = true
  }

  set = [
    {
      integer = 0
    },
    {
      integer = 1
    }
  ]
}
