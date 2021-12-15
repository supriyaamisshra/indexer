package fetcher

const (
	RARIBLE    = "Rarible"
	CONTEXT    = "Context"
	CONVO      = "Convo"
	TWITTER    = "Twitter"
	OPENSEA    = "Opensea"
	ZORA       = "Zora"
	FOUNDATION = "Foundation"
	SHOWTIME   = "Showtime"
	SYBIL      = "Sybil"
	SUPERRARE  = "Superrare"
	INFURA     = "Infura"
	POAP       = "Poap"
)

const (
	SuperrareContractAddress  = "0x41a322b28d0ff354040e2cbc676f0320d8c8850d"
	OpenSeaContractAddress    = "0x495f947276749ce646f68ac8c248420045cb7b5e"
	RaribleContractAddress    = "0xd07dc4262bcdbf85190c01c996b4c06a461d2430"
	FoundationContractAddress = "0x3b3ee1931dc30c1957379fac9aba94d1c48a5405"
	ZoraContractAddress       = "0xabefbc9fd2f806065b4f3c237d4b59d9a97bcac7"
	ContextContractAddress    = "ctx"
)

const (
	ContextUrl          = "https://context.app/api/profile/%s"
	SuperrareUrl        = "https://superrare.com/api/v2/user?address=%s"
	RaribleFollowingUrl = "https://api-mainnet.rarible.com/marketplace/api/v4/followings?owner=%s"
	RaribleFollowerUrl  = "https://api-mainnet.rarible.com/marketplace/api/v4/followers?user=%s"
	SybilUrl            = "https://raw.githubusercontent.com/Uniswap/sybil-list/master/verified.json"
	PoapTokensFetchUrl  = "https://api.poap.xyz/actions/scan/%s"
	PoapXyzGraphUrl     = "https://api.thegraph.com/subgraphs/name/poap-xyz/poap"
)

type ConnectionEntryList struct {
	Conn []ConnectionEntry
	Err  error
	msg  string
}
type ConnectionEntry struct {
	From     string
	To       string
	Platform string
}

type IdentityEntryList struct {
	OpenSea    []UserOpenSeaIdentity
	Twitter    []UserTwitterIdentity
	Superrare  []UserSuperrareIdentity
	Rarible    []UserRaribleIdentity
	Context    []UserContextIdentity
	Zora       []UserZoraIdentity
	Foundation []UserFoundationIdentity
	Showtime   []UserShowtimeIdentity
	Poaps      []UserPoapIdentity
	Ens        string
}

type IdentityEntry struct {
	OpenSea    *UserOpenSeaIdentity
	Twitter    *UserTwitterIdentity
	Superrare  *UserSuperrareIdentity
	Rarible    *UserRaribleIdentity
	Context    *UserContextIdentity
	Zora       *UserZoraIdentity
	Ens        *UserEnsIdentity
	Foundation *UserFoundationIdentity
	Showtime   *UserShowtimeIdentity
	Poap       *UserPoapIdentity
	Err        error
	Msg        string
}

type UserTwitterIdentity struct {
	Handle     string
	DataSource string
}

type UserRaribleIdentity struct {
	Username        string
	Homepage        string
	ItemSold        int
	AmountSoldInEth float64
	DataSource      string
}

type UserOpenSeaIdentity struct {
	Username   string
	Homepage   string
	DataSource string
}

type UserEnsIdentity struct {
	Ens        string
	DataSource string
}

type UserContextIdentity struct {
	FollowerCount int
	Username      string
	Website       string
	DataSource    string
}

type UserSuperrareIdentity struct {
	Username       string
	Homepage       string
	Location       string
	Bio            string
	InstagramLink  string
	TwitterLink    string
	SteemitLink    string
	Website        string
	SpotifyLink    string
	SoundCloudLink string
	DataSource     string
}

type UserFoundationIdentity struct {
	Username   string
	Bio        string
	Tiktok     string
	Twitch     string
	Discord    string
	Twitter    string
	Website    string
	Youtube    string
	Facebook   string
	Snapchat   string
	Instagram  string
	DataSource string
}

type UserZoraIdentity struct {
	Username   string
	Website    string
	DataSource string
}

type UserShowtimeIdentity struct {
	Name             string
	Username         string
	Bio              string
	TwitterHandle    string
	LinkTreeHandle   string
	CryptoArtHandle  string
	FoundationHandle string
	HicetnuncHandle  string
	OpenseaHandle    string
	RaribleHandle    string
	DataSource       string
}

