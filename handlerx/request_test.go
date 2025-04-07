package handlerx

import (
	"reflect"
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
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				bodyBytes: []byte(`{
  "base": {
    "username": "test10@myibc.net",
    "nickname": "",
    "groupId": "1",
    "phone": "",
    "email": "",
    "password": "042161e8476a02f81d759b8da24aada789ab855a7a83c2060da6a6266c7693b77259136a7426e2b73ff1a776fb9a749a1db3cebcfacc159184eded0039539ec467EFC94222E6F56E720D7F7F4AB99A631BCC5C519795A98D490804156390F3D15D20DE2BF3347A7267B5",
    "enableDeviceCodeAuth": true,
    "accessAlgorithm": "SM9",
    "sm2UsernameSource": "1",
    "sm2Cert": ""
  }
}`),
				req: CreateRequest{
					Base: CreateBase{
						UserBase: UserBase{
							Nickname:             "test",
							GroupId:              1,
							Phone:                "",
							Email:                "",
							Status:               0,
							EnableDeviceCodeAuth: true,
						},
						Username: "test10@myibc.net",
						Password: "042161e8476a02f81d759b8da24aada789ab855a7a83c2060da6a",
					},
				},
			},
			want: []byte(`{"base":{"nickname":"test","groupId":1,"phone":"","email":"","status":0,"enableDeviceCodeAuth":true,"username":"test10@myibc.net","password":"042161e8476a02f81d759b8da24aada789ab855a7a83c2060da6a6266c7693b77259136a7426e2b73ff1a776fb9a749a1db3cebcfacc159184eded0039539ec467EFC94222E6F56E720D7F7F4AB99A631BCC5C519795A98D490804156390F3D15D20DE2BF3347A7267B5","accessAlgorithm":"SM9","sm2UsernameSource":"1","sm2Cert":""}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := weaklyDecodeRequest(tt.args.bodyBytes, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("weaklyDecodeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("weaklyDecodeRequest() got = %s, want %s", got, tt.want)
			}
		})
	}
}
