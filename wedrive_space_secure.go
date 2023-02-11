package workwx

type WedriveSpaceSecureAddRequest struct {
	SpaceID  string                `json:"spaceid"`
	AuthInfo []WedriveAuthInfoItem `json:"auth_info"`
}

type WedriveSpaceSecureDelRequest struct {
	SpaceID  string                `json:"spaceid"`
	AuthInfo []WedriveUserInfoItem `json:"auth_info"`
}

type WedriveUserInfoItem struct {
	Type         uint32  `json:"type"` // 1: person, 2: department
	UserID       *string `json:"userid,omitempty"`
	DepartmentID *uint32 `json:"departmentid,omitempty"`
}

type WedriveSpaceSecureSetRequest struct {
	SpaceID                string  `json:"spaceid"`
	EnableWatermark        *bool   `json:"enable_watermark,omitempty"`         //（本字段仅专业版企业可设置）启用水印。false:关 true:开 ;如果不填充此字段为保持原有状态
	EnableConfidentialMode *bool   `json:"enable_confidential_mode,omitempty"` // 是否开启保密模式。false:关 true:开 如果不填充此字段为保持原有状态
	DefaultFileScope       *uint32 `json:"default_file_scope,omitempty"`       // 文件默认可查看范围。1:仅成员；2:企业内。如果不填充此字段为保持原有状态
	BanShareExternal       *bool   `json:"ban_share_external,omitempty"`       //	是否禁止文件分享到企业外｜false:关 true:开 如果不填充此字段为保持原有状态
}

type WedriveSpaceSecureShareRequest struct {
	SpaceID string `json:"spaceid"`
}

type WedriveSpaceSecureShareResponse struct {
	Error
	Url string `json:"space_share_url"` // 邀请链接
}

type WedriveSpaceSecureInfoRequest WedriveSpaceSecureShareRequest

type WedriveSpaceSecureInfoResponse struct {
	Error
	SpaceInfo WedriveSpaceSecureInfo `json:"space_info"`
}

type WedriveSpaceSecureInfo struct {
	SpaceID       string                    `json:"spaceid"`
	SpaceName     string                    `json:"space_name"`
	AuthList      WedriveAutoList           `json:"auth_list"`
	SpaceSubType  uint32                    `json:"space_sub_type"`
	SecureSetting WedriveSpaceSecureSetting `json:"secure_setting"`
}

type WedriveSpaceSecureSetting struct {
	EnableWatermark                   bool   `json:"enable_watermark"`
	AddMemberOnlyAdmin                bool   `json:"add_member_only_admin"`
	EnableShareUrl                    bool   `json:"enable_share_url"`
	ShareUrlNoApprove                 bool   `json:"share_url_no_approve"`
	ShareUrlNoApproveDefaultAuth      int    `json:"share_url_no_approve_default_auth"`
	EnableShareExternal               bool   `json:"enable_share_external"`
	EnableShareExternalAdmin          bool   `json:"enable_share_external_admin"`
	EnableSpaceAddExternalMember      bool   `json:"enable_space_add_external_member"`
	EnableSpaceAddExternalMemberAdmin bool   `json:"enable_space_add_external_member_admin"`
	EnableConfidentialMode            bool   `json:"enable_confidential_mode"`
	DefaultFileScope                  uint32 `json:"default_file_scope"`
	CreateFileOnlyAdmin               bool   `json:"create_file_only_admin"`
}

func (p *WedriveSpace) SecureAdd(param *WedriveSpaceSecureAddRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiSpaceSecureAdd, param, &ret)

	return
}

func (p *WedriveSpace) SecureDel(param *WedriveSpaceSecureDelRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiSpaceSecureDel, param, &ret)

	return
}

func (p *WedriveSpace) SecureSet(param *WedriveSpaceSecureSetRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiSpaceSecureSet, param, &ret)

	return
}

func (p *WedriveSpace) SecureShare(param *WedriveSpaceSecureShareRequest) (ret WedriveSpaceSecureShareResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiSpaceSecureShare, param, &ret)

	return
}

func (p *WedriveSpace) SecureInfo(param *WedriveSpaceSecureInfoRequest) (ret WedriveSpaceSecureInfoResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiSpaceSecureInfo, param, &ret)

	return
}

const (
	wedriveApiSpaceSecureAdd   = "space_acl_add"
	wedriveApiSpaceSecureDel   = "space_acl_del"
	wedriveApiSpaceSecureSet   = "space_setting"
	wedriveApiSpaceSecureShare = "space_share"
	wedriveApiSpaceSecureInfo  = "new_space_info"
)
