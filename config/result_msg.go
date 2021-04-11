package config

import "encoding/json"

//http返回消息格式
type ResultMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func (r ResultMsg) ToMap() map[string]interface{}{
	j, _ := json.Marshal(r)
	var m = make(map[string]interface{})
	_ = json.Unmarshal(j, &m)
	return m
}