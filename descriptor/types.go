package descriptor

// UPDATE: For the first repo only release not needed, therefore not active

// import (
// 	"github.com/aas-core-works/aas-core3.1-golang/jsonization"
// 	"github.com/aas-core-works/aas-core3.1-golang/reporting"
// 	"github.com/aas-core-works/aas-core3.1-golang/types"
// )

// /**
// IMPORTANT NOTE:
// These are only aimed at reading shell descriptor and submodel descriptor endpoints.
// They will not work with the "ToJsonable" method from aas-core-works!
// Maybe that will get added here in the future.
// */

// // these are used to check wether the interfaces are actually implemented by their structs
// // if there is something wrong here, check the structs!
// var (
// 	_ IDescriptor                           = (*Descriptor)(nil)
// 	_ IAssetAdministrationShellDescriptor   = (*AssetAdministrationShellDescriptor)(nil)
// 	_ IEndpoint                             = (*Endpoint)(nil)
// 	_ IProtocolInformation                  = (*ProtocolInformation)(nil)
// 	_ IProtocolInformationSecurityAttribute = (*ProtocolInformationSecurityAttribute)(nil)
// 	_ ISubmodelDescriptor                   = (*SubmodelDescriptor)(nil)
// )

// const (
// 	ModelTypeAssetAdministrationShellDescriptor types.ModelType = iota + 100
// 	ModelTypeSubmodelDescriptor
// 	ModelTypeDescriptor // just a filler, should never be used
// 	ModelTypeEndpoint
// 	ModelTypeProtocolInformation
// 	ModelTypeProtocolInformationSecurityAttribute
// )

// func newDeserializationError(message string) *jsonization.DeserializationError {
// 	return &jsonization.DeserializationError{
// 		Path:    &reporting.Path{},
// 		Message: message,
// 	}
// }

// // fake the descend methods
// type FakeDescendable struct{}

// func (f *FakeDescendable) Descend(action func(types.IClass) bool) (abort bool)     { return false }
// func (f *FakeDescendable) DescendOnce(action func(types.IClass) bool) (abort bool) { return false }

// // IDescriptor represents properties shared by AAS and submodel descriptors.
// type IDescriptor interface {
// 	types.IClass

// 	Description() []types.ILangStringTextType
// 	SetDescription([]types.ILangStringTextType)

// 	DisplayName() []types.ILangStringNameType
// 	SetDisplayName([]types.ILangStringNameType)

// 	Extensions() []types.IExtension
// 	SetExtensions([]types.IExtension)

// 	Administration() types.IAdministrativeInformation
// 	SetAdministration(types.IAdministrativeInformation)

// 	IDShort() *string
// 	SetIDShort(*string)
// }

// type Descriptor struct {
// 	FakeDescendable

// 	description    []types.ILangStringTextType
// 	displayName    []types.ILangStringNameType
// 	extensions     []types.IExtension
// 	administration types.IAdministrativeInformation
// 	idShort        *string
// }

// func (x *Descriptor) ModelType() types.ModelType { return ModelTypeDescriptor }

// func (x *Descriptor) Description() []types.ILangStringTextType     { return x.description }
// func (x *Descriptor) SetDescription(v []types.ILangStringTextType) { x.description = v }

// func (x *Descriptor) DisplayName() []types.ILangStringNameType     { return x.displayName }
// func (x *Descriptor) SetDisplayName(v []types.ILangStringNameType) { x.displayName = v }

// func (x *Descriptor) Extensions() []types.IExtension     { return x.extensions }
// func (x *Descriptor) SetExtensions(v []types.IExtension) { x.extensions = v }

// func (x *Descriptor) Administration() types.IAdministrativeInformation     { return x.administration }
// func (x *Descriptor) SetAdministration(v types.IAdministrativeInformation) { x.administration = v }

// func (x *Descriptor) IDShort() *string     { return x.idShort }
// func (x *Descriptor) SetIDShort(v *string) { x.idShort = v }

// type IAssetAdministrationShellDescriptor interface {
// 	IDescriptor

// 	ID() string
// 	SetID(string)

// 	AssetKind() *types.AssetKind
// 	SetAssetKind(*types.AssetKind)

// 	AssetType() *string
// 	SetAssetType(*string)

// 	GlobalAssetID() *string
// 	SetGlobalAssetID(*string)

// 	SpecificAssetIDs() []types.ISpecificAssetID
// 	SetSpecificAssetIDs([]types.ISpecificAssetID)

// 	Endpoints() []IEndpoint
// 	SetEndpoints([]IEndpoint)

// 	SubmodelDescriptors() []ISubmodelDescriptor
// 	SetSubmodelDescriptors([]ISubmodelDescriptor)
// }
// type AssetAdministrationShellDescriptor struct {
// 	Descriptor
// 	FakeDescendable

