package bot

import "testing"

func TestNewOptions(t *testing.T) {
	tests := [][]string{
		[]string{"name=name is", "hi, my name is GWoo", "name", "GWoo"},
		[]string{"start=start", "start 2016-09-28", "start", "2016-09-28"},
		[]string{"sql*=sql is", "sql is select * from users;", "sql", "select * from users;"},
		[]string{"query*=query:", "query: select * from users;", "query", "select * from users;"},
		[]string{"folder=folder is", "folder is /var/lib", "folder", "/var/lib"},
	}
	for _, test := range tests {
		op, err := NewOption(test[0], test[1])
		if err != nil {
			t.Error(err)
		}
		if op.Key != test[2] {
			t.Errorf("Expect key to be %s, was %s", test[2], op.Key)
		}
		if op.Value != test[3] {
			t.Errorf("Expect value to be %s, was %s", test[3], op.Value)
		}
	}
}
