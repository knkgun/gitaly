// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wiki.proto

package gitaly

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type WikiCommitDetails struct {
	Name    []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email   []byte `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Message []byte `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (m *WikiCommitDetails) Reset()                    { *m = WikiCommitDetails{} }
func (m *WikiCommitDetails) String() string            { return proto.CompactTextString(m) }
func (*WikiCommitDetails) ProtoMessage()               {}
func (*WikiCommitDetails) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{0} }

func (m *WikiCommitDetails) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *WikiCommitDetails) GetEmail() []byte {
	if m != nil {
		return m.Email
	}
	return nil
}

func (m *WikiCommitDetails) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

type WikiPageVersion struct {
	Commit *GitCommit `protobuf:"bytes,1,opt,name=commit" json:"commit,omitempty"`
	Format string     `protobuf:"bytes,2,opt,name=format" json:"format,omitempty"`
}

func (m *WikiPageVersion) Reset()                    { *m = WikiPageVersion{} }
func (m *WikiPageVersion) String() string            { return proto.CompactTextString(m) }
func (*WikiPageVersion) ProtoMessage()               {}
func (*WikiPageVersion) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{1} }

func (m *WikiPageVersion) GetCommit() *GitCommit {
	if m != nil {
		return m.Commit
	}
	return nil
}

func (m *WikiPageVersion) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

