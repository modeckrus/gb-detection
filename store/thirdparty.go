package store

type ThirdpartyDb interface {
	AddUrl(timeRecord.Url)
	DeleteUrl(timeRecord.Url)
	All()
}
