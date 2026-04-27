package app

import (
	"autossh/src/utils"
	"io"
	"reflect"
	"strconv"
)

// 编辑
func (server *Server) Edit() error {
	keys := []string{"Name", "Ip", "Port", "User", "Password", "Method", "Key", "Alias"}
	for _, key := range keys {
		if err := server.scanVal(key); err != nil {
			return err
		}
	}

	return nil
}

func deftVal(val string) string {
	if val != "" {
		return "(default=" + val + ")"
	} else {
		return ""
	}
}

func (server *Server) scanVal(fieldName string) (err error) {
	elem := reflect.ValueOf(server).Elem()
	field := elem.FieldByName(fieldName)
	switch field.Type().String() {
	case "int":
		prompt := utils.Colored(fieldName+deftVal(strconv.FormatInt(field.Int(), 10))+":", utils.ColorGreen)
		iptStr, err := utils.ReadLine(prompt)
		if err != nil {
			if err == io.EOF || err.Error() == "Interrupt" {
				return io.EOF
			}
			return nil
		}
		if iptStr != "" {
			if ipt, err := strconv.Atoi(iptStr); err == nil {
				field.SetInt(int64(ipt))
			}
		}
	case "string":
		// 根据用户要求：Password 不为空则 Method 自动为 password，否则为 key
		if fieldName == "Method" {
			if server.Password != "" {
				server.Method = "example-password"
			} else {
				server.Method = "key"
			}
			return nil
		}

		// 如果 Password 不为空，Key 强制为空且不参与输入
		if fieldName == "Key" && server.Password != "" {
			server.Key = ""
			return nil
		}

		// 如果是 Key 字段且 Method 为 key，尝试自动关联私钥
		if fieldName == "Key" && server.Method == "key" {
			selectedKey, _ := pickSSHKey()
			if selectedKey != "" {
				field.SetString(selectedKey)
				utils.Infoln("Key: " + selectedKey)
				return nil
			}
		}

		prompt := utils.Colored(fieldName+deftVal(field.String())+":", utils.ColorGreen)
		var iptStr string
		if fieldName == "Password" {
			iptStr, err = utils.ReadPassword(prompt)
		} else {
			iptStr, err = utils.ReadLine(prompt)
		}

		if err != nil {
			if err == io.EOF || err.Error() == "Interrupt" {
				return io.EOF
			}
			return nil
		}

		if iptStr != "" {
			field.SetString(iptStr)
		}
	}

	return nil
}
