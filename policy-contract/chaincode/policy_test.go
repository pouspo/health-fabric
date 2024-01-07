package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccess_add(t *testing.T) {
	type fields struct {
		Read  []string
		Write []string
	}
	type args struct {
		mode   string
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Access
	}{
		{
			name: "add_read",
			fields: fields{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
			args: args{
				mode:   "read",
				fields: []string{"field_2", "field3"},
			},
			want: Access{
				Read:  []string{"field_1", "field_2", "field3"},
				Write: []string{"field_1"},
			},
		},
		{
			name: "add_write",
			fields: fields{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
			args: args{
				mode:   "write",
				fields: []string{"field_2", "field3"},
			},
			want: Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1", "field_2", "field3"},
			},
		},
		{
			name: "add_invalid",
			fields: fields{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
			args: args{
				mode:   "invalid",
				fields: []string{"field_2", "field3"},
			},
			want: Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Access{
				Read:  tt.fields.Read,
				Write: tt.fields.Write,
			}
			a.add(tt.args.mode, tt.args.fields)

			if !cmp.Equal(a, tt.want) {
				t.Errorf("add() = %v, want %v", a, tt.want)
			}
		})
	}
}

func TestAccess_remove(t *testing.T) {
	type fields struct {
		Read  []string
		Write []string
	}
	type args struct {
		mode   string
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Access
	}{
		{
			name: "remove_write",
			fields: fields{
				Read:  []string{"field_1"},
				Write: []string{"field_1", "field_2", "field_3", "field_3"},
			},
			args: args{
				mode:   "write",
				fields: []string{"field_2", "field_3"},
			},
			want: Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
		},
		{
			name: "remove_read",
			fields: fields{
				Read:  []string{"field_1", "field_2", "field3"},
				Write: []string{"field_1"},
			},
			args: args{
				mode:   "read",
				fields: []string{"field_2", "field3"},
			},
			want: Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
		},
		{
			name: "add_invalid",
			fields: fields{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
			args: args{
				mode:   "invalid",
				fields: []string{"field_2", "field3"},
			},
			want: Access{
				Read:  []string{"field_1"},
				Write: []string{"field_1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Access{
				Read:  tt.fields.Read,
				Write: tt.fields.Write,
			}
			a.remove(tt.args.mode, tt.args.fields)

			if !cmp.Equal(a, tt.want) {
				t.Errorf("remove() = %v, want %v", a, tt.want)
			}
		})
	}
}

func TestSomething(t *testing.T) {
	var policy Policy
	err := json.Unmarshal([]byte(`{"id":"admin","policy_map":{"user_101":{"read":["field_1"],"write":null}}}`), &policy)
	require.NoError(t, err)
	fmt.Print(policy)

}
