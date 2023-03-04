package spstruct

import "time"

type FitBossPostData struct {
	Username  string                   `json:"username"`
	Type      string                   `json:"type"`
	Team      []map[string]interface{} `json:"team"`
	Boosts    []interface{}            `json:"boosts"`
	DeckPower int                      `json:"deckPower"`
	Memo      string                   `json:"memo"`
}

type FitBossRequestsData struct {
	Message struct {
		Method string `json:"method"`
		Params struct {
			DocumentURL    string `json:"documentURL"`
			FrameID        string `json:"frameId"`
			HasUserGesture bool   `json:"hasUserGesture"`
			Initiator      struct {
				Stack struct {
					CallFrames []struct {
						ColumnNumber int    `json:"columnNumber"`
						FunctionName string `json:"functionName"`
						LineNumber   int    `json:"lineNumber"`
						ScriptID     string `json:"scriptId"`
						URL          string `json:"url"`
					} `json:"callFrames"`
				} `json:"stack"`
				Type string `json:"type"`
			} `json:"initiator"`
			LoaderID             string `json:"loaderId"`
			RedirectHasExtraInfo bool   `json:"redirectHasExtraInfo"`
			Request              struct {
				HasPostData bool `json:"hasPostData"`
				Headers     struct {
					Accept          string `json:"Accept"`
					ContentType     string `json:"Content-Type"`
					Referer         string `json:"Referer"`
					UserAgent       string `json:"User-Agent"`
					SecChUa         string `json:"sec-ch-ua"`
					SecChUaMobile   string `json:"sec-ch-ua-mobile"`
					SecChUaPlatform string `json:"sec-ch-ua-platform"`
				} `json:"headers"`
				InitialPriority  string `json:"initialPriority"`
				IsSameSite       bool   `json:"isSameSite"`
				Method           string `json:"method"`
				MixedContentType string `json:"mixedContentType"`
				PostData         string `json:"postData"`
				PostDataEntries  []struct {
					Bytes string `json:"bytes"`
				} `json:"postDataEntries"`
				ReferrerPolicy string `json:"referrerPolicy"`
				URL            string `json:"url"`
			} `json:"request"`
			RequestID string  `json:"requestId"`
			Timestamp float64 `json:"timestamp"`
			Type      string  `json:"type"`
			WallTime  float64 `json:"wallTime"`
		} `json:"params"`
	} `json:"message"`
	Webview string `json:"webview"`
}

type FitReturnData struct {
	Date   int64  `json:"date"`
	Player string `json:"player"`
	Boss   string `json:"boss"`
	//Team     []map[string]interface{} `json:"team"`
	//Actions  []map[string]interface{} `json:"actions"`
	TotalDmg int        `json:"totalDmg"` // 总伤害
	Points   int        `json:"points"`   // 分数
	Rewards  []struct { // 奖励
		Type string  `json:"type"`
		Name string  `json:"name"`
		Qty  float64 `json:"qty"`
	} `json:"rewards"`
	UniqueRules []string `json:"uniqueRules"`
	Ppd         int      `json:"ppd"`
	ID          string   `json:"id"`
	NewRules    struct {
		Message string `json:"message"`
		Rules   struct {
			Active bool     `json:"active"`
			BossID string   `json:"boss_id"`
			Rules  []string `json:"rules"`
		} `json:"rules"`
	} `json:"newRules"`
}

