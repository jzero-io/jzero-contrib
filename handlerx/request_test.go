package handlerx

import (
	"testing"
)

type UserBase struct {
	Nickname             string `json:"nickname" olabel:"用户姓名"`                                            // 用户名称
	GroupId              int    `json:"groupId" olabel:"所属组织id"`                                           // 组织 id
	Phone                string `json:"phone" olabel:"手机号"`                                                // 手机号
	Email                string `json:"email" olabel:"邮箱"`                                                 // 邮箱
	Status               int8   `json:"status,optional" olabel:"状态" ovalue:"0=正常,1=已锁定,2=已失效,3=已注销,4=已禁用"` // 用户状态 0-正常 1-已锁定 2-已失效 3-已注销 4-已禁用
	EnableDeviceCodeAuth bool   `json:"enableDeviceCodeAuth" olabel:"开启设备特征码认证"`
}

type CreateBase struct {
	UserBase
	Username          string `json:"username"`                    // 用户账号
	Password          string `json:"password"`                    // 登陆密码
	AccessAlgorithm   string `json:"accessAlgorithm,default=SM9"` // 接入算法 SM2/SM9
	Sm2UsernameSource string `json:"sm2UsernameSource,optional"`  // sm2 账号来源 1-手动输入 2-UKEY 获取 3-上传证书
	Sm2Cert           string `json:"sm2Cert,optional"`            // 当账号来源为上传证书则为文件路径, 账号来源为 UKEY 获取则为证书内容
}

type CreateRequest struct {
	Base CreateBase `json:"base"` // 用户基本信息
}

func Test_weaklyDecodeRequest(t *testing.T) {
	type args struct {
		bodyBytes []byte
		req       any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				bodyBytes: []byte(`{"base":{"username":"test10@myibc.net","nickname":"","groupId":"1","phone":"","email":"","password":"043064917508e684c87033314275c923b1df07341cad8bec49a4464389954d893cf778252bf7d740047369a7d0f544f193e0e97db69c83755c71d483978b62df966ACBC745216B44585622F2372D38CEF785786C4295772F08E71C18116D3E65B913C04C2815CA275725","enableDeviceCodeAuth":true,"accessAlgorithm":"SM9","sm2UsernameSource":"1","sm2Cert":""}}`),
				req:       &CreateRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := weaklyDecodeRequest(tt.args.bodyBytes, tt.args.req)
			if err != nil {
				t.Errorf("weaklyDecodeRequest() error = %v", err)
				return
			}
			t.Log(string(got))
		})
	}
}
