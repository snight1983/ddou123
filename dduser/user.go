package dduser

//InitUser .
func InitUser(userFolder string) error {

	return nil
}

//Login .
func Login(ID string, passWD string, eqID string, Ver string, Addr string, IP string,
	Lon float32, Lan float32, acType AccountType, routeid string) (string, string, bool) {
	return userLogin(ID, passWD, eqID, Ver, Addr, IP, Lon, Lan, acType, routeid)
}

//Regedit .
func Regedit(id string, passwd string, eqid string, ver string, addr string, ip string,
	lon float32, lan float32, tp AccountType, routeID string) (string, bool) {
	res, ok := userRegedit(id, passwd, eqid, ver, addr, ip, lon, lan, tp, routeID)
	if true == ok {
	}
	return res, ok
}

//SaveUserFile 保存用户文件
func SaveUserFile(uuid string) error {
	return nil
}
