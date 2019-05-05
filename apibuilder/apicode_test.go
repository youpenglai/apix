package apibuilder

import "testing"

type testReader struct {
	v map[string]map[string]interface{}
}

func (tr *testReader) Get(name, from string) (v interface{}) {
	fromData , _ := tr.v[from]
	v, _ = fromData[name]
	return
}

func newTestReader() *testReader {
	t := make(map[string]map[string]interface{})
	t["body"] = map[string]interface{}{
		"userName": "",
		"password": "123456",
	}
	return &testReader{v:t}
}

func TestGenApiCode(t *testing.T) {
	apiDoc := NewApiDoc()
	err := apiDoc.Parse([]byte(testApiDoc))
	if err != nil {
		t.Error(err)
		return
	}
	code, err := GenApiCode(apiDoc)
	if err != nil {
		t.Error(err)
		return
	}

	reader := newTestReader()
	codeBlock, _ := code.GetApiCode("/auth/login")
	codeBlock.BindParamReader(reader)
	params, err := codeBlock.ReadParams()
	if err != nil {
		t.Error(err)
		return
	}
	paramJson, err := params.MarshalJSON()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(paramJson))
	err = params.Validation()
	t.Log("error:", err.Error())
}
