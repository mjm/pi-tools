package backup

var pgDatabases = []postgresDatabase{
	{
		Name:        "trips",
		vaultPolicy: "presence",
	},
	{
		Name:   "go-links",
		dbName: "go_links",
	},
	//{
	//	Name:   "homebase-bot",
	//	dbName: "homebase_bot",
	//},
	{
		Name: "grafana",
	},
	{
		Name: "paperless",
	},
}

var phabricatorDatabases = []string{
	"phabricator_almanac",
	"phabricator_application",
	"phabricator_audit",
	"phabricator_auth",
	"phabricator_badges",
	"phabricator_cache",
	"phabricator_calendar",
	"phabricator_chatlog",
	"phabricator_conduit",
	"phabricator_config",
	"phabricator_conpherence",
	"phabricator_countdown",
	"phabricator_daemon",
	"phabricator_dashboard",
	"phabricator_differential",
	"phabricator_diviner",
	"phabricator_doorkeeper",
	"phabricator_draft",
	"phabricator_drydock",
	"phabricator_fact",
	"phabricator_feed",
	"phabricator_file",
	"phabricator_flag",
	"phabricator_fund",
	"phabricator_harbormaster",
	"phabricator_herald",
	"phabricator_legalpad",
	"phabricator_maniphest",
	"phabricator_meta_data",
	"phabricator_metamta",
	"phabricator_multimeter",
	"phabricator_nuance",
	"phabricator_oauth_server",
	"phabricator_owners",
	"phabricator_packages",
	"phabricator_passphrase",
	"phabricator_paste",
	"phabricator_pastebin",
	"phabricator_phame",
	"phabricator_phlux",
	"phabricator_pholio",
	"phabricator_phortune",
	"phabricator_phragment",
	"phabricator_phrequent",
	"phabricator_phriction",
	"phabricator_phurl",
	"phabricator_policy",
	"phabricator_ponder",
	"phabricator_project",
	"phabricator_releeph",
	"phabricator_repository",
	"phabricator_search",
	"phabricator_slowvote",
	"phabricator_spaces",
	"phabricator_system",
	"phabricator_token",
	"phabricator_user",
	"phabricator_worker",
	"phabricator_xhpast",
	"phabricator_xhprof",
}

type postgresDatabase struct {
	Name        string
	dbName      string
	vaultPolicy string
}

func (d postgresDatabase) DBName() string {
	if d.dbName == "" {
		return d.Name
	}
	return d.dbName
}

func (d postgresDatabase) VaultPolicy() string {
	if d.vaultPolicy == "" {
		return d.Name
	}
	return d.vaultPolicy
}
