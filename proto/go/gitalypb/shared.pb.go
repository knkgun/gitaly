// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shared.proto

package gitalypb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ObjectType int32

const (
	ObjectType_UNKNOWN ObjectType = 0
	ObjectType_COMMIT  ObjectType = 1
	ObjectType_BLOB    ObjectType = 2
	ObjectType_TREE    ObjectType = 3
	ObjectType_TAG     ObjectType = 4
)

var ObjectType_name = map[int32]string{
	0: "UNKNOWN",
	1: "COMMIT",
	2: "BLOB",
	3: "TREE",
	4: "TAG",
}

var ObjectType_value = map[string]int32{
	"UNKNOWN": 0,
	"COMMIT":  1,
	"BLOB":    2,
	"TREE":    3,
	"TAG":     4,
}

func (x ObjectType) String() string {
	return proto.EnumName(ObjectType_name, int32(x))
}

func (ObjectType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{0}
}

type SignatureType int32

const (
	SignatureType_NONE SignatureType = 0
	SignatureType_PGP  SignatureType = 1
	SignatureType_X509 SignatureType = 2
)

var SignatureType_name = map[int32]string{
	0: "NONE",
	1: "PGP",
	2: "X509",
}

var SignatureType_value = map[string]int32{
	"NONE": 0,
	"PGP":  1,
	"X509": 2,
}

func (x SignatureType) String() string {
	return proto.EnumName(SignatureType_name, int32(x))
}

func (SignatureType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{1}
}

type OperationMsg_Operation int32

const (
	OperationMsg_UNKNOWN  OperationMsg_Operation = 0
	OperationMsg_MUTATOR  OperationMsg_Operation = 1
	OperationMsg_ACCESSOR OperationMsg_Operation = 2
)

var OperationMsg_Operation_name = map[int32]string{
	0: "UNKNOWN",
	1: "MUTATOR",
	2: "ACCESSOR",
}

var OperationMsg_Operation_value = map[string]int32{
	"UNKNOWN":  0,
	"MUTATOR":  1,
	"ACCESSOR": 2,
}

func (x OperationMsg_Operation) String() string {
	return proto.EnumName(OperationMsg_Operation_name, int32(x))
}

func (OperationMsg_Operation) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{0, 0}
}

type OperationMsg_Scope int32

const (
	OperationMsg_REPOSITORY OperationMsg_Scope = 0
	OperationMsg_SERVER     OperationMsg_Scope = 1
	OperationMsg_STORAGE    OperationMsg_Scope = 2
)

var OperationMsg_Scope_name = map[int32]string{
	0: "REPOSITORY",
	1: "SERVER",
	2: "STORAGE",
}

var OperationMsg_Scope_value = map[string]int32{
	"REPOSITORY": 0,
	"SERVER":     1,
	"STORAGE":    2,
}

func (x OperationMsg_Scope) String() string {
	return proto.EnumName(OperationMsg_Scope_name, int32(x))
}

func (OperationMsg_Scope) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{0, 1}
}

