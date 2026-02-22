package cache

import "time"

const (
	KeyUserLoginToken             = "user:login:token:%s"
	KeySolanaSignatureProcess     = "solana:signature:process:%s"
	KeyUsernameModificationLimit  = "user:username:modification:%d"
	KeyDisconnection              = "solana:subevent:disconnect:%s"
	KeyCreateKline                = "kline:create:%s"
	KeyUpdateHeld                 = "token:held:%d:%s"
	KeySolanaSignature            = "solana:signature:%s"
	KeyUserRegisteredNotification = "notification:user:registered:%d"
	KeySearchTrend                = "token:search:trends"
	KeyUserSearch                 = "token:user:search:%d"
	KeyChatMsg                    = "chat:msg"
	KeyListOnRaydium              = "raydium:create:%s"
)

const (
	ExpireUserRegisteredNotification = 1 * time.Hour
)
