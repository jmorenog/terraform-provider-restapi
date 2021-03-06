package restapi

import (
  "os"
  "testing"
  "github.com/hashicorp/terraform/helper/resource"
  "github.com/Mastercard/terraform-provider-restapi/fakeserver"
)

func TestAccRestApiObject_importBasic(t *testing.T) {
  debug := false
  api_server_objects := make(map[string]map[string]interface{})

  svr := fakeserver.NewFakeServer(8082, api_server_objects, true, debug, "")
  os.Setenv("REST_API_URI", "http://127.0.0.1:8082")

  client, err := NewAPIClient("http://127.0.0.1:8082/", false, "", "", make(map[string]string, 0), 2, "id", make([]string, 0), false, false, debug)
  if err != nil { t.Fatal(err) }
  client.send_request("POST", "/api/objects", `{ "id": "1234", "first": "Foo", "last": "Bar" }`)

  resource.UnitTest(t, resource.TestCase{
    Providers:    testAccProviders,
    PreCheck:     func() { svr.StartInBackground() },
    Steps: []resource.TestStep{
      {
        Config: generate_test_resource(
          "Foo",
          `{ "id": "1234", "first": "Foo", "last": "Bar" }`,
          make(map[string]interface{}),
        ),
      },
      {
        ResourceName: "restapi_object.Foo",
        ImportState: true,
        ImportStateId: "1234",
        ImportStateIdPrefix: "/api/objects/",
        ImportStateVerify: true,
        ImportStateVerifyIgnore: []string{ "debug", "data" },
      },
    },
  })

  svr.Shutdown()
}
