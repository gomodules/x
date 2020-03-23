/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRemoveNestedField(t *testing.T) {
	obj := yaml.MapSlice{
		{
			Key: "x",
			Value: yaml.MapSlice{
				{Key: "y", Value: 1},
				{Key: "a", Value: "foo"},
			},
		},
	}

	RemoveNestedField(&obj, "x", "a")
	assert.Len(t, _entry(obj, "x"), 1)
	RemoveNestedField(&obj, "x", "y")
	assert.Empty(t, _entry(obj, "x"))
	RemoveNestedField(&obj, "x")
	assert.Empty(t, obj)
	RemoveNestedField(&obj, "x") // Remove of a non-existent field
	assert.Empty(t, obj)
}

func TestNestedFieldNoCopy(t *testing.T) {
	target := yaml.MapSlice{
		{Key: "foo", Value: "bar"},
	}

	obj := yaml.MapSlice{
		{
			Key: "a",
			Value: yaml.MapSlice{
				{Key: "b", Value: target},
				{Key: "c", Value: nil},
				{Key: "d", Value: []interface{}{"foo"}},
				{Key: "e", Value: yaml.MapSlice{
					{Key: "f", Value: "bar"},
				}},
			},
		},
	}

	// case 1: field exists and is non-nil
	res, exists, err := NestedFieldNoCopy(obj, "a", "b")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, target, res)
	_set(&target, "foo", "baz")
	assert.Equal(t, _entry(target, "foo"), _entry(res.(yaml.MapSlice), "foo"), "result should be a reference to the expected item")

	// case 2: field exists and is nil
	res, exists, err = NestedFieldNoCopy(obj, "a", "c")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 3: error traversing obj
	res, exists, err = NestedFieldNoCopy(obj, "a", "d", "foo")
	assert.False(t, exists)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// case 4: field does not exist
	res, exists, err = NestedFieldNoCopy(obj, "a", "g")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 5: intermediate field does not exist
	res, exists, err = NestedFieldNoCopy(obj, "a", "g", "f")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 6: intermediate field is null
	//         (background: happens easily in YAML)
	res, exists, err = NestedFieldNoCopy(obj, "a", "c", "f")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 7: array/slice syntax is not supported
	//         (background: users may expect this to be supported)
	res, exists, err = NestedFieldNoCopy(obj, "a", "e[0]")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestNestedFieldCopy(t *testing.T) {
	target := yaml.MapSlice{
		{Key: "foo", Value: "bar"},
	}

	obj := yaml.MapSlice{
		{
			Key: "a",
			Value: yaml.MapSlice{
				{Key: "b", Value: target},
				{Key: "c", Value: nil},
				{Key: "d", Value: []interface{}{"foo"}},
			},
		},
	}

	// case 1: field exists and is non-nil
	res, exists, err := NestedFieldCopy(obj, "a", "b")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Equal(t, target, res)
	_set(&target, "foo", "baz")
	assert.NotEqual(t, _entry(target, "foo"), _entry(res.(yaml.MapSlice), "foo"), "result should be a copy of the expected item")

	// case 2: field exists and is nil
	res, exists, err = NestedFieldCopy(obj, "a", "c")
	assert.True(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)

	// case 3: error traversing obj
	res, exists, err = NestedFieldCopy(obj, "a", "d", "foo")
	assert.False(t, exists)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	// case 4: field does not exist
	res, exists, err = NestedFieldCopy(obj, "a", "e")
	assert.False(t, exists)
	assert.Nil(t, err)
	assert.Nil(t, res)
}
