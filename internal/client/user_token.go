package client

import (
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/riba2534/feishu-cli/internal/auth"
	"github.com/riba2534/feishu-cli/internal/config"
)

// noopRequestOption 是一个空的请求选项，用于 tenant 模式
// 必须返回非 nil 的函数，因为 Go 中将 nil 传给可变参数会导致切片包含 nil 元素
func noopRequestOption(*larkcore.RequestOption) {}

// GetUserTokenOption 根据当前配置的 token-mode 返回请求选项
// 如果应该使用 user_access_token，返回 WithUserAccessToken 选项
// 如果应该使用 tenant_access_token（默认），返回 no-op 选项（不能返回 nil）
//
// 优先级：
//   1. 如果配置文件/环境变量中设置了静态 user_access_token，且 token-mode 不是 tenant，直接使用
//   2. 如果通过 OAuth 登录过，根据 token-mode 决定是否使用
//   3. 否则使用默认的 tenant_access_token
func GetUserTokenOption() (larkcore.RequestOptionFunc, error) {
	cfg := config.Get()

	mode, err := auth.ParseTokenMode(cfg.TokenMode)
	if err != nil {
		return nil, err
	}

	// 如果强制使用 tenant 模式，返回 no-op 选项（不能返回 nil）
	if mode == auth.TokenModeTenant {
		return noopRequestOption, nil
	}

	// 优先使用配置中的静态 user_access_token（向后兼容）
	if cfg.UserAccessToken != "" {
		return larkcore.WithUserAccessToken(cfg.UserAccessToken), nil
	}

	// 使用 OAuth 管理的 user_access_token
	opt, err := auth.UserTokenRequestOption(mode)
	if err != nil {
		return nil, err
	}
	// 如果 OAuth 返回 nil（auto 模式下未登录），使用 no-op 选项
	if opt == nil {
		return noopRequestOption, nil
	}
	return opt, nil
}
