package cac

import (
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestNewCacTable(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *cacTable
	}{
		{
			name: "test_new",
			args: args{
				name: "test",
			},
			want: &cacTable{
				name:      "test",
				key2Cache: make(map[string]*cac),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCacTable(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCacTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacTable_Add(t *testing.T) {

	ct := NewCacTable("add")

	type args struct {
		key    string
		val    interface{}
		expire time.Duration
		cbs    []func()
	}
	tests := []struct {
		name string
		c    *cacTable
		args args
	}{
		{
			name: "test_add",
			c:    ct,
			args: args{
				key:    "key",
				val:    "val",
				expire: time.Second * 2,
				cbs: []func(){func() {
					log.Println("call back doing")
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Add(tt.args.key, tt.args.val, tt.args.expire, tt.args.cbs...)
		})
	}

	value, find := ct.Get("key")
	assert.Equal(t, value, "val")
	assert.True(t, find, true)

	time.Sleep(time.Second * 3)
	value_, find_ := ct.Get("key")
	assert.Nil(t, value_)
	assert.False(t, find_)

}

func Test_cacTable_Delete(t *testing.T) {
	ct := NewCacTable("delete")
	ct.Add("key", "delete", time.Second*10)

	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *cacTable
		args args
	}{
		{
			name: "test_delete",
			c:    ct,
			args: args{
				key: "key",
			},
		},
	}

	value, find := ct.Get("key")
	assert.Equal(t, value, "delete")
	assert.True(t, find, true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Delete(tt.args.key)
		})
	}
	value_, find_ := ct.Get("key")
	assert.Nil(t, value_)
	assert.False(t, find_)
}

func Test_cacTable_Get(t *testing.T) {

	ct := NewCacTable("get")
	ct.Add("k", "v", time.Second*2)

	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *cacTable
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name: "test_get",
			c:    ct,
			args: args{
				key: "k",
			},
			want:  "v",
			want1: true,
		},
		{
			name: "test_get_expire",
			c:    ct,
			args: args{
				key: "k",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cacTable.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cacTable.Get() got1 = %v, want %v", got1, tt.want1)
			}
			time.Sleep(time.Second * 2)
		})
	}
}

func Test_cacTable_Reset(t *testing.T) {
	ct := NewCacTable("reset")
	ct.Add("k", "v", time.Second)

	type args struct {
		key    string
		expire time.Duration
	}
	tests := []struct {
		name string
		c    *cacTable
		args args
	}{
		{
			name: "test_reset",
			c:    ct,
			args: args{
				key:    "k",
				expire: time.Second * 3,
			},
		},
	}

	value, find := ct.Get("k")
	assert.True(t, find)
	assert.Equal(t, value, "v")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := tt.c.Reset(tt.args.key, tt.args.expire)
			assert.True(t, success)

		})
	}

	time.Sleep(time.Second * 2)
	value_, find_ := ct.Get("k")
	assert.True(t, find_)
	assert.Equal(t, value_, "v")
}
