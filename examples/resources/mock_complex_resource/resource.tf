resource "mock_complex_resource" "complex" {
  object = {
    bool = true

    object = {
      string = "nested_object"
    }
  }

  list_block {
    string = "one"
  }

  list_block {
    string = "two"
  }
}
