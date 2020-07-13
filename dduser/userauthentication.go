package dduser

import (
	"crypto/sha1"
	"ddou123/ddlib"
	"fmt"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	guuid "github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

func userLogin(uid string, passwd string, equipmentID string, version string,
	addr string, ip string, lon float32, lan float32,
	acType AccountType, routeID string) (string, string, bool) {

	if len(passwd) > 0 {
		h := sha1.New()
		if _, err := h.Write([]byte(passwd)); nil != err {
			return "", "", false
		}
		passwd = fmt.Sprintf("%x", h.Sum(nil))
	}

	var ok bool = false
	var uuid interface{} = nil
	var userInfo interface{} = nil

	if uuid, ok = gUserIDMap.Load(uid); false == ok {

		ddlib.GLoger.Error("login_unregedit",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID))
		return "登录失败.1", "", false
	}

	if userInfo, ok = gUserInfoMap.Load(uuid); false == ok {
		ddlib.GLoger.Error("login_uninit",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID))
		return "登录失败.2", "", false
	}

	if acType == ACEmail || acType == ACPhone {
		if userInfo.(*UserInfo).UBase.PASSWD != passwd {
			ddlib.GLoger.Info("login_passwd",
				zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
				zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
				zap.Float32("lan", lan), zap.String("routeid", routeID))
			return "登录失败.3", "", false
		}
	}

	if ok = ValidateToken(userInfo.(*UserInfo).Token); true == ok {

		userInfo.(*UserInfo).LoginAddr = addr
		userInfo.(*UserInfo).LoginIP = ip
		userInfo.(*UserInfo).EquipmentID = equipmentID
		userInfo.(*UserInfo).Version = version
		userInfo.(*UserInfo).Longitude = lon
		userInfo.(*UserInfo).Lalatitude = lan
		userInfo.(*UserInfo).RouteID = routeID

		ddlib.GLoger.Info("login_validateToken",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID))

		return userInfo.(*UserInfo).UBase.UUID, userInfo.(UserInfo).Token, true
	}

	if signedToken, ok := CreateToken(userInfo.(*UserInfo).UBase.UUID, equipmentID, version); true == ok {

		userInfo.(*UserInfo).LoginAddr = addr
		userInfo.(*UserInfo).LoginIP = ip
		userInfo.(*UserInfo).EquipmentID = equipmentID
		userInfo.(*UserInfo).Version = version
		userInfo.(*UserInfo).Longitude = lon
		userInfo.(*UserInfo).Lalatitude = lan
		userInfo.(*UserInfo).RouteID = routeID
		userInfo.(*UserInfo).Token = signedToken
		userInfo.(*UserInfo).LoginTime = time.Now().Unix()
		userInfo.(*UserInfo).UserFileDBPath = filepath.Join(gUserBaseDBPath, "userfile")
		userInfo.(*UserInfo).UserRelationShipBPath = filepath.Join(gUserBaseDBPath, "userrelationship")

		ddlib.GLoger.Info("login_createToken",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID), zap.String("token", signedToken),
			zap.Int64("time", userInfo.(*UserInfo).LoginTime))

		go userInfo.(*UserInfo).loadUserRes()

		return userInfo.(*UserInfo).UBase.UUID, userInfo.(UserInfo).Token, true
	}

	//if nil == userInfo.(*UserInfo) {
	//}
	return "", "", false
}

func userRegedit(uid string, passwd string, equipmentID string, version string,
	addr string, ip string, lon float32, lan float32,
	acType AccountType, routeID string) (string, bool) {

	if len(passwd) > 0 {
		h := sha1.New()
		if _, err := h.Write([]byte(passwd)); nil != err {
			return "", false
		}
		passwd = fmt.Sprintf("%x", h.Sum(nil))
	}
	var err error = nil
	//var ok bool = false
	var uuid interface{} = nil
	var userIDMapDB *leveldb.DB = nil
	var userBaseDB *leveldb.DB = nil
	//if uuid, ok = gUserIDMap.Load(uid); true == ok {
	//	return "账户已注册,请登录", false
	//}
	uuid = guuid.New()
	userBase := getUserBase()
	userBase.UUID = uuid.(guuid.UUID).String()
	userBase.PASSWD = passwd
	userBase.Actype = acType

	switch acType {
	case ACPhone:
		userBase.Phone = uid
	case ACEmail:
		userBase.Email = uid
	case ACQQ:
		userBase.QQ = uid
	case ACWX:
		userBase.WX = uid
	default:
		{
			ddlib.GLoger.Error("reg_unknowactype",
				zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
				zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
				zap.Float32("lan", lan), zap.String("routeid", routeID))
			return "注册失败,请重试.1", false
		}

	}

	userBase.PASSWD = passwd

	if userIDMapDB, err = leveldb.OpenFile(gUserIDMapDBPath, nil); nil != err {
		ddlib.GLoger.Error("reg_idmapdb",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID), zap.String("path", gUserIDMapDBPath))
		return "注册失败,请重试.2", false
	}
	defer userIDMapDB.Close()

	if userBaseDB, err = leveldb.OpenFile(gUserBaseDBPath, nil); nil != err {
		ddlib.GLoger.Error("reg_ubasedb",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID),
			zap.String("path", gUserBaseDBPath), zap.String("error", fmt.Sprintf("%s", err)))
		return "注册失败,请重试.3", false
	}
	defer userBaseDB.Close()

	if err = updateUserBase(userBaseDB, userBase); err != nil {
		ddlib.GLoger.Error("reg_ubaseupdate",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID),
			zap.String("path", gUserBaseDBPath), zap.String("error", fmt.Sprintf("%s", err)))
		return "注册失败,请重试.4", false
	}

	if nil != userIDMapDB.Put([]byte(uid), []byte(uuid.(guuid.UUID).String()), nil) {
		ddlib.GLoger.Error("reg_idmapupdate",
			zap.String("uid", uid), zap.String("addr", addr), zap.String("ip", ip),
			zap.String("eq", equipmentID), zap.String("ver", version), zap.Float32("lon", lon),
			zap.Float32("lan", lan), zap.String("routeid", routeID),
			zap.String("path", gUserBaseDBPath), zap.String("error", fmt.Sprintf("%s", err)))
		return "注册失败,请重试.5", false
	}

	if nil != createRoot(uuid.(guuid.UUID).String()) {
		return "注册失败,请重试.6", false
	}

	return "", true
}

type jwtCustomClaims struct {
	UUID        string `json:"uid"`
	EquipmentID string `json:"equipment"`
	Version     string `json:"version"`
	jwt.StandardClaims
}

// CreateToken .
func CreateToken(uuid string, equipmentID string, version string) (signedToken string, success bool) {
	claims := &jwtCustomClaims{
		uuid,
		equipmentID,
		version,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(ddlib.SIGNINGKEY))
	if nil != err {
		return
	}
	success = true
	return
}

// ValidateToken .
func ValidateToken(signedToken string) (success bool) {
	token, err := jwt.ParseWithClaims(signedToken, &jwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected login method %v", token.Header["alg"])
			}
			return []byte(ddlib.SIGNINGKEY), nil
		})

	if err != nil {
		return
	}

	_, ok := token.Claims.(*jwtCustomClaims)
	if ok && token.Valid {
		success = true
		return
	}
	return
}
