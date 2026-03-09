package client

import (
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// CreateChat 创建群聊
func CreateChat(name, description, ownerID string, userIDs []string, chatType string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return "", tokenErr
	}

	bodyBuilder := larkim.NewCreateChatReqBodyBuilder()
	if name != "" {
		bodyBuilder.Name(name)
	}
	if description != "" {
		bodyBuilder.Description(description)
	}
	if ownerID != "" {
		bodyBuilder.OwnerId(ownerID)
	}
	if len(userIDs) > 0 {
		bodyBuilder.UserIdList(userIDs)
	}
	if chatType != "" {
		bodyBuilder.ChatType(chatType)
	}

	req := larkim.NewCreateChatReqBuilder().
		Body(bodyBuilder.Build()).
		Build()

	resp, err := client.Im.Chat.Create(Context(), req, tokenOpt)
	if err != nil {
		return "", fmt.Errorf("创建群聊失败: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("创建群聊失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data.ChatId == nil {
		return "", fmt.Errorf("群聊已创建但未返回群 ID")
	}

	return *resp.Data.ChatId, nil
}

// GetChat 获取群聊信息
func GetChat(chatID string) (*larkim.GetChatRespData, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	req := larkim.NewGetChatReqBuilder().
		ChatId(chatID).
		Build()

	resp, err := client.Im.Chat.Get(Context(), req, tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("获取群聊信息失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("获取群聊信息失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return resp.Data, nil
}

// UpdateChat 更新群聊信息
func UpdateChat(chatID, name, description, ownerID string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	bodyBuilder := larkim.NewUpdateChatReqBodyBuilder()
	if name != "" {
		bodyBuilder.Name(name)
	}
	if description != "" {
		bodyBuilder.Description(description)
	}
	if ownerID != "" {
		bodyBuilder.OwnerId(ownerID)
	}

	req := larkim.NewUpdateChatReqBuilder().
		ChatId(chatID).
		Body(bodyBuilder.Build()).
		Build()

	resp, err := client.Im.Chat.Update(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("更新群聊信息失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("更新群聊信息失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// DeleteChat 解散群聊
func DeleteChat(chatID string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	req := larkim.NewDeleteChatReqBuilder().
		ChatId(chatID).
		Build()

	resp, err := client.Im.Chat.Delete(Context(), req, tokenOpt)
	if err != nil {
		return fmt.Errorf("解散群聊失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("解散群聊失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// GetChatLink 获取群分享链接
func GetChatLink(chatID string, validityPeriod string) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return "", tokenErr
	}

	bodyBuilder := larkim.NewLinkChatReqBodyBuilder()
	if validityPeriod != "" {
		bodyBuilder.ValidityPeriod(validityPeriod)
	}

	req := larkim.NewLinkChatReqBuilder().
		ChatId(chatID).
		Body(bodyBuilder.Build()).
		Build()

	resp, err := client.Im.Chat.Link(Context(), req, tokenOpt)
	if err != nil {
		return "", fmt.Errorf("获取群分享链接失败: %w", err)
	}

	if !resp.Success() {
		return "", fmt.Errorf("获取群分享链接失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	if resp.Data.ShareLink == nil {
		return "", fmt.Errorf("获取群分享链接成功但未返回链接")
	}

	return *resp.Data.ShareLink, nil
}

// ChatMemberInfo 群成员信息
type ChatMemberInfo struct {
	MemberIDType string `json:"member_id_type"`
	MemberID     string `json:"member_id"`
	Name         string `json:"name"`
	TenantKey    string `json:"tenant_key"`
}

// ListChatMembersResult 群成员列表结果
type ListChatMembersResult struct {
	Items     []*ChatMemberInfo `json:"items"`
	PageToken string            `json:"page_token,omitempty"`
	HasMore   bool              `json:"has_more"`
}

// ListChatMembers 获取群成员列表
func ListChatMembers(chatID, memberIDType string, pageSize int, pageToken string) (*ListChatMembersResult, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return nil, tokenErr
	}

	reqBuilder := larkim.NewGetChatMembersReqBuilder().
		ChatId(chatID)

	if memberIDType != "" {
		reqBuilder.MemberIdType(memberIDType)
	}
	if pageSize > 0 {
		reqBuilder.PageSize(pageSize)
	}
	if pageToken != "" {
		reqBuilder.PageToken(pageToken)
	}

	resp, err := client.Im.ChatMembers.Get(Context(), reqBuilder.Build(), tokenOpt)
	if err != nil {
		return nil, fmt.Errorf("获取群成员列表失败: %w", err)
	}

	if !resp.Success() {
		return nil, fmt.Errorf("获取群成员列表失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	result := &ListChatMembersResult{
		PageToken: StringVal(resp.Data.PageToken),
		HasMore:   BoolVal(resp.Data.HasMore),
	}
	for _, item := range resp.Data.Items {
		result.Items = append(result.Items, &ChatMemberInfo{
			MemberIDType: StringVal(item.MemberIdType),
			MemberID:     StringVal(item.MemberId),
			Name:         StringVal(item.Name),
			TenantKey:    StringVal(item.TenantKey),
		})
	}

	return result, nil
}

// AddChatMembers 添加群成员
func AddChatMembers(chatID, memberIDType string, idList []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	reqBuilder := larkim.NewCreateChatMembersReqBuilder().
		ChatId(chatID).
		Body(larkim.NewCreateChatMembersReqBodyBuilder().
			IdList(idList).
			Build())

	if memberIDType != "" {
		reqBuilder.MemberIdType(memberIDType)
	}

	resp, err := client.Im.ChatMembers.Create(Context(), reqBuilder.Build(), tokenOpt)
	if err != nil {
		return fmt.Errorf("添加群成员失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("添加群成员失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

// RemoveChatMembers 移除群成员
func RemoveChatMembers(chatID, memberIDType string, idList []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tokenOpt, tokenErr := GetUserTokenOption()
	if tokenErr != nil {
		return tokenErr
	}

	reqBuilder := larkim.NewDeleteChatMembersReqBuilder().
		ChatId(chatID).
		Body(larkim.NewDeleteChatMembersReqBodyBuilder().
			IdList(idList).
			Build())

	if memberIDType != "" {
		reqBuilder.MemberIdType(memberIDType)
	}

	resp, err := client.Im.ChatMembers.Delete(Context(), reqBuilder.Build(), tokenOpt)
	if err != nil {
		return fmt.Errorf("移除群成员失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("移除群成员失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	return nil
}
