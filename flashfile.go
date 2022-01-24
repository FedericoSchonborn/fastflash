package main

import "encoding/xml"

type Flashfile struct {
	XMLName xml.Name `xml:"flashing"`
	Header  *Header  `xml:"header"`
	Steps   []*Step  `xml:"steps>step"`
}

type Header struct {
	PhoneModel        *PhoneModel        `xml:"phone_model"`
	SoftwareVersion   *SoftwareVersion   `xml:"software_version"`
	SubsidyLockConfig *SubsidyLockConfig `xml:"subsidy_lock_config"`
	RegulatoryConfig  *RegulatoryConfig  `xml:"regulatory_config"`
	Sparsing          *Sparsing          `xml:"sparsing"`
	Interfaces        []*Interface       `xml:"interfaces>interface"`
}

type PhoneModel struct {
	Model string `xml:"model,attr"`
}

type SoftwareVersion struct {
	Version string `xml:"version,attr"`
}

type SubsidyLockConfig struct {
	MD5  string `xml:"MD5,attr"`
	Name string `xml:"name,attr"`
}

type RegulatoryConfig struct {
	SHA1 string `xml:"SHA1,attr"`
	Name string `xml:"name,attr"`
}

type Sparsing struct {
	Enabled       string `xml:"enabled,attr"`
	MaxSparseSize string `xml:"max_sparse_size,attr"`
}

type Interface struct {
	Name string `xml:"name,attr"`
}

type Step struct {
	MD5       string `xml:"MD5,attr"`
	Operation string `xml:"operation,attr"`
	Var       string `xml:"var,attr"`
	Filename  string `xml:"filename,attr"`
	Partition string `xml:"partition,attr"`
}
