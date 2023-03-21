package app

func (this_app *PWMan_App) VerifyKey(opt_key ...string) bool {
	if len(opt_key) == 0 {
		return this_app.VerifyKey(this_app.Key)
	}
	if len(opt_key) == 1 {
		for _, ch := range opt_key[0] {
			if (ch < 48 || ch > 57) && ch != 0 && ch != 10 {
				return false
			}
		}
		return true
	}
	return false
}
