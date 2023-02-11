package workwx

import (
	"crypto/sha1"
	"encoding"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
)

type WedriveFileListRequest struct {
	SpaceID  string `json:"spaceid"`
	FatherID string `json:"fatherid"`  // 当前目录的fileid,根目录时为空间spaceid
	SortType uint32 `json:"sort_type"` // 1:名字升序；2:名字降序；3:大小升序；4:大小降序；5:修改时间升序；6:修改时间降序
	Start    uint32 `json:"start"`     // 首次填0, 后续填上一次请求返回的next_start
	Limit    uint32 `json:"limit"`     // 分批拉取最大文件数, 不超过1000
}

type WedriveFileListResponse struct {
	Error
	More      bool            `json:"has_more"`
	NextStart uint32          `json:"next_start"`
	FileList  WedriveFileList `json:"file_list"`
}

type WedriveFileList struct {
	Item []WedriveFileItem `json:"item"`
}

type WedriveFileItem struct {
	FileID     string `json:"fileid"`
	FileName   string `json:"file_name"`
	SpaceID    string `json:"spaceid"`
	FatherID   string `json:"fatherid"` // 当前目录的fileid,根目录时为空间fileid
	FileSize   uint64 `json:"file_size"`
	CTime      uint64 `json:"ctime"`
	MTime      uint64 `json:"mtime"`
	FileType   uint32 `json:"file_type"`   // 1:文件夹 2:文件 3:微文档(文档) 4:微文档(表格) 5:微文档(收集表)
	FileStatus uint32 `json:"file_status"` // 文件状态, 1:正常 2:删除
	Sha        string `json:"sha"`
	Md5        string `json:"md5"`
	Url        string `json:"url"`
}

type WedriveFileInfoRequest struct {
	FileID string `json:"fileid"`
}

type WedriveFileInfoResponse struct {
	Error
	FileInfo WedriveFileItem `json:"file_info"`
}

type WedriveFileCreateRequest struct {
	SpaceID  string `json:"spaceid"`
	FatherID string `json:"fatherid"`  // 当前目录的fileid,根目录时为空间spaceid
	FileType uint32 `json:"file_type"` // 1:文件夹 3:微文档(文档) 4:微文档(表格)
	FileName string `json:"file_name"`
}

type WedriveFileCreateResponse struct {
	Error
	FileID string `json:"fileid"`
	Url    string `json:"url"`
}

type WedriveFileDeleteRequest struct {
	FileID []string `json:"fileid"`
}

type WedriveFileRenameRequest struct {
	FileID  string `json:"fileid"`
	NewName string `json:"new_name"`
}

type WedriveFileMoveRequest struct {
	FatherID string   `json:"fatherid"` // 当前目录的fileid,根目录时为空间spaceid
	Replace  bool     `json:"replace"`
	FileID   []string `json:"fileid"`
}

type WedriveFileMoveResponse struct {
	Error
	FileList WedriveFileList `json:"file_list"`
}

type WedriveFileUploadRequest struct {
	SpaceID    string `json:"spaceid"`
	FatherID   string `json:"fatherid"` // 当前目录的fileid,根目录时为空间spaceid
	FileName   string `json:"file_name"`
	DataBase64 string `json:"file_base64_content"`
}

type WedriveFileUploadResponse struct {
	Error
	FileID string `json:"fileid"`
}

type WedriveFileBlockUploadRequest struct {
	SpaceID      string `json:"spaceid"`
	FatherID     string `json:"fatherid"` // 当前目录的fileid,根目录时为空间spaceid
	FileName     string `json:"file_name"`
	Size         uint64 `json:"size"` //max 20G
	SkipPushCard bool   `json:"skip_push_card"`
	Data         io.ReadSeeker
	Concurrent   int
}

type blockUploadInitRequest struct {
	SpaceID      string   `json:"spaceid"`
	FatherID     string   `json:"fatherid"` // 当前目录的fileid,根目录时为空间spaceid
	FileName     string   `json:"file_name"`
	Size         uint64   `json:"size"` //max 20G
	SkipPushCard bool     `json:"skip_push_card"`
	BlockSha     []string `json:"block_sha"`
}

type blockUploadInitResponse struct {
	Error
	Hit    bool   `json:"hit_exist"`
	Key    string `json:"upload_key"`
	FileID string `json:"fileid"`
}