type OperationMsg struct {
	Op OperationMsg_Operation `protobuf:"varint,1,opt,name=op,proto3,enum=gitaly.OperationMsg_Operation" json:"op,omitempty"`
	// Scope level indicates what level an RPC interacts with a server:
	//   - REPOSITORY: scoped to only a single repo
	//   - SERVER: affects the entire server and potentially all repos
	//   - STORAGE: scoped to a specific storage location and all repos within
	ScopeLevel           OperationMsg_Scope `protobuf:"varint,2,opt,name=scope_level,json=scopeLevel,proto3,enum=gitaly.OperationMsg_Scope" json:"scope_level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *OperationMsg) Reset()         { *m = OperationMsg{} }
func (m *OperationMsg) String() string { return proto.CompactTextString(m) }
func (*OperationMsg) ProtoMessage()    {}
func (*OperationMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{0}
}

func (m *OperationMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OperationMsg.Unmarshal(m, b)
}
func (m *OperationMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OperationMsg.Marshal(b, m, deterministic)
}
func (m *OperationMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OperationMsg.Merge(m, src)
}
func (m *OperationMsg) XXX_Size() int {
	return xxx_messageInfo_OperationMsg.Size(m)
}
func (m *OperationMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_OperationMsg.DiscardUnknown(m)
}

var xxx_messageInfo_OperationMsg proto.InternalMessageInfo

func (m *OperationMsg) GetOp() OperationMsg_Operation {
	if m != nil {
		return m.Op
	}
	return OperationMsg_UNKNOWN
}

func (m *OperationMsg) GetScopeLevel() OperationMsg_Scope {
	if m != nil {
		return m.ScopeLevel
	}
	return OperationMsg_REPOSITORY
}

type Repository struct {
	StorageName  string `protobuf:"bytes,2,opt,name=storage_name,json=storageName,proto3" json:"storage_name,omitempty"`
	RelativePath string `protobuf:"bytes,3,opt,name=relative_path,json=relativePath,proto3" json:"relative_path,omitempty"`
	// Sets the GIT_OBJECT_DIRECTORY envvar on git commands to the value of this field.
	// It influences the object storage directory the SHA1 directories are created underneath.
	GitObjectDirectory string `protobuf:"bytes,4,opt,name=git_object_directory,json=gitObjectDirectory,proto3" json:"git_object_directory,omitempty"`
	// Sets the GIT_ALTERNATE_OBJECT_DIRECTORIES envvar on git commands to the values of this field.
	// It influences the list of Git object directories which can be used to search for Git objects.
	GitAlternateObjectDirectories []string `protobuf:"bytes,5,rep,name=git_alternate_object_directories,json=gitAlternateObjectDirectories,proto3" json:"git_alternate_object_directories,omitempty"`
	// Used in callbacks to GitLab so that it knows what repository the event is
	// associated with. May be left empty on RPC's that do not perform callbacks.
	// During project creation, `gl_repository` may not be known.
	GlRepository string `protobuf:"bytes,6,opt,name=gl_repository,json=glRepository,proto3" json:"gl_repository,omitempty"`
	// The human-readable GitLab project path (e.g. gitlab-org/gitlab-ce).
	// When hashed storage is use, this associates a project path with its
	// path on disk. The name can change over time (e.g. when a project is
	// renamed). This is primarily used for logging/debugging at the
	// moment.
	GlProjectPath        string   `protobuf:"bytes,8,opt,name=gl_project_path,json=glProjectPath,proto3" json:"gl_project_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Repository) Reset()         { *m = Repository{} }
func (m *Repository) String() string { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()    {}
func (*Repository) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{1}
}

func (m *Repository) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Repository.Unmarshal(m, b)
}
func (m *Repository) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Repository.Marshal(b, m, deterministic)
}
func (m *Repository) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Repository.Merge(m, src)
}
func (m *Repository) XXX_Size() int {
	return xxx_messageInfo_Repository.Size(m)
}
func (m *Repository) XXX_DiscardUnknown() {
	xxx_messageInfo_Repository.DiscardUnknown(m)
}

var xxx_messageInfo_Repository proto.InternalMessageInfo

func (m *Repository) GetStorageName() string {
	if m != nil {
		return m.StorageName
	}
	return ""
}

func (m *Repository) GetRelativePath() string {
	if m != nil {
		return m.RelativePath
	}
	return ""
}

func (m *Repository) GetGitObjectDirectory() string {
	if m != nil {
		return m.GitObjectDirectory
	}
	return ""
}

func (m *Repository) GetGitAlternateObjectDirectories() []string {
	if m != nil {
		return m.GitAlternateObjectDirectories
	}
	return nil
}

func (m *Repository) GetGlRepository() string {
	if m != nil {
		return m.GlRepository
	}
	return ""
}

func (m *Repository) GetGlProjectPath() string {
	if m != nil {
		return m.GlProjectPath
	}
	return ""
}