// 	id                  string
// 	assetKind           *types.AssetKind
// 	assetType           *string
// 	endpoints           []IEndpoint
// 	globalAssetID       *string
// 	specificAssetIDs    []types.ISpecificAssetID
// 	submodelDescriptors []ISubmodelDescriptor
// }

// func NewAssetAdministrationShellDescriptor(id string) *AssetAdministrationShellDescriptor {
// 	return &AssetAdministrationShellDescriptor{id: id}
// }

// func (x *AssetAdministrationShellDescriptor) ModelType() types.ModelType {
// 	return ModelTypeAssetAdministrationShellDescriptor
// }

// func (x *AssetAdministrationShellDescriptor) ID() string {
// 	return x.id
// }
// func (x *AssetAdministrationShellDescriptor) SetID(v string) {
// 	x.id = v
// }

// func (x *AssetAdministrationShellDescriptor) AssetKind() *types.AssetKind {
// 	return x.assetKind
// }
// func (x *AssetAdministrationShellDescriptor) SetAssetKind(v *types.AssetKind) {
// 	x.assetKind = v
// }

// func (x *AssetAdministrationShellDescriptor) AssetType() *string {
// 	return x.assetType
// }
// func (x *AssetAdministrationShellDescriptor) SetAssetType(v *string) {
// 	x.assetType = v
// }

// func (x *AssetAdministrationShellDescriptor) GlobalAssetID() *string {
// 	return x.globalAssetID
// }
// func (x *AssetAdministrationShellDescriptor) SetGlobalAssetID(v *string) {
// 	x.globalAssetID = v
// }

// func (x *AssetAdministrationShellDescriptor) SpecificAssetIDs() []types.ISpecificAssetID {
// 	return x.specificAssetIDs
// }
// func (x *AssetAdministrationShellDescriptor) SetSpecificAssetIDs(v []types.ISpecificAssetID) {
// 	x.specificAssetIDs = v
// }

// func (x *AssetAdministrationShellDescriptor) Endpoints() []IEndpoint {
// 	return x.endpoints
// }
// func (x *AssetAdministrationShellDescriptor) SetEndpoints(v []IEndpoint) {
// 	x.endpoints = v
// }

// func (x *AssetAdministrationShellDescriptor) SubmodelDescriptors() []ISubmodelDescriptor {
// 	return x.submodelDescriptors
// }
// func (x *AssetAdministrationShellDescriptor) SetSubmodelDescriptors(v []ISubmodelDescriptor) {
// 	x.submodelDescriptors = v
// }

// type IEndpoint interface {
// 	types.IClass

// 	Interface() string
// 	SetInterface(string)

// 	ProtocolInformation() IProtocolInformation
// 	SetProtocolInformation(IProtocolInformation)
// }

// type Endpoint struct {
// 	FakeDescendable

// 	theInterface        string
// 	protocolInformation IProtocolInformation
// }

// func NewEndpoint(theInterface string, protocolInformation IProtocolInformation) *Endpoint {
// 	return &Endpoint{theInterface: theInterface, protocolInformation: protocolInformation}
// }

// func (x *Endpoint) ModelType() types.ModelType { return ModelTypeEndpoint }

// func (x *Endpoint) Interface() string {
// 	return x.theInterface
// }
// func (x *Endpoint) SetInterface(v string) {
// 	x.theInterface = v
// }

// func (x *Endpoint) ProtocolInformation() IProtocolInformation {
// 	return x.protocolInformation
// }
// func (x *Endpoint) SetProtocolInformation(v IProtocolInformation) {
// 	x.protocolInformation = v
// }

// type IProtocolInformation interface {
// 	types.IClass

// 	Href() string
// 	SetHref(string)

// 	EndpointProtocol() *string
// 	SetEndpointProtocol(*string)

// 	EndpointProtocolVersion() []string
// 	SetEndpointProtocolVersion([]string)

// 	Subprotocol() *string
// 	SetSubprotocol(*string)

// 	SubprotocolBody() *string
// 	SetSubprotocolBody(*string)

// 	SubprotocolBodyEncoding() *string
// 	SetSubprotocolBodyEncoding(*string)

// 	SecurityAttributes() []IProtocolInformationSecurityAttribute
// 	SetSecurityAttributes([]IProtocolInformationSecurityAttribute)
// }

// type ProtocolInformation struct {
// 	FakeDescendable

// 	href                    string
// 	endpointProtocol        *string
// 	endpointProtocolVersion []string
// 	subprotocol             *string
// 	subprotocolBody         *string
// 	subprotocolBodyEncoding *string
// 	securityAttributes      []IProtocolInformationSecurityAttribute
// }

// func NewProtocolInformation(href string) *ProtocolInformation {
// 	return &ProtocolInformation{
// 		href: href,
// 	}
// }

// func (x *ProtocolInformation) ModelType() types.ModelType { return ModelTypeProtocolInformation }

// func (p *ProtocolInformation) Href() string {
// 	return p.href
// }
// func (p *ProtocolInformation) SetHref(v string) {
// 	p.href = v
// }

