resource "jira_project" "my_project" {
  name = "Terraform test"
  key = "TFTEST"
  url = "http://example.org"
}