type UserPoapIdentity struct {
	EventID         string
	EventDesc       string
	TokenID         string
	Recommendations []PoapRecommendation
}

type PoapRecommendation struct {
	Address string
	EventID string
	TokenID string
}

// DTOs below

type RaribleConnectionResp struct {
	Following struct {
		From string `json:"owner"`
		To   string `json:"user"`
	} `json:"following"`
}

type ContextAppResp struct {
	FollowerCount int               `json:"followerCount"`
	Ens           map[string]string `json:"ens"`
	Profiles      map[string]([]struct {
		Address  string `json:"address,omitempty"`
		Contract string `json:"contract,omitempty"`
		Url      string `json:"url,omitempty"`
		Website  string `json:"website,omitempty"`
		Username string `json:"username,omitempty"`
	}) `json:"profiles"`
}

type ContextConnection struct {
	Relationships []struct {
		Actor string `json:"actor"`
	} `json:"relationships"`
	Profiles map[string]([]struct {
		Address string `json:"address"`
	}) `json:"profiles"`
}

type SuperrareProfile struct {
	Result struct {
		Username       string `json:"username"`
		Location       string `json:"location"`
		Bio            string `json:"bio"`
		InstagramLink  string `json:"instagramLink"`
		TwitterLink    string `json:"twitterLink"`
		SteemitLink    string `json:"steemitLink"`
		Website        string `json:"website"`
		SpotifyLink    string `json:"spotifyLink"`
		SoundCloudLink string `json:"soundcloudLink"`
	} `json:"result"`
}

type SybilVerifiedList map[string]struct {
	Twitter struct {
		Timestamp uint64 `json:"timestamp"`
		TweetId   string `json:"tweetID"`
		Handle    string `json:"handle"`
	} `json:"twitter"`
}

type FoundationIdentity struct {
	Data struct {
		User struct {
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Links    struct {
				Tiktok struct {
					Handle string `json:"handle"`
				} `json:"tiktok"`
				Twitch struct {
					Handle string `json:"handle"`
				} `json:"twitch"`
				Discord struct {
					Handle string `json:"handle"`
				} `json:"discord"`
				Twitter struct {
					Handle string `json:"handle"`
				} `json:"twitter"`
				Website struct {
					Handle string `json:"handle"`
				} `json:"website"`
				Youtube struct {
					Handle string `json:"handle"`
				} `json:"youtube"`
				Facebook struct {
					Handle string `json:"handle"`
				} `json:"facebook"`
				Snapchat struct {
					Handle string `json:"handle"`
				} `json:"snapchat"`
				Instagram struct {
					Handle string `json:"handle"`
				} `json:"instagram"`
			} `json:"links"`
			TwitSocialVerifs []struct {
				Username string `json:"username"`
			} `json:"twitSocialVerifs"`
			InstaSocialVerifs []struct {
				Username string `json:"username"`
			} `json:"instaSocialVerifs"`
		} `json:"user"`
	} `json:"data"`
}

type PoapActionScanResultEntry struct {
	Event struct {
		Id          uint64  `json:"id"`
		FancyId     string  `json:"fancy_id"`
		Name        string  `json:"name"`
		EventUrl    string  `json:"event_url"`
		ImageUrl    string  `json:"image_url"`
		Country     string  `json:"country"`
		City        string  `json:"city"`
		Description string  `json:"description"`
		Year        uint32  `json:"year"`
		StartDate   string  `json:"start_date"`
		EndDate     string  `json:"end_date"`
		ExpiryDate  string  `json:"expiry_date"`
		CreatedDate string  `json:"created_date"`
		Supply      *uint64 `json:"supply"`
	} `json:"event"`

	TokenID string  `json:"tokenId"`
	Owner   string  `json:"owner"`
	Supply  uint64  `json:"supply"`
	Created *string `json:"created"`
}

type PoapGraphTokenOwnerData struct {
	Id string `json:"id"`
}

type PoapGraphEventTokensData struct {
	Tokens []struct {
		Id    string                  `json:"id"`
		Owner PoapGraphTokenOwnerData `json:"owner"`
	} `json:"tokens"`
}

type PoapGraphEventData struct {
	Event PoapGraphEventTokensData `json:"event"`
}

type PoapTokenAndOwnerGraphQueryResult struct {
	Data PoapGraphEventData `json:"data"`
}
