resource "pools_string" "example" {
  resources = ["one", "two"]
  borrowers = ["alice", "bob"]
}