// Corresponds to Gitlab::Git::Commit
type GitCommit struct {
	Id        string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Subject   []byte        `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Body      []byte        `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	Author    *CommitAuthor `protobuf:"bytes,4,opt,name=author,proto3" json:"author,omitempty"`
	Committer *CommitAuthor `protobuf:"bytes,5,opt,name=committer,proto3" json:"committer,omitempty"`
	ParentIds []string      `protobuf:"bytes,6,rep,name=parent_ids,json=parentIds,proto3" json:"parent_ids,omitempty"`
	// If body exceeds a certain threshold, it will be nullified,
	// but its size will be set in body_size so we can know if
	// a commit had a body in the first place.
	BodySize             int64         `protobuf:"varint,7,opt,name=body_size,json=bodySize,proto3" json:"body_size,omitempty"`
	SignatureType        SignatureType `protobuf:"varint,8,opt,name=signature_type,json=signatureType,proto3,enum=gitaly.SignatureType" json:"signature_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GitCommit) Reset()         { *m = GitCommit{} }
func (m *GitCommit) String() string { return proto.CompactTextString(m) }
func (*GitCommit) ProtoMessage()    {}
func (*GitCommit) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{2}
}

func (m *GitCommit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GitCommit.Unmarshal(m, b)
}
func (m *GitCommit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GitCommit.Marshal(b, m, deterministic)
}
func (m *GitCommit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GitCommit.Merge(m, src)
}
func (m *GitCommit) XXX_Size() int {
	return xxx_messageInfo_GitCommit.Size(m)
}
func (m *GitCommit) XXX_DiscardUnknown() {
	xxx_messageInfo_GitCommit.DiscardUnknown(m)
}

var xxx_messageInfo_GitCommit proto.InternalMessageInfo

func (m *GitCommit) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GitCommit) GetSubject() []byte {
	if m != nil {
		return m.Subject
	}
	return nil
}

func (m *GitCommit) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *GitCommit) GetAuthor() *CommitAuthor {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *GitCommit) GetCommitter() *CommitAuthor {
	if m != nil {
		return m.Committer
	}
	return nil
}

func (m *GitCommit) GetParentIds() []string {
	if m != nil {
		return m.ParentIds
	}
	return nil
}

func (m *GitCommit) GetBodySize() int64 {
	if m != nil {
		return m.BodySize
	}
	return 0
}

func (m *GitCommit) GetSignatureType() SignatureType {
	if m != nil {
		return m.SignatureType
	}
	return SignatureType_NONE
}

type CommitAuthor struct {
	Name                 []byte               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email                []byte               `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	Timezone             []byte               `protobuf:"bytes,4,opt,name=timezone,proto3" json:"timezone,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CommitAuthor) Reset()         { *m = CommitAuthor{} }
func (m *CommitAuthor) String() string { return proto.CompactTextString(m) }
func (*CommitAuthor) ProtoMessage()    {}
func (*CommitAuthor) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{3}
}

func (m *CommitAuthor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommitAuthor.Unmarshal(m, b)
}
func (m *CommitAuthor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommitAuthor.Marshal(b, m, deterministic)
}
func (m *CommitAuthor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommitAuthor.Merge(m, src)
}
func (m *CommitAuthor) XXX_Size() int {
	return xxx_messageInfo_CommitAuthor.Size(m)
}
func (m *CommitAuthor) XXX_DiscardUnknown() {
	xxx_messageInfo_CommitAuthor.DiscardUnknown(m)
}

var xxx_messageInfo_CommitAuthor proto.InternalMessageInfo

func (m *CommitAuthor) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *CommitAuthor) GetEmail() []byte {
	if m != nil {
		return m.Email
	}
	return nil
}

func (m *CommitAuthor) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *CommitAuthor) GetTimezone() []byte {
	if m != nil {
		return m.Timezone
	}
	return nil
}

type ExitStatus struct {
	Value                int32    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExitStatus) Reset()         { *m = ExitStatus{} }
func (m *ExitStatus) String() string { return proto.CompactTextString(m) }
func (*ExitStatus) ProtoMessage()    {}
func (*ExitStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{4}
}

func (m *ExitStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExitStatus.Unmarshal(m, b)
}
func (m *ExitStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExitStatus.Marshal(b, m, deterministic)
}
func (m *ExitStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExitStatus.Merge(m, src)
}
func (m *ExitStatus) XXX_Size() int {
	return xxx_messageInfo_ExitStatus.Size(m)
}
func (m *ExitStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ExitStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ExitStatus proto.InternalMessageInfo

func (m *ExitStatus) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

// Corresponds to Gitlab::Git::Branch
type Branch struct {
	Name                 []byte     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	TargetCommit         *GitCommit `protobuf:"bytes,2,opt,name=target_commit,json=targetCommit,proto3" json:"target_commit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Branch) Reset()         { *m = Branch{} }
func (m *Branch) String() string { return proto.CompactTextString(m) }
func (*Branch) ProtoMessage()    {}
func (*Branch) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{5}
}

