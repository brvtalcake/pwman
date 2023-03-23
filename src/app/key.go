package app

func (pwman PWMan_App) ConvertKey() []byte {
	return []byte(pwman.Key)
}

func (this_app *PWMan_App) VerifyKey(opt_key ...string) bool {
	if len(opt_key) == 0 {
		return this_app.VerifyKey(this_app.Key)
	}
	if len(opt_key) >= 1 {
		this_app.Byte_key = this_app.ConvertKey()
		if len(this_app.Byte_key) <= 1 || this_app.Byte_key == nil {
			goto false_ret
		} else if len(this_app.Byte_key) > 32 {
			this_app.Byte_key = this_app.Byte_key[:32]
		} else if len(this_app.Byte_key) < 32 {
			byte_key_cpy := make([]byte, len(this_app.Byte_key))
			copy(byte_key_cpy, this_app.Byte_key)
			for len(this_app.Byte_key) < 32 {
				this_app.Byte_key = append(this_app.Byte_key, byte_key_cpy...)
				if len(this_app.Byte_key) > 32 {
					this_app.Byte_key = this_app.Byte_key[:32]
				}
			}
		}
		return len(opt_key[0]) > 0 && len(this_app.Byte_key) == 32
	} else {
		return false
	}
false_ret:
	return false
}