// func (p *ProtocolInformation) EndpointProtocol() *string {
// 	return p.endpointProtocol
// }
// func (p *ProtocolInformation) SetEndpointProtocol(v *string) {
// 	p.endpointProtocol = v
// }

// func (p *ProtocolInformation) EndpointProtocolVersion() []string {
// 	return p.endpointProtocolVersion
// }
// func (p *ProtocolInformation) SetEndpointProtocolVersion(v []string) {
// 	p.endpointProtocolVersion = v
// }

// func (p *ProtocolInformation) Subprotocol() *string {
// 	return p.subprotocol
// }
// func (p *ProtocolInformation) SetSubprotocol(v *string) {
// 	p.subprotocol = v
// }

// func (p *ProtocolInformation) SubprotocolBody() *string {
// 	return p.subprotocolBody
// }
// func (p *ProtocolInformation) SetSubprotocolBody(v *string) {
// 	p.subprotocolBody = v
// }

// func (p *ProtocolInformation) SubprotocolBodyEncoding() *string {
// 	return p.subprotocolBodyEncoding
// }
// func (p *ProtocolInformation) SetSubprotocolBodyEncoding(v *string) {
// 	p.subprotocolBodyEncoding = v
// }

// func (p *ProtocolInformation) SecurityAttributes() []IProtocolInformationSecurityAttribute {
// 	return p.securityAttributes
// }
// func (p *ProtocolInformation) SetSecurityAttributes(
// 	v []IProtocolInformationSecurityAttribute,
// ) {
// 	p.securityAttributes = v
// }

// type IProtocolInformationSecurityAttribute interface {
// 	types.IClass

// 	Type() string
// 	SetType(string)

// 	Key() string
// 	SetKey(string)

// 	Value() string
// 	SetValue(string)
// }

// type ProtocolInformationSecurityAttribute struct {
// 	FakeDescendable

// 	theType string
// 	key     string
// 	value   string
// }

// func NewProtocolInformationSecurityAttribute(theType string, theKey string, value string) *ProtocolInformationSecurityAttribute {
// 	return &ProtocolInformationSecurityAttribute{
// 		theType: theType,
// 		key:     theKey,
// 		value:   value,
// 	}
// }

// func (x *ProtocolInformationSecurityAttribute) ModelType() types.ModelType {
// 	return ModelTypeProtocolInformationSecurityAttribute
// }

// func (p *ProtocolInformationSecurityAttribute) Type() string {
// 	return p.theType
// }
// func (p *ProtocolInformationSecurityAttribute) SetType(v string) {
// 	p.theType = v
// }

// func (p *ProtocolInformationSecurityAttribute) Key() string {
// 	return p.key
// }
// func (p *ProtocolInformationSecurityAttribute) SetKey(v string) {
// 	p.key = v
// }

// func (p *ProtocolInformationSecurityAttribute) Value() string {
// 	return p.value
// }
// func (p *ProtocolInformationSecurityAttribute) SetValue(v string) {
// 	p.value = v
// }

// type ISubmodelDescriptor interface {
// 	IDescriptor

// 	ID() string
// 	SetID(string)

// 	SemanticID() types.IReference
// 	SetSemanticID(types.IReference)

// 	SupplementalSemanticIDs() []types.IReference
// 	SetSupplementalSemanticIDs([]types.IReference)

// 	Endpoints() []IEndpoint
// 	SetEndpoints([]IEndpoint)
// }

// type SubmodelDescriptor struct {
// 	Descriptor
// 	id                      string
// 	semanticID              types.IReference
// 	supplementalSemanticIDs []types.IReference
// 	endpoints               []IEndpoint
// }

// func NewSubmodelDescriptor(id string) *SubmodelDescriptor {
// 	return &SubmodelDescriptor{id: id}
// }

// func (x *SubmodelDescriptor) ModelType() types.ModelType { return ModelTypeSubmodelDescriptor }

// func (x *SubmodelDescriptor) ID() string {
// 	return x.id
// }
// func (x *SubmodelDescriptor) SetID(v string) {
// 	x.id = v
// }

// func (x *SubmodelDescriptor) SemanticID() types.IReference {
// 	return x.semanticID
// }
// func (x *SubmodelDescriptor) SetSemanticID(v types.IReference) {
// 	x.semanticID = v
// }

// func (x *SubmodelDescriptor) SupplementalSemanticIDs() []types.IReference {
// 	return x.supplementalSemanticIDs
// }
// func (x *SubmodelDescriptor) SetSupplementalSemanticIDs(v []types.IReference) {
// 	x.supplementalSemanticIDs = v
// }

// func (x *SubmodelDescriptor) Endpoints() []IEndpoint {
// 	return x.endpoints
// }
// func (x *SubmodelDescriptor) SetEndpoints(v []IEndpoint) {
// 	x.endpoints = v
// }
