package datastruct

// Error code definition
const (
	//App Error
	//Status OK value <=500
	ErrSuccess                  int = 00
	ErrFailedQuery              int = 100
	ErrInvalidParameter         int = 101
	ErrReadDBValue              int = 102
	ErrFailedPreparedStatements int = 103
	ErrNotInitializedTableName  int = 104
	ErrNotInitializedCondition  int = 105
	ErrFTPConnection            int = 300
	ErrFTPAuth                  int = 301
	ErrFTPChangeDir             int = 302
	ErrFTPCopyFile              int = 303
	ErrStringConvert            int = 305
	ErrInvalidFormat            int = 306
	ErrGetVoiceResponsiveFailed int = 301
	//Bad Request value 501-900

	//value 901-998 Unauthorized
	ErrUnauthorized      int = 901
	ErrWrongUserPassword int = 902

	//value >999 //Internal Server Error
	ErrOthers int = 999

	//MYSQL SP Return Error
	ErrNoData                int = -1000
	ErrUnknownAction         int = -1100
	ErrOnlyNoShowCanRecall   int = -1200
	ErrCannotModifyDoneQueue int = -1500
)

// Error message definition
const (
	DescSuccess                   string = "Success"
	DescInvalidParameter          string = "Invalid Parameter"
	DescCommonError               string = "Common Error"
	DescGalleryNotFound           string = "Gallery Not Found"
	DescQueueCountFailed          string = "Get Queue Count Failed"
	DescQueueCounterFailed        string = "Get Queue Counter Failed"
	DescQueueListFailed           string = "Get Queue List Failed"
	DescServiceGroupListFailed    string = "Get Service Group List Failed"
	DescQueueGetGalleryNameFailed string = "Get Gallery Name Failed"
	DescStringConvert             string = "String Convert Error"
	DescReadDatabaseError         string = "Read Database Value Error"
	DescGetMDNFromZSmartFailed    string = "Get MDN Info Failed"
	DescVoiceResponsiveFailed     string = "Get Voice Responsive Failed"
)
