package datastruct

// Constant Value
type AuthRequest struct {
	Uname  string `json:"username"`
	Passwd string `json:"password"`
}

type RequiredResponse struct {
	ResponseCode int `json:"respCode"`
}

type FileResponse struct {
	RequiredResponse
	ResponseDesc []string `json:"responseDesc"`
}

type NotFileResponse struct {
	RequiredResponse
	ResponseDesc string `json:"responseDesc"`
}

const (

	//DBName Added by Raditya Pratama
	DbName string = "ecluster"

	//QueuePriority
	PriorityRecall     int = 100
	PriorityNormal     int = 100
	PriorityTransfer   int = 100
	PriorityMobile     int = 100
	PriorityWalkDirect int = 100 // tidak dipanggil, namun langsung menjadi status QSLobbyOfficer
	//khsusu utk mobile yang booking jauh dr waktu yang diminta, menggunakan code M
	PriorityVIP int = 1
	//Lobby officer has right to add new que

	//Do not change this value, add new value instead
	//These values use on stored procedure on mysql server

	//Queue Status < 100 will show in queue list
	QSQueuing  int = 10
	QSTransfer int = 20
	QSReCall   int = 30

	//Queue Status 100-500 will show in serving list
	QSCalling        int = 100
	QSDirectServing  int = 105
	QSAcceptTransfer int = 110
	QSServing        int = 200
	QSTagging        int = 300

	//Queue Status 800-1000 will not show in list
	QSNoShow int = 800

	//Queue Status >900 cannot modify queue
	QSDone int = 900

	//Queue Status >= 2000 Agent Related Functions
	QSLogin  int = 2100
	QSReturn int = 2200
	QSLogout int = 3100
	QSBreak  int = 3200
	// //Serv Status
	// //Do not change this value, add new value instead
	// SVCalling    int = 100
	// SVNoShow     int = 900
	// SVTagging    int = 120
	// SVTransfered int = 20

	EclusterSecretKey      string = "eClu5t3RaP!"
	ResponsiveVoicePattern string = "pelanggandengannomorantrian %s %s silahkankekonter %s"

	ConstSys string = "system"

	//TransactionType Const
	ConstAlokasi string = "1"
	ConstRetur   string = "2"
	ConstCancel  string = "3"

	//SodogiFlag Const
	ConstExist          string = "1"
	ConstNotExist       string = "0"
	ConstCompleteManual string = "2"

	//HistoryFlag Const
	ConstYes string = "1"
	ConstNo  string = "0"

	//SendFlag Const
	ConstBlmDikirim   string = "0"
	ConstSudahDikirim string = "1"
	ConstWaiting      string = "2"
	ConstSkip         string = "3"

	//DetailStatus Const
	ConstComplete    string = "1"
	ConstNotComplete string = "2"

	ConstSodogi                 string = "SODOGI"
	ConstDso                    string = "DSO"
	ConstGto                    string = "GTO"
	ConstCso                    string = "CSO"
	ConstRto                    string = "RTO"
	ConstRso                    string = "RSO"
	ConstTabelDetailTransaction string = "detail_transaction"
	ConstTabelSummary           string = "sodogi_summary"
	ConstTabelMasterFile        string = "master_file"
	ConstNull                   string = ""
	ConstAreaName               string = "GALLERY"
	ConstChannel                string = "50"
	ConstCustGroup              string = "6"
	ConstClusterNoGto           string = "0"
	ConstIsValid                string = "1"
	ConstNotValid               string = "0"

	ConstEncyprtPass			string = "eClu73r4ppS"

	ConstStringZero string = "0"
	ConstStringOne  string = "1"
)