func (m *Branch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Branch.Unmarshal(m, b)
}
func (m *Branch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Branch.Marshal(b, m, deterministic)
}
func (m *Branch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Branch.Merge(m, src)
}
func (m *Branch) XXX_Size() int {
	return xxx_messageInfo_Branch.Size(m)
}
func (m *Branch) XXX_DiscardUnknown() {
	xxx_messageInfo_Branch.DiscardUnknown(m)
}

var xxx_messageInfo_Branch proto.InternalMessageInfo

func (m *Branch) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Branch) GetTargetCommit() *GitCommit {
	if m != nil {
		return m.TargetCommit
	}
	return nil
}

type Tag struct {
	Name         []byte     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id           string     `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	TargetCommit *GitCommit `protobuf:"bytes,3,opt,name=target_commit,json=targetCommit,proto3" json:"target_commit,omitempty"`
	// If message exceeds a certain threshold, it will be nullified,
	// but its size will be set in message_size so we can know if
	// a tag had a message in the first place.
	Message              []byte        `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	MessageSize          int64         `protobuf:"varint,5,opt,name=message_size,json=messageSize,proto3" json:"message_size,omitempty"`
	Tagger               *CommitAuthor `protobuf:"bytes,6,opt,name=tagger,proto3" json:"tagger,omitempty"`
	SignatureType        SignatureType `protobuf:"varint,7,opt,name=signature_type,json=signatureType,proto3,enum=gitaly.SignatureType" json:"signature_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{6}
}

func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (m *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(m, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Tag) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Tag) GetTargetCommit() *GitCommit {
	if m != nil {
		return m.TargetCommit
	}
	return nil
}

func (m *Tag) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *Tag) GetMessageSize() int64 {
	if m != nil {
		return m.MessageSize
	}
	return 0
}

func (m *Tag) GetTagger() *CommitAuthor {
	if m != nil {
		return m.Tagger
	}
	return nil
}

func (m *Tag) GetSignatureType() SignatureType {
	if m != nil {
		return m.SignatureType
	}
	return SignatureType_NONE
}

type User struct {
	GlId                 string   `protobuf:"bytes,1,opt,name=gl_id,json=glId,proto3" json:"gl_id,omitempty"`
	Name                 []byte   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email                []byte   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	GlUsername           string   `protobuf:"bytes,4,opt,name=gl_username,json=glUsername,proto3" json:"gl_username,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{7}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetGlId() string {
	if m != nil {
		return m.GlId
	}
	return ""
}

func (m *User) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *User) GetEmail() []byte {
	if m != nil {
		return m.Email
	}
	return nil
}

func (m *User) GetGlUsername() string {
	if m != nil {
		return m.GlUsername
	}
	return ""
}

type ObjectPool struct {
	Repository           *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ObjectPool) Reset()         { *m = ObjectPool{} }
func (m *ObjectPool) String() string { return proto.CompactTextString(m) }
func (*ObjectPool) ProtoMessage()    {}
func (*ObjectPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_d8a4e87e678c5ced, []int{8}
}

func (m *ObjectPool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ObjectPool.Unmarshal(m, b)
}
func (m *ObjectPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ObjectPool.Marshal(b, m, deterministic)
}
func (m *ObjectPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObjectPool.Merge(m, src)
}
func (m *ObjectPool) XXX_Size() int {
	return xxx_messageInfo_ObjectPool.Size(m)
}
func (m *ObjectPool) XXX_DiscardUnknown() {
	xxx_messageInfo_ObjectPool.DiscardUnknown(m)
}

var xxx_messageInfo_ObjectPool proto.InternalMessageInfo

func (m *ObjectPool) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

var E_OpType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*OperationMsg)(nil),
	Field:         82303,
	Name:          "gitaly.op_type",
	Tag:           "bytes,82303,opt,name=op_type",
	Filename:      "shared.proto",
}