type WikiPage struct {
	// These fields are only present in the first message of a WikiPage stream
	Version    *WikiPageVersion `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	Format     string           `protobuf:"bytes,2,opt,name=format" json:"format,omitempty"`
	Title      []byte           `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	UrlPath    string           `protobuf:"bytes,4,opt,name=url_path,json=urlPath" json:"url_path,omitempty"`
	Path       []byte           `protobuf:"bytes,5,opt,name=path,proto3" json:"path,omitempty"`
	Name       []byte           `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Historical bool             `protobuf:"varint,7,opt,name=historical" json:"historical,omitempty"`
	// This field is present in all messages of a WikiPage stream
	RawData []byte `protobuf:"bytes,8,opt,name=raw_data,json=rawData,proto3" json:"raw_data,omitempty"`
}

func (m *WikiPage) Reset()                    { *m = WikiPage{} }
func (m *WikiPage) String() string            { return proto.CompactTextString(m) }
func (*WikiPage) ProtoMessage()               {}
func (*WikiPage) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{2} }

func (m *WikiPage) GetVersion() *WikiPageVersion {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *WikiPage) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *WikiPage) GetTitle() []byte {
	if m != nil {
		return m.Title
	}
	return nil
}

func (m *WikiPage) GetUrlPath() string {
	if m != nil {
		return m.UrlPath
	}
	return ""
}

func (m *WikiPage) GetPath() []byte {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *WikiPage) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *WikiPage) GetHistorical() bool {
	if m != nil {
		return m.Historical
	}
	return false
}

func (m *WikiPage) GetRawData() []byte {
	if m != nil {
		return m.RawData
	}
	return nil
}

type WikiGetPageVersionsRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	PagePath   []byte      `protobuf:"bytes,2,opt,name=page_path,json=pagePath,proto3" json:"page_path,omitempty"`
	Page       int32       `protobuf:"varint,3,opt,name=page" json:"page,omitempty"`
	PerPage    int32       `protobuf:"varint,4,opt,name=per_page,json=perPage" json:"per_page,omitempty"`
}

func (m *WikiGetPageVersionsRequest) Reset()                    { *m = WikiGetPageVersionsRequest{} }
func (m *WikiGetPageVersionsRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiGetPageVersionsRequest) ProtoMessage()               {}
func (*WikiGetPageVersionsRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{3} }

func (m *WikiGetPageVersionsRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiGetPageVersionsRequest) GetPagePath() []byte {
	if m != nil {
		return m.PagePath
	}
	return nil
}

func (m *WikiGetPageVersionsRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *WikiGetPageVersionsRequest) GetPerPage() int32 {
	if m != nil {
		return m.PerPage
	}
	return 0
}

type WikiGetPageVersionsResponse struct {
	Versions []*WikiPageVersion `protobuf:"bytes,1,rep,name=versions" json:"versions,omitempty"`
}

func (m *WikiGetPageVersionsResponse) Reset()                    { *m = WikiGetPageVersionsResponse{} }
func (m *WikiGetPageVersionsResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiGetPageVersionsResponse) ProtoMessage()               {}
func (*WikiGetPageVersionsResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{4} }

func (m *WikiGetPageVersionsResponse) GetVersions() []*WikiPageVersion {
	if m != nil {
		return m.Versions
	}
	return nil
}

// This message is sent in a stream because the 'content' field may be large.
type WikiWritePageRequest struct {
	// These following fields are only present in the first message.
	Repository    *Repository        `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	Name          []byte             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Format        string             `protobuf:"bytes,3,opt,name=format" json:"format,omitempty"`
	CommitDetails *WikiCommitDetails `protobuf:"bytes,4,opt,name=commit_details,json=commitDetails" json:"commit_details,omitempty"`
	// This field is present in all messages.
	Content []byte `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *WikiWritePageRequest) Reset()                    { *m = WikiWritePageRequest{} }
func (m *WikiWritePageRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiWritePageRequest) ProtoMessage()               {}
func (*WikiWritePageRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{5} }

func (m *WikiWritePageRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiWritePageRequest) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *WikiWritePageRequest) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *WikiWritePageRequest) GetCommitDetails() *WikiCommitDetails {
	if m != nil {
		return m.CommitDetails
	}
	return nil
}

func (m *WikiWritePageRequest) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type WikiWritePageResponse struct {
	DuplicateError []byte `protobuf:"bytes,1,opt,name=duplicate_error,json=duplicateError,proto3" json:"duplicate_error,omitempty"`
}

func (m *WikiWritePageResponse) Reset()                    { *m = WikiWritePageResponse{} }
func (m *WikiWritePageResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiWritePageResponse) ProtoMessage()               {}
func (*WikiWritePageResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{6} }

func (m *WikiWritePageResponse) GetDuplicateError() []byte {
	if m != nil {
		return m.DuplicateError
	}
	return nil
}

type WikiUpdatePageRequest struct {
	// There fields are only present in the first message of the stream
	Repository    *Repository        `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	PagePath      []byte             `protobuf:"bytes,2,opt,name=page_path,json=pagePath,proto3" json:"page_path,omitempty"`
	Title         []byte             `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Format        string             `protobuf:"bytes,4,opt,name=format" json:"format,omitempty"`
	CommitDetails *WikiCommitDetails `protobuf:"bytes,5,opt,name=commit_details,json=commitDetails" json:"commit_details,omitempty"`
	// This field is present in all messages
	Content []byte `protobuf:"bytes,6,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *WikiUpdatePageRequest) Reset()                    { *m = WikiUpdatePageRequest{} }
func (m *WikiUpdatePageRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiUpdatePageRequest) ProtoMessage()               {}
func (*WikiUpdatePageRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{7} }

func (m *WikiUpdatePageRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiUpdatePageRequest) GetPagePath() []byte {
	if m != nil {
		return m.PagePath
	}
	return nil
}

func (m *WikiUpdatePageRequest) GetTitle() []byte {
	if m != nil {
		return m.Title
	}
	return nil
}

func (m *WikiUpdatePageRequest) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func (m *WikiUpdatePageRequest) GetCommitDetails() *WikiCommitDetails {
	if m != nil {
		return m.CommitDetails
	}
	return nil
}

func (m *WikiUpdatePageRequest) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type WikiUpdatePageResponse struct {
	Error []byte `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *WikiUpdatePageResponse) Reset()                    { *m = WikiUpdatePageResponse{} }
func (m *WikiUpdatePageResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiUpdatePageResponse) ProtoMessage()               {}
func (*WikiUpdatePageResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{8} }

func (m *WikiUpdatePageResponse) GetError() []byte {
	if m != nil {
		return m.Error
	}
	return nil
}

type WikiDeletePageRequest struct {
	Repository    *Repository        `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	PagePath      []byte             `protobuf:"bytes,2,opt,name=page_path,json=pagePath,proto3" json:"page_path,omitempty"`
	CommitDetails *WikiCommitDetails `protobuf:"bytes,3,opt,name=commit_details,json=commitDetails" json:"commit_details,omitempty"`
}

func (m *WikiDeletePageRequest) Reset()                    { *m = WikiDeletePageRequest{} }
func (m *WikiDeletePageRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiDeletePageRequest) ProtoMessage()               {}
func (*WikiDeletePageRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{9} }

func (m *WikiDeletePageRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiDeletePageRequest) GetPagePath() []byte {
	if m != nil {
		return m.PagePath
	}
	return nil
}

func (m *WikiDeletePageRequest) GetCommitDetails() *WikiCommitDetails {
	if m != nil {
		return m.CommitDetails
	}
	return nil
}

type WikiDeletePageResponse struct {
}

func (m *WikiDeletePageResponse) Reset()                    { *m = WikiDeletePageResponse{} }
func (m *WikiDeletePageResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiDeletePageResponse) ProtoMessage()               {}
func (*WikiDeletePageResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{10} }

type WikiFindPageRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	Title      []byte      `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Revision   []byte      `protobuf:"bytes,3,opt,name=revision,proto3" json:"revision,omitempty"`
	Directory  []byte      `protobuf:"bytes,4,opt,name=directory,proto3" json:"directory,omitempty"`
}

func (m *WikiFindPageRequest) Reset()                    { *m = WikiFindPageRequest{} }
func (m *WikiFindPageRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiFindPageRequest) ProtoMessage()               {}
func (*WikiFindPageRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{11} }

func (m *WikiFindPageRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiFindPageRequest) GetTitle() []byte {
	if m != nil {
		return m.Title
	}
	return nil
}

func (m *WikiFindPageRequest) GetRevision() []byte {
	if m != nil {
		return m.Revision
	}
	return nil
}

func (m *WikiFindPageRequest) GetDirectory() []byte {
	if m != nil {
		return m.Directory
	}
	return nil
}

// WikiFindPageResponse is a stream because we need multiple WikiPage
// messages to send the raw_data field.
type WikiFindPageResponse struct {
	Page *WikiPage `protobuf:"bytes,1,opt,name=page" json:"page,omitempty"`
}

func (m *WikiFindPageResponse) Reset()                    { *m = WikiFindPageResponse{} }
func (m *WikiFindPageResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiFindPageResponse) ProtoMessage()               {}
func (*WikiFindPageResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{12} }

func (m *WikiFindPageResponse) GetPage() *WikiPage {
	if m != nil {
		return m.Page
	}
	return nil
}

type WikiFindFileRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	Name       []byte      `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Optional: revision
	Revision []byte `protobuf:"bytes,3,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (m *WikiFindFileRequest) Reset()                    { *m = WikiFindFileRequest{} }
func (m *WikiFindFileRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiFindFileRequest) ProtoMessage()               {}
func (*WikiFindFileRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{13} }

func (m *WikiFindFileRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *WikiFindFileRequest) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *WikiFindFileRequest) GetRevision() []byte {
	if m != nil {
		return m.Revision
	}
	return nil
}

type WikiFindFileResponse struct {
	// If 'name' is empty, the file was not found.
	Name     []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	MimeType string `protobuf:"bytes,2,opt,name=mime_type,json=mimeType" json:"mime_type,omitempty"`
	RawData  []byte `protobuf:"bytes,3,opt,name=raw_data,json=rawData,proto3" json:"raw_data,omitempty"`
	Path     []byte `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
}

