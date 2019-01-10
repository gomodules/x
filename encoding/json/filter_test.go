package json_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/appscode/go/encoding/json"
	"github.com/stretchr/testify/assert"
)

var (
	obj = `
{
  "apiVersion": "extensions/v1beta1",
  "kind": "DaemonSet",
  "metadata": {
    "name": "busy-dm",
    "namespace": "default",
    "labels": {
      "app": "busy-dm",
      "release": "latest"
    }
  },
  "spec": {
    "template": {
      "metadata": {
        "labels": {
          "name": "busy-dm"
        }
      },
      "spec": {
        "nodeSelector": {
          "kubernetes.io/hostname": "ip-172-20-53-35.ec2.internal"
        },
        "containers": [
          {
            "image": "busybox",
            "command": [
              "sleep",
              "3600"
            ],
            "imagePullPolicy": "IfNotPresent",
            "name": "busybox"
          },
          {
            "image": "nginx",
            "command": [
              "sleep",
              "3600"
            ],
            "imagePullPolicy": "IfNotPresent",
            "name": "nginx"
          }
        ]
      }
    }
  }
}`
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name   string
		obj    string
		filter string
		result string
		rErr error
	}{
		{
			name: `non-nested filter`,
			obj:  obj,
			filter: `
{
  "apiVersion": null,
  "kind": null
}
`,
			result: `
{
  "apiVersion": "extensions/v1beta1",
  "kind": "DaemonSet"
}
`,
		},
		{
			name: `nested object filter`,
			obj:  obj,
			filter: `
{
  "apiVersion": null,
  "kind": null,
  "metadata": {
	"name": null,
    "labels": {
      "app2": null,
      "release": null
    }
  }
}
`,
			result: `
{
  "apiVersion": "extensions/v1beta1",
  "kind": "DaemonSet",
  "metadata": {
    "labels": {
      "release": "latest"
    },
    "name": "busy-dm"
  }
}
`,
		},
		{
			name: `nested array filter`,
			obj:  obj,
			filter: `
{
  "apiVersion": null,
  "kind": null,
  "metadata": {
	"name": null,
	"namespace": null,
    "labels": {
      "app2": null
    }
  },
  "spec": {
    "template": {
      "spec": {
        "containers": {
          "name": null,
          "image": null
        }
      }
    }
  }
}
`,
			result: `
{
  "apiVersion": "extensions/v1beta1",
  "kind": "DaemonSet",
  "metadata": {
    "labels": {},
    "name": "busy-dm",
    "namespace": "default"
  },
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "image": "busybox",
            "name": "busybox"
          },
          {
            "image": "nginx",
            "name": "nginx"
          }
        ]
      }
    }
  }
}
`,
		},
		{
			name:   `should fail to filter non-nested key`,
			obj:    obj,
			filter: `
{
  "apiVersion": null,
  "kind": null,
  "metadata": {
	"name": {
      "app2": null
    }
  }
}
`,
			rErr: errors.New(`can't apply filter {"app2":null} on metadata.name: busy-dm`),
		},
		{
			name:   `should fail to filter non-nested array element`,
			obj:    obj,
			filter: `
{
  "apiVersion": null,
  "kind": null,
  "metadata": {
	"name": null,
	"namespace": null
  },
  "spec": {
    "template": {
      "spec": {
        "containers": {
          "name": {
            "app2": null
          }
        }
      }
    }
  }
}`,
			rErr: errors.New(`can't apply filter {"app2":null} on spec.template.spec.containers[0].name: busybox`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(test.obj), &obj)
			if err != nil {
				t.Fatal(err)
			}

			var filter map[string]interface{}
			err = json.Unmarshal([]byte(test.filter), &filter)
			if err != nil {
				t.Fatal(err)
			}

			out, err := Filter(obj, filter)
			if err != nil {
				if test.rErr != nil {
					assert.EqualError(t, err, test.rErr.Error())
					return
				}
				t.Fatal(err)
			}
			outBytes, err := json.Marshal(out)
			if err != nil {
				t.Fatal(err)
			}
			assert.JSONEq(t, test.result, string(outBytes))
		})
	}
}