type blockUploadPartRequest struct {
	Key   string `json:"upload_key"`
	Index int32  `json:"index"`
	Data  string `json:"file_base64_content"`
}

type WedriveFileDownloadRequest WedriveFileInfoRequest

type WedriveFile struct {
	token
}

func (p *WedriveFile) List(param *WedriveFileListRequest) (ret WedriveFileListResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileList, param, &ret)

	return
}

func (p *WedriveFile) Info(param *WedriveFileInfoRequest) (ret WedriveFileInfoResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileInfo, param, &ret)

	return
}

func (p *WedriveFile) Create(param *WedriveFileCreateRequest) (ret WedriveFileCreateResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileCreate, param, &ret)

	return
}

func (p *WedriveFile) Delete(param *WedriveFileDeleteRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileDelete, param, &ret)

	return
}

func (p *WedriveFile) Rename(param *WedriveFileRenameRequest) (ret Error, err error) {
	err = wedrivePost(&p.token, wedriveApiFileRename, param, &ret)

	return
}

func (p *WedriveFile) Move(param *WedriveFileMoveRequest) (ret WedriveFileMoveResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileMove, param, &ret)

	return
}

// Small files upload. FileSize <= 10M.
func (p *WedriveFile) Upload(param *WedriveFileUploadRequest) (ret WedriveFileUploadResponse, err error) {
	err = wedrivePost(&p.token, wedriveApiFileUpload, param, &ret)

	return
}

// Big files upload.
func (p *WedriveFile) BlockUpload(param *WedriveFileBlockUploadRequest) (ret WedriveFileUploadResponse, err error) {
	var initRet blockUploadInitResponse
	if initRet, err = p.blockUploadInit(param); err != nil {
		return
	} else if initRet.Hit {
		ret.Error = initRet.Error
		ret.FileID = initRet.FileID

		return
	}

	var partRet Error
	blocks := wedriveFileBlocks(param.Size)
	if partRet, err = p.blockUploadPart(param.Data, initRet.Key, blocks, param.Concurrent); err != nil {
		return
	} else if partRet.ErrCode != 0 {
		ret.Error = partRet

		return
	}

	return p.blockUploadFinish(initRet.Key)
}

func wedriveFileBlocks(size uint64) int {
	return int((size-1)/uint64(wedriveBlockSize)) + 1
}

func (p *WedriveFile) blockUploadInit(param *WedriveFileBlockUploadRequest) (ret blockUploadInitResponse, err error) {
	blocks := wedriveFileBlocks(param.Size)
	sha := make([]string, blocks)

	req := blockUploadInitRequest{
		SpaceID:      param.SpaceID,
		FatherID:     param.FatherID,
		FileName:     param.FileName,
		Size:         param.Size,
		SkipPushCard: param.SkipPushCard,
		BlockSha:     sha,
	}

	data := make([]byte, wedriveBlockSize)

	r := param.Data
	r.Seek(0, io.SeekStart)
	h := sha1.New()

	for i := 0; i < blocks; i++ {
		var n int
		if n, err = r.Read(data); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				break
			}
		}

		h.Write(data[:n])
		if i == blocks-1 {
			sum := h.Sum(nil)
			sha[i] = hex.EncodeToString(sum)
		} else {
			sha[i] = getHashState(h)
		}
	}

	if err != nil {
		return
	}

	err = wedrivePost(&p.token, wedriveApiFileBlockUploadInit, &req, &ret)

	return
}

func getHashState(h hash.Hash) (ret string) {
	if m, ok := h.(encoding.BinaryMarshaler); ok {
		if binaryData, err := m.MarshalBinary(); err == nil {
			data := [...]uint32{
				binary.BigEndian.Uint32(binaryData[4:8]),
				binary.BigEndian.Uint32(binaryData[8:12]),
				binary.BigEndian.Uint32(binaryData[12:16]),
				binary.BigEndian.Uint32(binaryData[16:20]),
				binary.BigEndian.Uint32(binaryData[20:24]),
			}

			bytes := make([]byte, 0, 20)
			bytes = binary.LittleEndian.AppendUint32(bytes, data[0])
			bytes = binary.LittleEndian.AppendUint32(bytes, data[1])
			bytes = binary.LittleEndian.AppendUint32(bytes, data[2])
			bytes = binary.LittleEndian.AppendUint32(bytes, data[3])
			bytes = binary.LittleEndian.AppendUint32(bytes, data[4])

			ret = hex.EncodeToString(bytes)
		}
	}

	return
}

