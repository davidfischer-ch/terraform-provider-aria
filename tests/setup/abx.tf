# ABX action referenced by the subscription tests as a runnable (TF_VAR_test_abx_action_id).

resource "aria_abx_action" "test" {
  name            = "${local.prefix}_ACTION"
  description     = "Action referenced by Aria provider's acceptance tests, does nothing."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  shared          = true
  project_id      = vra_project.test[0].id

  inputs = {
    SomeString = jsonencode("")
    SomeNumber = jsonencode(42)
  }

  source = <<-EOT
    def handler(*args, **kwargs):
        print('Called with', args, kwargs)
  EOT
}