var E_Storage = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         91233,
	Name:          "gitaly.storage",
	Tag:           "varint,91233,opt,name=storage",
	Filename:      "shared.proto",
}

var E_Repository = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         91234,
	Name:          "gitaly.repository",
	Tag:           "varint,91234,opt,name=repository",
	Filename:      "shared.proto",
}

var E_TargetRepository = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         91235,
	Name:          "gitaly.target_repository",
	Tag:           "varint,91235,opt,name=target_repository",
	Filename:      "shared.proto",
}

var E_AdditionalRepository = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         91236,
	Name:          "gitaly.additional_repository",
	Tag:           "varint,91236,opt,name=additional_repository",
	Filename:      "shared.proto",
}

func init() {
	proto.RegisterEnum("gitaly.ObjectType", ObjectType_name, ObjectType_value)
	proto.RegisterEnum("gitaly.SignatureType", SignatureType_name, SignatureType_value)
	proto.RegisterEnum("gitaly.OperationMsg_Operation", OperationMsg_Operation_name, OperationMsg_Operation_value)
	proto.RegisterEnum("gitaly.OperationMsg_Scope", OperationMsg_Scope_name, OperationMsg_Scope_value)
	proto.RegisterType((*OperationMsg)(nil), "gitaly.OperationMsg")
	proto.RegisterType((*Repository)(nil), "gitaly.Repository")
	proto.RegisterType((*GitCommit)(nil), "gitaly.GitCommit")
	proto.RegisterType((*CommitAuthor)(nil), "gitaly.CommitAuthor")
	proto.RegisterType((*ExitStatus)(nil), "gitaly.ExitStatus")
	proto.RegisterType((*Branch)(nil), "gitaly.Branch")
	proto.RegisterType((*Tag)(nil), "gitaly.Tag")
	proto.RegisterType((*User)(nil), "gitaly.User")
	proto.RegisterType((*ObjectPool)(nil), "gitaly.ObjectPool")
	proto.RegisterExtension(E_OpType)
	proto.RegisterExtension(E_Storage)
	proto.RegisterExtension(E_Repository)
	proto.RegisterExtension(E_TargetRepository)
	proto.RegisterExtension(E_AdditionalRepository)
}

func init() { proto.RegisterFile("shared.proto", fileDescriptor_d8a4e87e678c5ced) }

