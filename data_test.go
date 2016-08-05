package gohttp

import (
	"encoding/json"

	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/gopkg.in/check.v1"
)

type DataTest struct{}

var _ = check.Suite(&DataTest{})

func (d *DataTest) TestPreparingJSON(c *check.C) {
	jsonToSerialize := map[string]interface{}{"test": "testing"}
	data, err := JSONData(jsonToSerialize)
	c.Assert(err, check.IsNil)
	c.Assert(data, check.NotNil)

	decoder := json.NewDecoder(data)
	var decodedJSON interface{}
	err = decoder.Decode(&decodedJSON)
	c.Assert(err, check.NotNil)
	c.Assert(jsonToSerialize, check.DeepEquals, decodedJSON)
}

func (d *DataTest) TestParsingJSON(c *check.C) {

}

func (d *DataTest) TestPreparingFormData(c *check.C) {
	data := map[string]interface{}{
		"test": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}
	formData, err := FormData(data)
	c.Assert(err, check.NotNil)
	c.Assert(formData, check.NotNil)
}

func (d *DataTest) TestParsingHeaders(c *check.C) {

}

func (d *DataTest) TestParsingFormData(c *check.C) {

}