type KeyLoginResData struct {
	Stamina struct {
		Last    string `json:"last"`
		Current int    `json:"current"`
		Max     int    `json:"max"`
	} `json:"stamina"`
	Sc struct {
		Balance float64 `json:"balance"`
	} `json:"sc"`
	Heroes struct {
		Warrior struct {
			Weapon   string `json:"weapon"`
			Offhand  string `json:"offhand"`
			Head     string `json:"head"`
			Necklace string `json:"necklace"`
			Body     string `json:"body"`
			Hands    string `json:"hands"`
			Ring     string `json:"ring"`
			Legs     string `json:"legs"`
			Feet     string `json:"feet"`
			Back     string `json:"back"`
		} `json:"warrior"`
		Wizard struct {
			Weapon   string `json:"weapon"`
			Offhand  string `json:"offhand"`
			Head     string `json:"head"`
			Necklace string `json:"necklace"`
			Body     string `json:"body"`
			Hands    string `json:"hands"`
			Ring     string `json:"ring"`
			Legs     string `json:"legs"`
			Feet     string `json:"feet"`
			Back     string `json:"back"`
		} `json:"wizard"`
		Ranger struct {
			Weapon   string `json:"weapon"`
			Offhand  string `json:"offhand"`
			Head     string `json:"head"`
			Necklace string `json:"necklace"`
			Body     string `json:"body"`
			Hands    string `json:"hands"`
			Ring     string `json:"ring"`
			Legs     string `json:"legs"`
			Feet     string `json:"feet"`
			Back     string `json:"back"`
		} `json:"ranger"`
	} `json:"heroes"`
	CreatedDate  time.Time `json:"createdDate"`
	Username     string    `json:"username"`
	InGameAssets []struct {
		ID          string `json:"_id"`
		Name        string `json:"name"`
		Symbol      string `json:"symbol,omitempty"`
		Type        string `json:"type"`
		Qty         int    `json:"qty"`
		Max         int    `json:"max"`
		Description string `json:"description"`
	} `json:"inGameAssets"`
	LastLogin time.Time `json:"lastLogin"`
	Token     string    `json:"token"`
	Stats     struct {
		Logins     []int64 `json:"logins"`
		BossFights int     `json:"bossFights"`
		BronzeDmg  int     `json:"bronze_dmg"`
		SilverDmg  int     `json:"silver_dmg"`
		GoldDmg    int     `json:"gold_dmg"`
		DiamondDmg int     `json:"diamond_dmg"`
	} `json:"stats"`
	CurrentClass string `json:"currentClass"`
	Settings     struct {
		LogSpeed int    `json:"logSpeed"`
		Tutorial []bool `json:"tutorial"`
		Frame    string `json:"frame"`
	} `json:"settings"`
	Electrum    int `json:"electrum"`
	UniqueRules struct {
		T1 struct {
			Active bool     `json:"active"`
			BossID string   `json:"boss_id"`
			Rules  []string `json:"rules"`
		} `json:"t1"`
		T2 struct {
			Active bool     `json:"active"`
			BossID string   `json:"boss_id"`
			Rules  []string `json:"rules"`
		} `json:"t2"`
	} `json:"uniqueRules"`
	ID         string `json:"id"`
	ServerTime int64  `json:"serverTime"`
	Config     struct {
		ActiveUser       string  `json:"activeUser"`
		PostingUser      string  `json:"postingUser"`
		CollectionURL    string  `json:"collectionURL"`
		TokenPricesURL   string  `json:"tokenPricesURL"`
		SfAPI            string  `json:"sfAPI"`
		Mode             string  `json:"mode"`
		MaxMana          int     `json:"maxMana"`
		MaxMonsters      int     `json:"maxMonsters"`
		MaxLevelsPerTier [][]int `json:"maxLevelsPerTier"`
		SocketItems      []struct {
			Type  string `json:"type"`
			Bonus []int  `json:"bonus"`
		} `json:"socketItems"`
		StatLevels       [][]int `json:"statLevels"`
		SfCardLevels     [][]int `json:"sfCardLevels"`
		SfFoilCardLevels [][]int `json:"sfFoilCardLevels"`
		BurnTable        struct {
			Num1 int `json:"1"`
			Num2 int `json:"2"`
			Num3 int `json:"3"`
			Num4 int `json:"4"`
			Foil int `json:"foil"`
		} `json:"burnTable"`
		SocketOdds       [][]int `json:"socketOdds"`
		HeroStatMaxTiers []struct {
			Armor  int `json:"armor"`
			Health int `json:"health"`
			Speed  int `json:"speed"`
			Dmg    int `json:"dmg"`
		} `json:"heroStatMaxTiers"`
		ReplacedSkills []string `json:"replacedSkills"`
		AirdropCards   []struct {
			Cid         string   `json:"cid"`
			Player      string   `json:"player"`
			Name        string   `json:"name"`
			Type        string   `json:"type"`
			Slot        string   `json:"slot"`
			Stat        string   `json:"stat"`
			Rarity      int      `json:"rarity"`
			Combined    int      `json:"combined"`
			Sockets     int      `json:"sockets"`
			SocketItems []string `json:"socket_items"`
			Level       int      `json:"level"`
			Edition     int      `json:"edition"`
			Burnt       bool     `json:"burnt"`
			Equipped    bool     `json:"equipped"`
			ImgURL      string   `json:"imgURL"`
		} `json:"airdropCards"`
		GemAirdropCards []struct {
			Cid         string        `json:"cid"`
			Player      string        `json:"player"`
			Name        string        `json:"name"`
			Type        string        `json:"type"`
			Slot        string        `json:"slot"`
			Stat        string        `json:"stat"`
			Rarity      int           `json:"rarity"`
			Combined    int           `json:"combined"`
			Sockets     int           `json:"sockets"`
			SocketItems []interface{} `json:"socket_items"`
			Level       int           `json:"level"`
			Edition     int           `json:"edition"`
			Burnt       bool          `json:"burnt"`
			Equipped    bool          `json:"equipped"`
			ImgURL      string        `json:"imgURL"`
		} `json:"gemAirdropCards"`
		AirdropMilestones []int `json:"airdropMilestones"`
		AirdropGuarantees []int `json:"airdropGuarantees"`
		AlphaExclusive    struct {
			Cid         string   `json:"cid"`
			Player      string   `json:"player"`
			Name        string   `json:"name"`
			Type        string   `json:"type"`
			Slot        string   `json:"slot"`
			Stat        string   `json:"stat"`
			Rarity      int      `json:"rarity"`
			Combined    int      `json:"combined"`
			Sockets     int      `json:"sockets"`
			SocketItems []string `json:"socket_items"`
			Level       int      `json:"level"`
			Edition     int      `json:"edition"`
			Burnt       bool     `json:"burnt"`
			Equipped    bool     `json:"equipped"`
			ImgURL      string   `json:"imgURL"`
		} `json:"alphaExclusive"`
		AvailableAirdrops []string `json:"availableAirdrops"`
		Whitelist         []string `json:"whitelist"`
	} `json:"config"`
	Message string `json:"message"`
}

