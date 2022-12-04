package parse

import (
	. "github.com/ilius/ls-go/common"
)

var Titles = [][2]string{
	{C_Name, "Name"},
	{C_Mode, "Mode"},
	{C_Owner, "Owner"},
	{C_Group, "Group"},
	{C_Size, "Size"},
	{C_MTime, "Modified Time"},
	{C_CTime, "Change Time"},
	{C_ATime, "Access Time"},
	{C_Inode, "inode"},
	{C_ModeOct, "Oct"},
	{C_HardLinks, "Hard Links"},
	{C_Blocks, "Blocks"},
}
