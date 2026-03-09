package client

import (
	"fmt"

	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

// PermissionMember represents a permission member
type PermissionMember struct {
	MemberType string `json:"member_type"` // "email", "openid", "userid", "unionid", "openchat", "opendepartmentid", "groupid", "wikispaceid"
	MemberID   string `json:"member_id"`
	Perm       string `json:"perm"` // "view", "edit", "full_access"
}

// AddPermission adds permission to a document
func AddPermission(docToken string, docType string, member PermissionMember, notify bool) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	memberObj := larkdrive.NewBaseMemberBuilder().
		MemberType(member.MemberType).
		MemberId(member.MemberID).
		Perm(member.Perm).
		Build()

	req := larkdrive.NewCreatePermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		NeedNotification(notify).
		BaseMember(memberObj).
		Build()

	resp, err := client.Drive.PermissionMember.Create(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("添加权限失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("添加权限失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// ListPermission lists all permissions for a document
func ListPermission(docToken string, docType string) ([]*larkdrive.Member, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	req := larkdrive.NewListPermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		Build()

	resp, err := client.Drive.PermissionMember.List(Context(), req, tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("获取权限列表失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("获取权限列表失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return resp.Data.Items, nil
}

// DeletePermission removes permission from a document
func DeletePermission(docToken string, docType string, memberType string, memberID string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	req := larkdrive.NewDeletePermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		MemberId(memberID).
		MemberType(memberType).
		Build()

	resp, err := client.Drive.PermissionMember.Delete(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("删除权限失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("删除权限失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// TransferOwnership 转移文档所有权
func TransferOwnership(docToken string, docType string, memberType string, memberID string, notify bool, removeOldOwner bool, stayPut bool, oldOwnerPerm string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	owner := larkdrive.NewOwnerBuilder().
		MemberType(memberType).
		MemberId(memberID).
		Build()

	req := larkdrive.NewTransferOwnerPermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		NeedNotification(notify).
		RemoveOldOwner(removeOldOwner).
		StayPut(stayPut).
		OldOwnerPerm(oldOwnerPerm).
		Owner(owner).
		Build()

	resp, err := client.Drive.PermissionMember.TransferOwner(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("转移所有权失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("转移所有权失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// PublicPermissionUpdate 公共权限更新参数
type PublicPermissionUpdate struct {
	ExternalAccess  *bool   // 外部访问
	SecurityEntity  *string // 谁可以复制内容、创建副本、打印、下载
	CommentEntity   *string // 谁可以评论
	ShareEntity     *string // 谁可以添加和管理协作者
	LinkShareEntity *string // 链接分享范围
	InviteExternal  *bool   // 是否允许邀请外部人
}

// GetPublicPermission 获取文档公共权限设置
func GetPublicPermission(docToken, docType string) (*larkdrive.PermissionPublic, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	req := larkdrive.NewGetPermissionPublicReqBuilder().
		Token(docToken).
		Type(docType).
		Build()

	resp, err := client.Drive.PermissionPublic.Get(Context(), req, tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("获取公共权限设置失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("获取公共权限设置失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("获取公共权限设置返回数据为空")
	}

	return resp.Data.PermissionPublic, nil
}

// UpdatePublicPermissionV2 更新文档公共权限设置（支持所有字段）
func UpdatePublicPermissionV2(docToken, docType string, update PublicPermissionUpdate) (*larkdrive.PermissionPublic, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	builder := larkdrive.NewPermissionPublicRequestBuilder()
	if update.ExternalAccess != nil {
		builder.ExternalAccess(*update.ExternalAccess)
	}
	if update.SecurityEntity != nil {
		builder.SecurityEntity(*update.SecurityEntity)
	}
	if update.CommentEntity != nil {
		builder.CommentEntity(*update.CommentEntity)
	}
	if update.ShareEntity != nil {
		builder.ShareEntity(*update.ShareEntity)
	}
	if update.LinkShareEntity != nil {
		builder.LinkShareEntity(*update.LinkShareEntity)
	}
	if update.InviteExternal != nil {
		builder.InviteExternal(*update.InviteExternal)
	}

	req := larkdrive.NewPatchPermissionPublicReqBuilder().
		Token(docToken).
		Type(docType).
		PermissionPublicRequest(builder.Build()).
		Build()

	resp, err := client.Drive.PermissionPublic.Patch(Context(), req, tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("更新公共权限失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("更新公共权限失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("更新公共权限返回数据为空")
	}

	return resp.Data.PermissionPublic, nil
}

// CreatePublicPassword 创建文档分享密码
func CreatePublicPassword(docToken, docType string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return "", tokenErr
	}

	req := larkdrive.NewCreatePermissionPublicPasswordReqBuilder().
		Token(docToken).
		Type(docType).
		Build()

	resp, err := client.Drive.PermissionPublicPassword.Create(Context(), req, tokenOpt)
	if err != nil {
		return "", fmt.Errorf("创建文档密码失败: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("创建文档密码失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil {
		return "", fmt.Errorf("创建文档密码返回数据为空")
	}

	return StringVal(resp.Data.Password), nil
}

// DeletePublicPassword 删除文档分享密码
func DeletePublicPassword(docToken, docType string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	req := larkdrive.NewDeletePermissionPublicPasswordReqBuilder().
		Token(docToken).
		Type(docType).
		Build()

	resp, err := client.Drive.PermissionPublicPassword.Delete(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("删除文档密码失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("删除文档密码失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// UpdatePublicPassword 刷新文档分享密码
func UpdatePublicPassword(docToken, docType string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return "", tokenErr
	}

	req := larkdrive.NewUpdatePermissionPublicPasswordReqBuilder().
		Token(docToken).
		Type(docType).
		Build()

	resp, err := client.Drive.PermissionPublicPassword.Update(Context(), req, tokenOpt)
	if err != nil {
		return "", fmt.Errorf("刷新文档密码失败: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("刷新文档密码失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil {
		return "", fmt.Errorf("刷新文档密码返回数据为空")
	}

	return StringVal(resp.Data.Password), nil
}

// BatchAddPermission 批量添加协作者权限
func BatchAddPermission(docToken, docType string, members []*PermissionMember, notify bool) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	baseMembers := make([]*larkdrive.BaseMember, 0, len(members))
	for _, m := range members {
		baseMembers = append(baseMembers, larkdrive.NewBaseMemberBuilder().
			MemberType(m.MemberType).
			MemberId(m.MemberID).
			Perm(m.Perm).
			Build())
	}

	req := larkdrive.NewBatchCreatePermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		NeedNotification(notify).
		Body(larkdrive.NewBatchCreatePermissionMemberReqBodyBuilder().
			Members(baseMembers).
			Build()).
		Build()

	resp, err := client.Drive.PermissionMember.BatchCreate(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("批量添加权限失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("批量添加权限失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// AuthPermission 判断当前用户对文档的权限
func AuthPermission(docToken, docType, action string) (*larkdrive.AuthPermissionMemberRespData, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	req := larkdrive.NewAuthPermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		Action(action).
		Build()

	resp, err := client.Drive.PermissionMember.Auth(Context(), req, tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("权限判断失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("权限判断失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("权限判断返回数据为空")
	}

	return resp.Data, nil
}

// UpdatePermission 更新协作者权限
func UpdatePermission(docToken string, docType string, memberID string, memberType string, perm string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	req := larkdrive.NewUpdatePermissionMemberReqBuilder().
		Token(docToken).
		Type(docType).
		MemberId(memberID).
		BaseMember(larkdrive.NewBaseMemberBuilder().
			MemberType(memberType).
			MemberId(memberID).
			Perm(perm).
			Build()).
		Build()

	resp, err := client.Drive.PermissionMember.Update(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("更新权限失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("更新权限失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}
