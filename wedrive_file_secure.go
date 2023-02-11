package workwx

type WedriveFileSecureAddRequest struct {
	FileID   string                `json:"fileid"`
	AuthInfo []WedriveAuthInfoItem `json:"auth_info"`
}

type WedriveFileSecureDelRequest struct {
	FileID   string                `json:"fileid"`
	AuthInfo []WedriveUserInfoItem `json:"auth_info"`
}

type WedriveFileSecureSetRequest struct {
	FileID    string  `json:"fileid"`
	AuthScope uint32  `json:"auth_scope"`     //	权限范围：1:指定人 2:企业内 3:企业外 4: 企业内需管理员审批（仅有管理员时可设置） 5: 企业外需管理员审批（仅有管理员时可设置）
	Auth      *uint32 `json:"auth,omitempty"` //	权限信息 普通文档： 1:仅浏览（可下载) 4:仅预览（仅专业版企业可设置）；如果不填充此字段为保持原有状态 微文档： 1:仅浏览（可下载）；如果不填充此字段为保持原有状态
}

type WedriveFileSecureShareRequest struct {
	FileID string `json:"fileid"`
}

type WedriveFileSecureShareResponse struct {
	Error
	Url string `json:"share_url"` // 邀请链接
}

type WedriveFileSecurePermissionRequest WedriveFileSecureShareRequest

type WedriveFileSecurePermissionResponse struct {
	Error
	ShareRange        WedriveFileSecureShareRange       `json:"share_range"`
	SecureSetting     WedriveFileSecureSetting          `json:"secure_setting"`
	InheritFatherAuth WedriveFileSecureIneritFatherAuth `json:"inherit_father_auth"`
	FileMemberList    []WedriveAuthInfoItem             `json:"member_list,omitempty"`
	Watermark         WedriveFileSecureWatermark        `json:"watermark"`
}

// 文件分享设置
type WedriveFileSecureShareRange struct {
	EnableCorpInternal bool    `json:"enable_corp_internal"`         // 是否为企业内可访问
	CorpInternalAuth   *uint32 `json:"corp_internal_auth,omitempty"` //	企业内权限信息 普通文档： 1:仅浏览（可下载) 4:仅预览（仅专业版企业可设置）；如果不填充此字段为保持原有状态 微文档： 1:仅浏览（可下载）；如果不填充此字段为保持原有状态
	EnableCorpExternal bool    `json:"enable_corp_external"`         // 是否为企业外可访问
	CorpExternalAuth   *uint32 `json:"corp_external_auth,omitempty"` //	企业外权限信息 普通文档： 1:仅浏览（可下载) 4:仅预览（仅专业版企业可设置）；如果不填充此字段为保持原有状态 微文档： 1:仅浏览（可下载）；如果不填充此字段为保持原有状态
}

type WedriveFileSecureSetting struct {
	EnableReadonly        bool `json:"enable_readonly_copy"`
	ModifyOnlyByAdmin     bool `json:"modify_only_by_admin"`
	EnableReadonlyComment bool `json:"enable_readonly_comment"`
	BanShareExternal      bool `json:"ban_share_external"`
}

type WedriveFileSecureIneritFatherAuth struct {
	AuthList []WedriveAuthInfoItem `json:"auth_list"`
	Inherit  bool                  `json:"inherit"`
}

type WedriveFileSecureWatermark struct {
	WedriveFileSecureWatermarkSetting
	ForceByAdmin      *bool `json:"force_by_admin,omitempty"`       // 管理员是否强制要求使用水印，此字段不填则保持原样
	ForceBySpaceAdmin *bool `json:"force_by_space_admin,omitempty"` // 空间管理员是否强制要求使用水印，此字段不填则保持原样
}

type WedriveFileSecureWatermarkSetting struct {
	Text            *string `json:"text,omitempty"`              // 水印文字，此字段不填则保持原样
	MarginType      *uint32 `json:"margin_type,omitempty"`       // 水印类型。1：低密度水印， 2： 高密度水印，此字段不填则保持原样
	ShowVisitorName *bool   `json:"show_visitor_name,omitempty"` //	是否显示访问人名称，此字段不填则保持原样（仅专业版支持）
	ShowText        *bool   `json:"show_text,omitempty"`         // 是否展示水印文本，此字段不填则保持原样
}

type WedriveFileSecureSetPermissionRequest struct {
	FileID    string                            `json:"fileid"`
	Watermark WedriveFileSecureWatermarkSetting `json:"watermark"`
}

func (p *WedriveFile) SecureAdd(param *WedriveFileSecureAddRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecureAdd, param, &ret)

	return
}

func (p *WedriveFile) SecureDel(param *WedriveFileSecureDelRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecureDel, param, &ret)

	return
}

func (p *WedriveFile) SecureSet(param *WedriveFileSecureSetRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecureSet, param, &ret)

	return
}

func (p *WedriveFile) SecureShare(param *WedriveFileSecureShareRequest) (ret WedriveFileSecureShareResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecureShare, param, &ret)

	return
}

func (p *WedriveFile) SecurePermission(param *WedriveFileSecurePermissionRequest) (ret WedriveFileSecurePermissionResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecurePermission, param, &ret)

	return
}

func (p *WedriveFile) SecureSetPermission(param *WedriveFileSecureSetPermissionRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileSecureSetPermission, param, &ret)

	return
}

const (
	wedriveApiFileSecureAdd           = "file_acl_add"
	wedriveApiFileSecureDel           = "file_acl_del"
	wedriveApiFileSecureSet           = "file_setting"
	wedriveApiFileSecureShare         = "file_share"
	wedriveApiFileSecurePermission    = "get_file_permission"
	wedriveApiFileSecureSetPermission = "file_secure_setting"
)
