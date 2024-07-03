package config

import (
	"go.mau.fi/whatsmeow/proto/waCompanionReg"
)

var (
	AppVersion             = "v1"
	AppPort                = "3030"
	AppDebug               = false
	AppOs                  = "Midincihuy"
	AppPlatform            = waCompanionReg.DeviceProps_PlatformType(1)
	AppBasicAuthCredential string

	PathQrCode    = "statics/qrcode"
	PathSendItems = "statics/senditems"
	PathMedia     = "statics/media"
	PathStorages  = "storages"

	DBName = "whatsapp.db"

	WhatsappAutoReplyMessage    string
	WhatsappWebhook             string
	WhatsappLogLevel                  = "ERROR"
	WhatsappSettingMaxFileSize  int64 = 50000000  // 50MB
	WhatsappSettingMaxVideoSize int64 = 100000000 // 100MB
	WhatsappTypeUser                  = "@s.whatsapp.net"
	WhatsappTypeGroup                 = "@g.us"
)
