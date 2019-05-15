package util

import "errors"

var (
	ErrUnknownAncestor = errors.New("父块异常")

	ErrBlockTime = errors.New("区块时间异常")

	ErrInvalidNumber = errors.New("无效的区块高度")

	ErrUnknowBlock = errors.New("未知的区块")

	ErrPublicKeyRecoverFail = errors.New("公钥回复失败")

	ErrPublicKeyCompareFail = errors.New("公钥比对失败")

	ErrPrivateKeySignFail = errors.New("私钥加签失败")

	ErrSignerNotExist = errors.New("签名者不存在于签名列表")

	ErrSignerNotLocal = errors.New("签名者不是当前节点")
)
