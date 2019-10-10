package swlib

const FamilyName = "switch"

type CommandType int

const (
	CmdUnspec CommandType = iota
	CmdGetSwitch
	CmdNewAttr
	CmdListGlobal
	CmdGetGlobal
	CmdSetGlobal
	CmdListPort
	CmdGetPort
	CmdSetPort
	CmdListVLan
	CmdGetVLan
	CmdSetVLan
)

type AttributeType int

const (
	AttrUnspec AttributeType = iota
	AttrType
	AttrID
	AttrDevName
	AttrAlias
	AttrName
	AttrVLANs
	AttrPorts
	AttrPortMap
	AttrCPUPort
	AttrOPID
	AttrOPType
	AttrOPName
	AttrOPPort
	AttrOPVLAN
	AttrOPValueInt
	AttrOPValueStr
	AttrOPValuePorts
	AttrOPValueLink
	AttrOPDescription
	AttrPort
)