type KeyLoginPostData struct {
	Message struct {
		Method string `json:"method"`
		Params struct {
			DocumentURL    string `json:"documentURL"`
			FrameID        string `json:"frameId"`
			HasUserGesture bool   `json:"hasUserGesture"`
			Initiator      struct {
				Stack struct {
					CallFrames []struct {
						ColumnNumber int    `json:"columnNumber"`
						FunctionName string `json:"functionName"`
						LineNumber   int    `json:"lineNumber"`
						ScriptID     string `json:"scriptId"`
						URL          string `json:"url"`
					} `json:"callFrames"`
				} `json:"stack"`
				Type string `json:"type"`
			} `json:"initiator"`
			LoaderID             string `json:"loaderId"`
			RedirectHasExtraInfo bool   `json:"redirectHasExtraInfo"`
			Request              struct {
				HasPostData bool `json:"hasPostData"`
				Headers     struct {
					Accept          string `json:"Accept"`
					ContentType     string `json:"Content-Type"`
					Referer         string `json:"Referer"`
					UserAgent       string `json:"User-Agent"`
					SecChUa         string `json:"sec-ch-ua"`
					SecChUaMobile   string `json:"sec-ch-ua-mobile"`
					SecChUaPlatform string `json:"sec-ch-ua-platform"`
				} `json:"headers"`
				InitialPriority  string `json:"initialPriority"`
				IsSameSite       bool   `json:"isSameSite"`
				Method           string `json:"method"`
				MixedContentType string `json:"mixedContentType"`
				PostData         string `json:"postData"`
				PostDataEntries  []struct {
					Bytes string `json:"bytes"`
				} `json:"postDataEntries"`
				ReferrerPolicy string `json:"referrerPolicy"`
				URL            string `json:"url"`
			} `json:"request"`
			RequestID string  `json:"requestId"`
			Timestamp float64 `json:"timestamp"`
			Type      string  `json:"type"`
			WallTime  float64 `json:"wallTime"`
		} `json:"params"`
	} `json:"message"`
	Webview string `json:"webview"`
}

type MonsterId struct {
	PlayingMonstersID   string `json:"playingMonstersId"`
	PlayingMonstersName string `json:"playingMonstersName"`
	PlayingMontersDiv   string `json:"playingMontersDiv"`
}

type Summoners struct {
	PlayingSummonersDiv  string `json:"playingSummonersDiv"`
	PlayingSummonersID   string `json:"playingSummonersId"`
	PlayingSummonersName string `json:"playingSummonersName"`
}

type CardSelection struct {
	PlayingMonsters  []MonsterId `json:"PlayingMonsters"`
	PlayingSummoners []Summoners `json:"playingSummoners"`
}

type UserData struct {
	BossID            string          `json:"bossId"`
	CardSelection     []CardSelection `json:"cardSelection"`
	HeroesType        string          `json:"heroesType"`
	PostingKey        string          `json:"postingKey"`
	TimeSleepInMinute int             `json:"timeSleepInMinute"`
	UserName          string          `json:"userName"`
}
type Card struct {
	CardDetailID int `json:"card_detail_id"`
}
type CardList struct {
	Cards []Card `json:"cards"`
}
type BattleCardsRequestBody struct {
	BossName string `json:"bossName"`
	BossId   string `json:"bossId"`
	Team     []int  `json:"team"`
}