func (m *WikiFindFileResponse) Reset()                    { *m = WikiFindFileResponse{} }
func (m *WikiFindFileResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiFindFileResponse) ProtoMessage()               {}
func (*WikiFindFileResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{14} }

func (m *WikiFindFileResponse) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *WikiFindFileResponse) GetMimeType() string {
	if m != nil {
		return m.MimeType
	}
	return ""
}

func (m *WikiFindFileResponse) GetRawData() []byte {
	if m != nil {
		return m.RawData
	}
	return nil
}

func (m *WikiFindFileResponse) GetPath() []byte {
	if m != nil {
		return m.Path
	}
	return nil
}

type WikiGetAllPagesRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
}

func (m *WikiGetAllPagesRequest) Reset()                    { *m = WikiGetAllPagesRequest{} }
func (m *WikiGetAllPagesRequest) String() string            { return proto.CompactTextString(m) }
func (*WikiGetAllPagesRequest) ProtoMessage()               {}
func (*WikiGetAllPagesRequest) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{15} }

func (m *WikiGetAllPagesRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

// The WikiGetAllPagesResponse stream is a concatenation of WikiPage streams
type WikiGetAllPagesResponse struct {
	Page *WikiPage `protobuf:"bytes,1,opt,name=page" json:"page,omitempty"`
	// When end_of_page is true it signals a change of page for the next Response message (if any)
	EndOfPage bool `protobuf:"varint,2,opt,name=end_of_page,json=endOfPage" json:"end_of_page,omitempty"`
}

func (m *WikiGetAllPagesResponse) Reset()                    { *m = WikiGetAllPagesResponse{} }
func (m *WikiGetAllPagesResponse) String() string            { return proto.CompactTextString(m) }
func (*WikiGetAllPagesResponse) ProtoMessage()               {}
func (*WikiGetAllPagesResponse) Descriptor() ([]byte, []int) { return fileDescriptor15, []int{16} }

func (m *WikiGetAllPagesResponse) GetPage() *WikiPage {
	if m != nil {
		return m.Page
	}
	return nil
}

func (m *WikiGetAllPagesResponse) GetEndOfPage() bool {
	if m != nil {
		return m.EndOfPage
	}
	return false
}

func init() {
	proto.RegisterType((*WikiCommitDetails)(nil), "gitaly.WikiCommitDetails")
	proto.RegisterType((*WikiPageVersion)(nil), "gitaly.WikiPageVersion")
	proto.RegisterType((*WikiPage)(nil), "gitaly.WikiPage")
	proto.RegisterType((*WikiGetPageVersionsRequest)(nil), "gitaly.WikiGetPageVersionsRequest")
	proto.RegisterType((*WikiGetPageVersionsResponse)(nil), "gitaly.WikiGetPageVersionsResponse")
	proto.RegisterType((*WikiWritePageRequest)(nil), "gitaly.WikiWritePageRequest")
	proto.RegisterType((*WikiWritePageResponse)(nil), "gitaly.WikiWritePageResponse")
	proto.RegisterType((*WikiUpdatePageRequest)(nil), "gitaly.WikiUpdatePageRequest")
	proto.RegisterType((*WikiUpdatePageResponse)(nil), "gitaly.WikiUpdatePageResponse")
	proto.RegisterType((*WikiDeletePageRequest)(nil), "gitaly.WikiDeletePageRequest")
	proto.RegisterType((*WikiDeletePageResponse)(nil), "gitaly.WikiDeletePageResponse")
	proto.RegisterType((*WikiFindPageRequest)(nil), "gitaly.WikiFindPageRequest")
	proto.RegisterType((*WikiFindPageResponse)(nil), "gitaly.WikiFindPageResponse")
	proto.RegisterType((*WikiFindFileRequest)(nil), "gitaly.WikiFindFileRequest")
	proto.RegisterType((*WikiFindFileResponse)(nil), "gitaly.WikiFindFileResponse")
	proto.RegisterType((*WikiGetAllPagesRequest)(nil), "gitaly.WikiGetAllPagesRequest")
	proto.RegisterType((*WikiGetAllPagesResponse)(nil), "gitaly.WikiGetAllPagesResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for WikiService service

type WikiServiceClient interface {
	WikiGetPageVersions(ctx context.Context, in *WikiGetPageVersionsRequest, opts ...grpc.CallOption) (WikiService_WikiGetPageVersionsClient, error)
	WikiWritePage(ctx context.Context, opts ...grpc.CallOption) (WikiService_WikiWritePageClient, error)
	WikiUpdatePage(ctx context.Context, opts ...grpc.CallOption) (WikiService_WikiUpdatePageClient, error)
	WikiDeletePage(ctx context.Context, in *WikiDeletePageRequest, opts ...grpc.CallOption) (*WikiDeletePageResponse, error)
	// WikiFindPage returns a stream because the page's raw_data field may be arbitrarily large.
	WikiFindPage(ctx context.Context, in *WikiFindPageRequest, opts ...grpc.CallOption) (WikiService_WikiFindPageClient, error)
	WikiFindFile(ctx context.Context, in *WikiFindFileRequest, opts ...grpc.CallOption) (WikiService_WikiFindFileClient, error)
	WikiGetAllPages(ctx context.Context, in *WikiGetAllPagesRequest, opts ...grpc.CallOption) (WikiService_WikiGetAllPagesClient, error)
}

type wikiServiceClient struct {
	cc *grpc.ClientConn
}

func NewWikiServiceClient(cc *grpc.ClientConn) WikiServiceClient {
	return &wikiServiceClient{cc}
}

func (c *wikiServiceClient) WikiGetPageVersions(ctx context.Context, in *WikiGetPageVersionsRequest, opts ...grpc.CallOption) (WikiService_WikiGetPageVersionsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[0], c.cc, "/gitaly.WikiService/WikiGetPageVersions", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiGetPageVersionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WikiService_WikiGetPageVersionsClient interface {
	Recv() (*WikiGetPageVersionsResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiGetPageVersionsClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiGetPageVersionsClient) Recv() (*WikiGetPageVersionsResponse, error) {
	m := new(WikiGetPageVersionsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wikiServiceClient) WikiWritePage(ctx context.Context, opts ...grpc.CallOption) (WikiService_WikiWritePageClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[1], c.cc, "/gitaly.WikiService/WikiWritePage", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiWritePageClient{stream}
	return x, nil
}

type WikiService_WikiWritePageClient interface {
	Send(*WikiWritePageRequest) error
	CloseAndRecv() (*WikiWritePageResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiWritePageClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiWritePageClient) Send(m *WikiWritePageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *wikiServiceWikiWritePageClient) CloseAndRecv() (*WikiWritePageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(WikiWritePageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wikiServiceClient) WikiUpdatePage(ctx context.Context, opts ...grpc.CallOption) (WikiService_WikiUpdatePageClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[2], c.cc, "/gitaly.WikiService/WikiUpdatePage", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiUpdatePageClient{stream}
	return x, nil
}

type WikiService_WikiUpdatePageClient interface {
	Send(*WikiUpdatePageRequest) error
	CloseAndRecv() (*WikiUpdatePageResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiUpdatePageClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiUpdatePageClient) Send(m *WikiUpdatePageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *wikiServiceWikiUpdatePageClient) CloseAndRecv() (*WikiUpdatePageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(WikiUpdatePageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wikiServiceClient) WikiDeletePage(ctx context.Context, in *WikiDeletePageRequest, opts ...grpc.CallOption) (*WikiDeletePageResponse, error) {
	out := new(WikiDeletePageResponse)
	err := grpc.Invoke(ctx, "/gitaly.WikiService/WikiDeletePage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wikiServiceClient) WikiFindPage(ctx context.Context, in *WikiFindPageRequest, opts ...grpc.CallOption) (WikiService_WikiFindPageClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[3], c.cc, "/gitaly.WikiService/WikiFindPage", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiFindPageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WikiService_WikiFindPageClient interface {
	Recv() (*WikiFindPageResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiFindPageClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiFindPageClient) Recv() (*WikiFindPageResponse, error) {
	m := new(WikiFindPageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wikiServiceClient) WikiFindFile(ctx context.Context, in *WikiFindFileRequest, opts ...grpc.CallOption) (WikiService_WikiFindFileClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[4], c.cc, "/gitaly.WikiService/WikiFindFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiFindFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WikiService_WikiFindFileClient interface {
	Recv() (*WikiFindFileResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiFindFileClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiFindFileClient) Recv() (*WikiFindFileResponse, error) {
	m := new(WikiFindFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wikiServiceClient) WikiGetAllPages(ctx context.Context, in *WikiGetAllPagesRequest, opts ...grpc.CallOption) (WikiService_WikiGetAllPagesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WikiService_serviceDesc.Streams[5], c.cc, "/gitaly.WikiService/WikiGetAllPages", opts...)
	if err != nil {
		return nil, err
	}
	x := &wikiServiceWikiGetAllPagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WikiService_WikiGetAllPagesClient interface {
	Recv() (*WikiGetAllPagesResponse, error)
	grpc.ClientStream
}

type wikiServiceWikiGetAllPagesClient struct {
	grpc.ClientStream
}

func (x *wikiServiceWikiGetAllPagesClient) Recv() (*WikiGetAllPagesResponse, error) {
	m := new(WikiGetAllPagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for WikiService service

type WikiServiceServer interface {
	WikiGetPageVersions(*WikiGetPageVersionsRequest, WikiService_WikiGetPageVersionsServer) error
	WikiWritePage(WikiService_WikiWritePageServer) error
	WikiUpdatePage(WikiService_WikiUpdatePageServer) error
	WikiDeletePage(context.Context, *WikiDeletePageRequest) (*WikiDeletePageResponse, error)
	// WikiFindPage returns a stream because the page's raw_data field may be arbitrarily large.
	WikiFindPage(*WikiFindPageRequest, WikiService_WikiFindPageServer) error
	WikiFindFile(*WikiFindFileRequest, WikiService_WikiFindFileServer) error
	WikiGetAllPages(*WikiGetAllPagesRequest, WikiService_WikiGetAllPagesServer) error
}

func RegisterWikiServiceServer(s *grpc.Server, srv WikiServiceServer) {
	s.RegisterService(&_WikiService_serviceDesc, srv)
}

func _WikiService_WikiGetPageVersions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WikiGetPageVersionsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WikiServiceServer).WikiGetPageVersions(m, &wikiServiceWikiGetPageVersionsServer{stream})
}

type WikiService_WikiGetPageVersionsServer interface {
	Send(*WikiGetPageVersionsResponse) error
	grpc.ServerStream
}

type wikiServiceWikiGetPageVersionsServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiGetPageVersionsServer) Send(m *WikiGetPageVersionsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WikiService_WikiWritePage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WikiServiceServer).WikiWritePage(&wikiServiceWikiWritePageServer{stream})
}

type WikiService_WikiWritePageServer interface {
	SendAndClose(*WikiWritePageResponse) error
	Recv() (*WikiWritePageRequest, error)
	grpc.ServerStream
}

type wikiServiceWikiWritePageServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiWritePageServer) SendAndClose(m *WikiWritePageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *wikiServiceWikiWritePageServer) Recv() (*WikiWritePageRequest, error) {
	m := new(WikiWritePageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _WikiService_WikiUpdatePage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WikiServiceServer).WikiUpdatePage(&wikiServiceWikiUpdatePageServer{stream})
}

type WikiService_WikiUpdatePageServer interface {
	SendAndClose(*WikiUpdatePageResponse) error
	Recv() (*WikiUpdatePageRequest, error)
	grpc.ServerStream
}

type wikiServiceWikiUpdatePageServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiUpdatePageServer) SendAndClose(m *WikiUpdatePageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *wikiServiceWikiUpdatePageServer) Recv() (*WikiUpdatePageRequest, error) {
	m := new(WikiUpdatePageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _WikiService_WikiDeletePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WikiDeletePageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WikiServiceServer).WikiDeletePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.WikiService/WikiDeletePage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WikiServiceServer).WikiDeletePage(ctx, req.(*WikiDeletePageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WikiService_WikiFindPage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WikiFindPageRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WikiServiceServer).WikiFindPage(m, &wikiServiceWikiFindPageServer{stream})
}

type WikiService_WikiFindPageServer interface {
	Send(*WikiFindPageResponse) error
	grpc.ServerStream
}

type wikiServiceWikiFindPageServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiFindPageServer) Send(m *WikiFindPageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WikiService_WikiFindFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WikiFindFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WikiServiceServer).WikiFindFile(m, &wikiServiceWikiFindFileServer{stream})
}

type WikiService_WikiFindFileServer interface {
	Send(*WikiFindFileResponse) error
	grpc.ServerStream
}

type wikiServiceWikiFindFileServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiFindFileServer) Send(m *WikiFindFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WikiService_WikiGetAllPages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WikiGetAllPagesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WikiServiceServer).WikiGetAllPages(m, &wikiServiceWikiGetAllPagesServer{stream})
}

type WikiService_WikiGetAllPagesServer interface {
	Send(*WikiGetAllPagesResponse) error
	grpc.ServerStream
}

type wikiServiceWikiGetAllPagesServer struct {
	grpc.ServerStream
}

func (x *wikiServiceWikiGetAllPagesServer) Send(m *WikiGetAllPagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _WikiService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.WikiService",
	HandlerType: (*WikiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WikiDeletePage",
			Handler:    _WikiService_WikiDeletePage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WikiGetPageVersions",
			Handler:       _WikiService_WikiGetPageVersions_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "WikiWritePage",
			Handler:       _WikiService_WikiWritePage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "WikiUpdatePage",
			Handler:       _WikiService_WikiUpdatePage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "WikiFindPage",
			Handler:       _WikiService_WikiFindPage_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "WikiFindFile",
			Handler:       _WikiService_WikiFindFile_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "WikiGetAllPages",
			Handler:       _WikiService_WikiGetAllPages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "wiki.proto",
}

func init() { proto.RegisterFile("wiki.proto", fileDescriptor15) }

var fileDescriptor15 = []byte{
	// 846 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xcd, 0x72, 0xdc, 0x44,
	0x10, 0x8e, 0xbc, 0xeb, 0x5d, 0x6d, 0xdb, 0x71, 0xc8, 0x10, 0x12, 0x45, 0x36, 0xc6, 0x35, 0x50,
	0x85, 0xb9, 0xb8, 0xc0, 0xb9, 0x72, 0x08, 0x85, 0x89, 0x2f, 0x50, 0x18, 0x25, 0xc4, 0x47, 0xd5,
	0x64, 0xd5, 0xf6, 0x4e, 0x45, 0x7f, 0x8c, 0x66, 0xd7, 0xb5, 0x8f, 0x42, 0x15, 0x4f, 0xc0, 0xe3,
	0xf0, 0x06, 0x1c, 0xb9, 0xf2, 0x04, 0xd4, 0xfc, 0x68, 0x35, 0xd2, 0xfe, 0x50, 0x61, 0xc9, 0x4d,
	0xd3, 0xdd, 0xf3, 0x4d, 0x7f, 0x5f, 0x4f, 0xf7, 0x08, 0xe0, 0x8e, 0xbf, 0xe5, 0x67, 0xa5, 0x28,
	0x64, 0x41, 0x06, 0xb7, 0x5c, 0xb2, 0x74, 0x1e, 0xee, 0x57, 0x13, 0x26, 0x30, 0x31, 0x56, 0x7a,
	0x0d, 0x0f, 0xaf, 0xf9, 0x5b, 0xfe, 0x6d, 0x91, 0x65, 0x5c, 0x5e, 0xa0, 0x64, 0x3c, 0xad, 0x08,
	0x81, 0x7e, 0xce, 0x32, 0x0c, 0xbc, 0x13, 0xef, 0x74, 0x3f, 0xd2, 0xdf, 0xe4, 0x11, 0xec, 0x62,
	0xc6, 0x78, 0x1a, 0xec, 0x68, 0xa3, 0x59, 0x90, 0x00, 0x86, 0x19, 0x56, 0x15, 0xbb, 0xc5, 0xa0,
	0xa7, 0xed, 0xf5, 0x92, 0xbe, 0x82, 0x07, 0x0a, 0xf8, 0x8a, 0xdd, 0xe2, 0x6b, 0x14, 0x15, 0x2f,
	0x72, 0xf2, 0x05, 0x0c, 0xc6, 0xfa, 0x1c, 0x0d, 0xbc, 0x77, 0xfe, 0xf0, 0xcc, 0xa4, 0x74, 0x76,
	0xc9, 0xa5, 0x49, 0x20, 0xb2, 0x01, 0xe4, 0x31, 0x0c, 0x6e, 0x0a, 0x91, 0x31, 0xa9, 0x8f, 0x1b,
	0x45, 0x76, 0x45, 0xff, 0xf2, 0xc0, 0xaf, 0x61, 0xc9, 0x57, 0x30, 0x9c, 0x19, 0x68, 0x0b, 0xf8,
	0xa4, 0x06, 0xec, 0x9c, 0x1c, 0xd5, 0x71, 0xeb, 0x70, 0x15, 0x3b, 0xc9, 0x65, 0x5a, 0xb3, 0x30,
	0x0b, 0xf2, 0x14, 0xfc, 0xa9, 0x48, 0xe3, 0x92, 0xc9, 0x49, 0xd0, 0xd7, 0xf1, 0xc3, 0xa9, 0x48,
	0xaf, 0x98, 0x9c, 0x28, 0x89, 0xb4, 0x79, 0xd7, 0x48, 0x54, 0x5a, 0x9b, 0x96, 0x6d, 0xe0, 0xc8,
	0x76, 0x0c, 0x30, 0xe1, 0x95, 0x2c, 0x04, 0x1f, 0xb3, 0x34, 0x18, 0x9e, 0x78, 0xa7, 0x7e, 0xe4,
	0x58, 0xd4, 0x11, 0x82, 0xdd, 0xc5, 0x09, 0x93, 0x2c, 0xf0, 0x8d, 0x82, 0x82, 0xdd, 0x5d, 0x30,
	0xc9, 0xe8, 0x6f, 0x1e, 0x84, 0x8a, 0xc8, 0x25, 0x4a, 0x87, 0x4b, 0x15, 0xe1, 0x2f, 0x53, 0xac,
	0x24, 0x39, 0x07, 0x10, 0x58, 0x16, 0x15, 0x97, 0x85, 0x98, 0x5b, 0x01, 0x48, 0x2d, 0x40, 0xb4,
	0xf0, 0x44, 0x4e, 0x14, 0x39, 0x84, 0x51, 0xc9, 0x6e, 0xd1, 0x30, 0x32, 0x85, 0xf4, 0x95, 0xa1,
	0xa1, 0x64, 0x0b, 0xb9, 0x1b, 0xe9, 0x6f, 0x95, 0x5e, 0x89, 0x22, 0xd6, 0xf6, 0xbe, 0xb6, 0x0f,
	0x4b, 0x14, 0x2a, 0x1d, 0x1a, 0xc1, 0xe1, 0xca, 0xec, 0xaa, 0xb2, 0xc8, 0x2b, 0x24, 0xcf, 0xc0,
	0xb7, 0xa2, 0x57, 0x81, 0x77, 0xd2, 0xdb, 0x54, 0x9d, 0x45, 0x20, 0xfd, 0xc3, 0x83, 0x47, 0xca,
	0x7b, 0x2d, 0xb8, 0x44, 0x15, 0xb2, 0x0d, 0xd9, 0xba, 0x1c, 0x3b, 0x4e, 0x39, 0x9a, 0xfa, 0xf7,
	0x5a, 0xf5, 0x7f, 0x0e, 0x07, 0xe6, 0xe6, 0xc5, 0x89, 0xe9, 0x01, 0xcd, 0x76, 0xef, 0xfc, 0xa9,
	0x9b, 0x73, 0xab, 0x49, 0xa2, 0xfb, 0xe3, 0x56, 0xcf, 0x04, 0x30, 0x1c, 0x17, 0xb9, 0xc4, 0x5c,
	0xda, 0x3b, 0x51, 0x2f, 0xe9, 0x73, 0xf8, 0xa8, 0xc3, 0xc9, 0x4a, 0xf4, 0x39, 0x3c, 0x48, 0xa6,
	0x65, 0xca, 0xc7, 0x4c, 0x62, 0x8c, 0x42, 0x14, 0xc2, 0x76, 0xdc, 0xc1, 0xc2, 0xfc, 0x9d, 0xb2,
	0xd2, 0xbf, 0x3d, 0x03, 0xf1, 0x73, 0x99, 0xb0, 0xed, 0x75, 0xd9, 0x78, 0x09, 0x56, 0x37, 0x42,
	0x23, 0x5b, 0xff, 0x5f, 0x64, 0xdb, 0xfd, 0xef, 0xb2, 0x0d, 0xda, 0xb2, 0x9d, 0xc1, 0xe3, 0x2e,
	0x67, 0xab, 0x9b, 0x1a, 0x45, 0x8e, 0x5a, 0x66, 0x41, 0x7f, 0xb7, 0x22, 0x5d, 0x60, 0x8a, 0xef,
	0x59, 0xa4, 0x65, 0xda, 0xbd, 0x77, 0xa3, 0x4d, 0x03, 0x43, 0xce, 0xcd, 0xd5, 0x90, 0xa3, 0xbf,
	0x7a, 0xf0, 0xa1, 0x72, 0xbd, 0xe0, 0x79, 0xb2, 0x2d, 0x89, 0x45, 0x31, 0x77, 0xdc, 0x62, 0x86,
	0xe0, 0x0b, 0x9c, 0x71, 0x3d, 0x37, 0x4d, 0x95, 0x17, 0x6b, 0x72, 0x04, 0xa3, 0x84, 0x0b, 0x1c,
	0xeb, 0x43, 0xfa, 0xda, 0xd9, 0x18, 0xe8, 0xd7, 0xa6, 0x3b, 0x9b, 0xd4, 0x6c, 0x41, 0x3e, 0xb3,
	0x93, 0xc3, 0x64, 0xf5, 0x41, 0xb7, 0xcf, 0xcd, 0x2c, 0xa1, 0xf3, 0x86, 0xd8, 0x0b, 0x9e, 0xfe,
	0xef, 0xad, 0xbd, 0x81, 0x16, 0x9d, 0x35, 0x89, 0x9b, 0xa3, 0x6d, 0xe2, 0xab, 0x1e, 0xba, 0x43,
	0x18, 0x65, 0x3c, 0xc3, 0x58, 0xce, 0x4b, 0xb4, 0xaf, 0x84, 0xaf, 0x0c, 0xaf, 0xe6, 0x25, 0xb6,
	0xc6, 0x75, 0xaf, 0x35, 0xae, 0x17, 0x2f, 0x42, 0xbf, 0x79, 0x11, 0xe8, 0xf7, 0xa6, 0xcc, 0x97,
	0x28, 0xbf, 0x49, 0x53, 0x25, 0xc5, 0x36, 0xd3, 0x9b, 0xc6, 0xf0, 0x64, 0x09, 0xed, 0x5d, 0x2a,
	0x40, 0x8e, 0x61, 0x0f, 0xf3, 0x24, 0x2e, 0x6e, 0xcc, 0x40, 0xdf, 0xd1, 0xaf, 0xd1, 0x08, 0xf3,
	0xe4, 0xc7, 0x1b, 0x15, 0x75, 0xfe, 0x67, 0x1f, 0xf6, 0xd4, 0x96, 0x97, 0x28, 0x66, 0x7c, 0x8c,
	0xe4, 0x8d, 0xa9, 0x58, 0x67, 0xc4, 0x13, 0xea, 0xc2, 0xaf, 0x7e, 0x9d, 0xc2, 0x4f, 0x37, 0xc6,
	0xd8, 0xbb, 0x7e, 0xef, 0x4b, 0x8f, 0x5c, 0xc1, 0xfd, 0xd6, 0x74, 0x24, 0x47, 0xee, 0xce, 0xee,
	0x43, 0x10, 0x7e, 0xbc, 0xc6, 0x5b, 0x23, 0x9e, 0x7a, 0xe4, 0x25, 0x1c, 0xb4, 0x07, 0x07, 0x69,
	0x6d, 0x5a, 0x1a, 0xa2, 0xe1, 0xf1, 0x3a, 0xb7, 0x03, 0xfa, 0x93, 0x01, 0x6d, 0x1a, 0xb6, 0x0d,
	0xba, 0x34, 0x74, 0xda, 0xa0, 0x2b, 0xfa, 0xfc, 0x1e, 0xf9, 0x01, 0xf6, 0xdd, 0x6e, 0x22, 0x87,
	0xee, 0x8e, 0x4e, 0xfb, 0x87, 0x47, 0xab, 0x9d, 0x8e, 0x90, 0x0e, 0x9c, 0xba, 0xe3, 0xcb, 0x70,
	0x4e, 0xd3, 0x2d, 0xc3, 0xb9, 0x6d, 0xa1, 0xe1, 0x5e, 0x9b, 0xff, 0x37, 0xe7, 0xb2, 0x91, 0xe3,
	0x4e, 0x4d, 0x3b, 0x77, 0x3a, 0xfc, 0x64, 0xad, 0xbf, 0xc1, 0x7d, 0x33, 0xd0, 0xff, 0x9d, 0xcf,
	0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0xc9, 0xbb, 0xf3, 0xde, 0x9b, 0x0a, 0x00, 0x00,
}
