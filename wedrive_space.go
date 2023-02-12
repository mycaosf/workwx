package workwx

type WedriveSpace struct {
	Token
}

type WedriveAuthInfoItem struct {
	Type         uint32  `json:"type"` // 1: person, 2: department
	UserID       *string `json:"userid,omitempty"`
	DepartmentID *uint32 `json:"departmentid,omitempty"`
	Auth         int     `json:"auth"` // 1: only download, 4: may preview, 7: administrator
}

type WedriveAutoList struct {
	AuthInfo   []WedriveAuthInfoItem `json:"auth_info"`
	QuitUserID []string              `json:"quit_userid"`
}

type WedriveSpaceInfo struct {
	SpaceID      string          `json:"spaceid"`
	SpaceName    string          `json:"space_name"`
	AuthList     WedriveAutoList `json:"auth_list"`
	SpaceSubType uint32          `json:"space_sub_type"`
}

type WedriveSpaceCreateRequest struct {
	Name         string                `json:"space_name"`
	AuthInfo     []WedriveAuthInfoItem `json:"auth_info"`
	SpaceSubType uint32                `json:"space_sub_type"`
}

type WedriveSpaceID struct {
	SpaceID string `json:"spaceid"`
}

type WedriveSpaceCreateResponse struct {
	Error
	WedriveSpaceID
}

type WedriveSpaceDeleteRequest WedriveSpaceID

type WedriveSpaceRenameRequest struct {
	SpaceID string `json:"spaceid"`
	Name    string `json:"space_name"`
}

type WedriveSpaceListRequest WedriveSpaceID

type WedriveSpaceListResponse struct {
	Error
	SpaceInfo WedriveSpaceInfo `json:"space_info"`
}

func (p *WedriveSpace) Create(param *WedriveSpaceCreateRequest) (ret WedriveSpaceCreateResponse) {
	wedrivePost(&p.Token, wedriveApiSpaceCreate, param, &ret)

	return
}

func (p *WedriveSpace) Delete(param *WedriveSpaceDeleteRequest) (ret Error) {
	wedrivePost(&p.Token, wedriveApiSpaceDelete, param, &ret)

	return
}

func (p *WedriveSpace) Rename(param *WedriveSpaceRenameRequest) (ret Error) {
	wedrivePost(&p.Token, wedriveApiSpaceRename, param, &ret)

	return
}

func (p *WedriveSpace) List(param *WedriveSpaceListRequest) (ret WedriveSpaceListResponse) {
	wedrivePost(&p.Token, wedriveApiSpaceInfo, param, &ret)

	return
}

const (
	wedriveApiSpaceCreate = "space_create"
	wedriveApiSpaceDelete = "space_dismiss"
	wedriveApiSpaceRename = "space_rename"
	wedriveApiSpaceInfo   = "space_info"
)