func (p *WedriveFile) blockUploadPart(r io.ReadSeeker, key string, blocks, concurrent int) (ret Error, err error) {
	r.Seek(0, io.SeekStart)

	if concurrent < 1 {
		concurrent = 1
	} else if concurrent > blocks {
		concurrent = blocks
	}

	if concurrent == 1 {
		ret, err = p.blockUploadPartSequent(r, key, blocks)
	} else {
		ret, err = p.blockUploadPartConcurrent(r, key, blocks, concurrent)
	}

	return
}

func (p *WedriveFile) blockUploadPartSequent(r io.ReadSeeker, key string, blocks int) (ret Error, err error) {
	data := make([]byte, wedriveBlockSize)

	for i := 1; i <= blocks; i++ {
		n, _ := r.Read(data)
		request := &blockUploadPartRequest{
			Key:   key,
			Index: int32(i),
			Data:  base64.StdEncoding.EncodeToString(data[:n]),
		}

		if err = wedrivePost(&p.token, wedriveApiFileBlockUploadPart, &request, &ret); err != nil || ret.ErrCode != 0 {
			break
		}
	}

	return
}

func (p *WedriveFile) blockUploadPartConcurrent(r io.ReadSeeker, key string, blocks, concurrent int) (ret Error, err error) {
	var mtx sync.Mutex
	var wg sync.WaitGroup
	current := 0
	errCount := int32(0)

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			data := make([]byte, wedriveBlockSize)

			for current < blocks && errCount == 0 {
				var index int32

				mtx.Lock()
				n, _ := r.Read(data)
				current++
				index = int32(current)
				mtx.Unlock()

				request := &blockUploadPartRequest{
					Key:   key,
					Index: index,
					Data:  base64.StdEncoding.EncodeToString(data[:n]),
				}

				var res Error
				if e := wedrivePost(&p.token, wedriveApiFileBlockUploadPart, &request, &res); e != nil || res.ErrCode != 0 {
					if atomic.AddInt32(&errCount, 1) == 1 {
						err = e
						ret = res
					}
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return
}

func (p *WedriveFile) blockUploadFinish(key string) (ret WedriveFileUploadResponse, err error) {
	type blockUploadFinishRequest struct {
		Key string `json:"upload_key"`
	}

	req := blockUploadFinishRequest{
		Key: key,
	}

	err = wedrivePost(&p.token, wedriveApiFileBlockUploadFinish, &req, &ret)

	return
}

func (p *WedriveFile) Download(param *WedriveFileDownloadRequest, to io.Writer) (ret Error, err error) {
	type fileDownloadResponse struct {
		Error
		Url         string `json:"download_url"`
		CookieName  string `json:"cookie_name"`
		CookieValue string `json:"cookie_value"`
	}

	var r fileDownloadResponse
	err = wedrivePost(&p.token, wedriveApiFileDownload, param, &r)
	ret = r.Error
	if err == nil && r.ErrCode == 0 {
		var req *http.Request

		if req, err = http.NewRequest("GET", r.Url, nil); err == nil {
			var client http.Client
			var res *http.Response

			req.AddCookie(&http.Cookie{Name: r.CookieName, Value: r.CookieValue})

			if res, err = client.Do(req); err == nil {
				defer res.Body.Close()
				io.Copy(to, res.Body)
			}
		}
	}

	return
}

const (
	wedriveApiFileList              = "file_list"
	wedriveApiFileInfo              = "file_info"
	wedriveApiFileCreate            = "file_create"
	wedriveApiFileDelete            = "file_delete"
	wedriveApiFileRename            = "file_rename"
	wedriveApiFileMove              = "file_move"
	wedriveApiFileUpload            = "file_upload"
	wedriveApiFileBlockUploadInit   = "file_upload_init"
	wedriveApiFileBlockUploadPart   = "file_upload_part"
	wedriveApiFileBlockUploadFinish = "file_upload_finish"
	wedriveApiFileDownload          = "file_download"
)

const (
	wedriveBlockSize int = 2 * 1024 * 1024
)
