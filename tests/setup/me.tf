# Who we are, used to make the approval policy tests approve on ourselves.

data "restful_resource" "me" {
  id = "/csp/gateway/am/api/loggedin/user"
}