var fileDescriptor_d8a4e87e678c5ced = []byte{
	// 1037 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x6e, 0x1c, 0xe7, 0xef, 0x24, 0x2d, 0xde, 0xa1, 0x2b, 0x59, 0x45, 0xdd, 0x2d, 0x46, 0x42,
	0xd5, 0xaa, 0xa4, 0x55, 0x57, 0x20, 0xb6, 0x20, 0xa1, 0xb4, 0x64, 0xab, 0x2e, 0x6d, 0x1c, 0x4d,
	0x5c, 0xfe, 0x6e, 0xac, 0x49, 0x3c, 0x3b, 0x19, 0xe4, 0x64, 0x2c, 0x7b, 0x52, 0xd1, 0x5e, 0x22,
	0xae, 0xb8, 0xe2, 0x25, 0xb8, 0xe5, 0x31, 0x10, 0xaf, 0xc1, 0xc2, 0x7b, 0x80, 0x66, 0xc6, 0x4e,
	0xdd, 0x1f, 0x60, 0xf7, 0x6e, 0xce, 0xf1, 0x77, 0xbe, 0x39, 0xf3, 0xcd, 0x77, 0x26, 0x81, 0x4e,
	0x36, 0x25, 0x29, 0x8d, 0xba, 0x49, 0x2a, 0xa4, 0x40, 0x75, 0xc6, 0x25, 0x89, 0x2f, 0x37, 0x1e,
	0x33, 0x21, 0x58, 0x4c, 0x77, 0x75, 0x76, 0xbc, 0x78, 0xb9, 0x2b, 0xf9, 0x8c, 0x66, 0x92, 0xcc,
	0x12, 0x03, 0xdc, 0xd8, 0xba, 0x0d, 0x88, 0x68, 0x36, 0x49, 0x79, 0x22, 0x45, 0x6a, 0x10, 0xde,
	0xab, 0x0a, 0x74, 0xfc, 0x84, 0xa6, 0x44, 0x72, 0x31, 0x3f, 0xcb, 0x18, 0xea, 0x82, 0x25, 0x12,
	0xb7, 0xb2, 0x55, 0xd9, 0x5e, 0xdb, 0x7f, 0xd4, 0x35, 0x1b, 0x75, 0xcb, 0x88, 0xeb, 0x00, 0x5b,
	0x22, 0x41, 0x9f, 0x40, 0x3b, 0x9b, 0x88, 0x84, 0x86, 0x31, 0xbd, 0xa0, 0xb1, 0x6b, 0xe9, 0xc2,
	0x8d, 0x7b, 0x0b, 0x47, 0x0a, 0x87, 0x41, 0xc3, 0x4f, 0x15, 0xda, 0x7b, 0x0a, 0xad, 0x25, 0x02,
	0xb5, 0xa1, 0x71, 0x3e, 0xf8, 0x62, 0xe0, 0x7f, 0x35, 0x70, 0x56, 0x54, 0x70, 0x76, 0x1e, 0xf4,
	0x02, 0x1f, 0x3b, 0x15, 0xd4, 0x81, 0x66, 0xef, 0xe8, 0xa8, 0x3f, 0x1a, 0xf9, 0xd8, 0xb1, 0xbc,
	0x3d, 0xa8, 0x69, 0x26, 0xb4, 0x06, 0x80, 0xfb, 0x43, 0x7f, 0x74, 0x12, 0xf8, 0xf8, 0x1b, 0x67,
	0x05, 0x01, 0xd4, 0x47, 0x7d, 0xfc, 0x65, 0x5f, 0x95, 0xb4, 0xa1, 0x31, 0x0a, 0x7c, 0xdc, 0x3b,
	0xee, 0x3b, 0x96, 0xf7, 0xab, 0x05, 0x80, 0x69, 0x22, 0x32, 0x2e, 0x45, 0x7a, 0x89, 0xde, 0x85,
	0x4e, 0x26, 0x45, 0x4a, 0x18, 0x0d, 0xe7, 0x64, 0x46, 0x75, 0xcf, 0x2d, 0xdc, 0xce, 0x73, 0x03,
	0x32, 0xa3, 0xe8, 0x3d, 0x58, 0x4d, 0x69, 0x4c, 0x24, 0xbf, 0xa0, 0x61, 0x42, 0xe4, 0xd4, 0xad,
	0x6a, 0x4c, 0xa7, 0x48, 0x0e, 0x89, 0x9c, 0xa2, 0x3d, 0x58, 0x67, 0x5c, 0x86, 0x62, 0xfc, 0x1d,
	0x9d, 0xc8, 0x30, 0xe2, 0x29, 0x9d, 0x28, 0x7e, 0xd7, 0xd6, 0x58, 0xc4, 0xb8, 0xf4, 0xf5, 0xa7,
	0xcf, 0x8b, 0x2f, 0xe8, 0x18, 0xb6, 0x54, 0x05, 0x89, 0x25, 0x4d, 0xe7, 0x44, 0xd2, 0xdb, 0xb5,
	0x9c, 0x66, 0x6e, 0x6d, 0xab, 0xba, 0xdd, 0xc2, 0x9b, 0x8c, 0xcb, 0x5e, 0x01, 0xbb, 0x49, 0xc3,
	0x69, 0xa6, 0xfa, 0x63, 0x71, 0x98, 0x2e, 0xcf, 0xe4, 0xd6, 0x4d, 0x7f, 0x2c, 0x2e, 0x9d, 0xf3,
	0x7d, 0x78, 0x8b, 0xc5, 0x61, 0x92, 0x0a, 0xbd, 0x87, 0x3e, 0x46, 0x53, 0xc3, 0x56, 0x59, 0x3c,
	0x34, 0x59, 0x75, 0x8e, 0x17, 0x76, 0xb3, 0xe2, 0x58, 0x2f, 0xec, 0x66, 0xc3, 0x69, 0x62, 0x5b,
	0xc1, 0xbc, 0x5f, 0x2c, 0x68, 0x1d, 0x73, 0x79, 0x24, 0x66, 0x33, 0x2e, 0xd1, 0x1a, 0x58, 0x3c,
	0xd2, 0x96, 0x68, 0x61, 0x8b, 0x47, 0xc8, 0x85, 0x46, 0xb6, 0xd0, 0x2d, 0x69, 0xe9, 0x3a, 0xb8,
	0x08, 0x11, 0x02, 0x7b, 0x2c, 0xa2, 0x4b, 0xad, 0x56, 0x07, 0xeb, 0x35, 0xda, 0x81, 0x3a, 0x59,
	0xc8, 0xa9, 0x48, 0xb5, 0x2e, 0xed, 0xfd, 0xf5, 0xc2, 0x1b, 0x86, 0xbd, 0xa7, 0xbf, 0xe1, 0x1c,
	0x83, 0xf6, 0xa1, 0x35, 0xd1, 0x79, 0x49, 0x53, 0xb7, 0xf6, 0x1f, 0x05, 0xd7, 0x30, 0xb4, 0x09,
	0x90, 0x90, 0x94, 0xce, 0x65, 0xc8, 0xa3, 0xcc, 0xad, 0x6b, 0xfd, 0x5a, 0x26, 0x73, 0x12, 0x65,
	0xe8, 0x1d, 0x68, 0xa9, 0x46, 0xc2, 0x8c, 0x5f, 0x51, 0xb7, 0xb1, 0x55, 0xd9, 0xae, 0xe2, 0xa6,
	0x4a, 0x8c, 0xf8, 0x15, 0x45, 0x9f, 0xc2, 0x5a, 0xc6, 0xd9, 0x9c, 0xc8, 0x45, 0x4a, 0x43, 0x79,
	0x99, 0x50, 0x2d, 0xd1, 0xda, 0xfe, 0xc3, 0x62, 0xd3, 0x51, 0xf1, 0x35, 0xb8, 0x4c, 0x28, 0x5e,
	0xcd, 0xca, 0xa1, 0xf7, 0x63, 0x05, 0x3a, 0xe5, 0xae, 0x94, 0x00, 0xda, 0x52, 0x15, 0x23, 0x80,
	0x5a, 0xa3, 0x75, 0xa8, 0xd1, 0x19, 0xe1, 0x71, 0x2e, 0x96, 0x09, 0x50, 0x17, 0xec, 0x88, 0x48,
	0xaa, 0xa5, 0x6a, 0xab, 0x81, 0xd1, 0x93, 0xda, 0x2d, 0x26, 0xb5, 0x1b, 0x14, 0xa3, 0x8c, 0x35,
	0x0e, 0x6d, 0x40, 0x53, 0x4d, 0xf7, 0x95, 0x98, 0x53, 0x2d, 0x64, 0x07, 0x2f, 0x63, 0xcf, 0x03,
	0xe8, 0x7f, 0xcf, 0xe5, 0x48, 0x12, 0xb9, 0xc8, 0xd4, 0x7e, 0x17, 0x24, 0x5e, 0x98, 0x26, 0x6a,
	0xd8, 0x04, 0x5e, 0x00, 0xf5, 0xc3, 0x94, 0xcc, 0x27, 0xd3, 0x7b, 0x7b, 0xfc, 0x08, 0x56, 0x25,
	0x49, 0x19, 0x95, 0xa1, 0x91, 0x55, 0xf7, 0xda, 0xde, 0x7f, 0x50, 0xa8, 0xb0, 0x34, 0x03, 0xee,
	0x18, 0x9c, 0x89, 0xbc, 0x9f, 0x2c, 0xa8, 0x06, 0x84, 0xdd, 0xcb, 0x69, 0x6c, 0x63, 0x2d, 0x6d,
	0x73, 0x67, 0x8f, 0xea, 0x6b, 0xed, 0xa1, 0xec, 0x36, 0xa3, 0x59, 0x46, 0x58, 0x71, 0xf0, 0x22,
	0x54, 0x83, 0x9c, 0x2f, 0xcd, 0xe5, 0xd6, 0xf4, 0xe5, 0xb6, 0xf3, 0x9c, 0xbe, 0xdf, 0x1d, 0xa8,
	0x4b, 0xc2, 0x18, 0x4d, 0xf5, 0x84, 0xfc, 0xab, 0xfb, 0x0c, 0xe6, 0x1e, 0x37, 0x34, 0xde, 0xc0,
	0x0d, 0x2f, 0xc1, 0x3e, 0xcf, 0x68, 0x8a, 0xde, 0x86, 0x1a, 0x8b, 0xc3, 0xe5, 0xc8, 0xd8, 0x2c,
	0x3e, 0x89, 0x96, 0x0a, 0x59, 0xf7, 0x39, 0xa3, 0x5a, 0x76, 0xc6, 0x63, 0x68, 0xb3, 0x38, 0x5c,
	0x64, 0x6a, 0xf6, 0x67, 0x34, 0x7f, 0x4d, 0x80, 0xc5, 0xe7, 0x79, 0xc6, 0x7b, 0x0e, 0x60, 0x5e,
	0x84, 0xa1, 0x10, 0x31, 0xfa, 0x18, 0xa0, 0xf4, 0x0e, 0x54, 0xf4, 0x29, 0x51, 0xd1, 0xef, 0xf5,
	0x6b, 0x70, 0x68, 0xff, 0xfc, 0xdb, 0x4e, 0x05, 0x97, 0xb0, 0x4f, 0x0e, 0x0b, 0x1e, 0xd5, 0xfd,
	0xcd, 0xe7, 0x17, 0xa0, 0x7e, 0xe4, 0x9f, 0x9d, 0x9d, 0x04, 0x4e, 0x05, 0x35, 0xc1, 0x3e, 0x3c,
	0xf5, 0x0f, 0x1d, 0x4b, 0xad, 0x02, 0xdc, 0xef, 0x3b, 0x55, 0xd4, 0x80, 0x6a, 0xd0, 0x3b, 0x76,
	0xec, 0x27, 0x3b, 0xb0, 0x7a, 0x43, 0x13, 0x85, 0x19, 0xf8, 0x83, 0xbe, 0xb3, 0xa2, 0x30, 0xc3,
	0xe3, 0xa1, 0x21, 0xf8, 0xfa, 0xc3, 0xbd, 0x67, 0x8e, 0x75, 0xe0, 0x43, 0x43, 0x24, 0x5a, 0x58,
	0xf4, 0xe8, 0x8e, 0xe3, 0xcf, 0xa8, 0x9c, 0x8a, 0xc8, 0x4f, 0xd4, 0x8f, 0x41, 0xe6, 0xfe, 0xfd,
	0xc3, 0xad, 0xe9, 0x2f, 0xff, 0x94, 0xe0, 0xba, 0x48, 0xd4, 0x6e, 0x07, 0xcf, 0xa0, 0x91, 0x3f,
	0xdb, 0x68, 0xf3, 0x0e, 0xe1, 0x73, 0x4e, 0xe3, 0x25, 0xdf, 0x1f, 0xbf, 0x2b, 0xbe, 0x26, 0x2e,
	0xf0, 0x07, 0x9f, 0x95, 0x75, 0xfb, 0xbf, 0xea, 0x57, 0x79, 0x75, 0xa9, 0xe4, 0xe0, 0x14, 0x1e,
	0xe4, 0x7e, 0x7e, 0x7d, 0x9e, 0x3f, 0x73, 0x1e, 0xc7, 0x54, 0x5e, 0x5f, 0xcf, 0x41, 0x00, 0x0f,
	0x49, 0x14, 0x71, 0x05, 0x23, 0xf1, 0x1b, 0x30, 0xfe, 0x95, 0x33, 0xae, 0x5f, 0x57, 0x97, 0x2e,
	0x7d, 0xef, 0x5b, 0x25, 0x5f, 0x4c, 0xc6, 0xdd, 0x89, 0x98, 0xed, 0x9a, 0xe5, 0x07, 0x22, 0x65,
	0xbb, 0x46, 0x54, 0xf3, 0xc7, 0x60, 0x97, 0x89, 0x3c, 0x4e, 0xc6, 0xe3, 0xba, 0x4e, 0x3d, 0xfd,
	0x27, 0x00, 0x00, 0xff, 0xff, 0xf3, 0x84, 0x06, 0x14, 0x72, 0x08, 0x00, 0x00,
}
