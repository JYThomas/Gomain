package itf_DataSources

type DataSources interface{
	GetDomainNames() []string
	GetDataSourceData() struct
} 