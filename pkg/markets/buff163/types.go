package buff163

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/arseniizyk/investor1337/pkg/markets"
	"go.uber.org/zap"
)

type buff163 struct {
	client *http.Client
	items  map[string]int
	l      *zap.Logger
}

type Response struct {
	Code string `json:"code"`
	Data struct {
		PageNum    int `json:"page_num"`
		PageSize   int `json:"page_size"`
		TotalPage  int `json:"total_page"`
		TotalCount int `json:"total_count"`
		Items      []struct {
			Appid       int    `json:"appid"`
			Game        string `json:"game"`
			ID          string `json:"id"`
			UserID      string `json:"user_id"`
			UserSteamid string `json:"user_steamid"`
			AssetInfo   struct {
				Appid      int    `json:"appid"`
				Contextid  int    `json:"contextid"`
				Assetid    string `json:"assetid"`
				Classid    string `json:"classid"`
				Instanceid string `json:"instanceid"`
				GoodsID    int    `json:"goods_id"`
				Paintwear  string `json:"paintwear"`
				ActionLink string `json:"action_link"`
				ID         string `json:"id"`
				Info       struct {
					IconURL         string `json:"icon_url"`
					Stickers        []any  `json:"stickers"`
					Keychains       []any  `json:"keychains"`
					Fraudwarnings   any    `json:"fraudwarnings"`
					TournamentTags  []any  `json:"tournament_tags"`
					OriginalIconURL string `json:"original_icon_url"`
				} `json:"info"`
				HasTradableCooldown  bool   `json:"has_tradable_cooldown"`
				TradableUnfrozenTime any    `json:"tradable_unfrozen_time"`
				TradableCooldownText string `json:"tradable_cooldown_text"`
			} `json:"asset_info"`
			Price                 string `json:"price"`
			Fee                   string `json:"fee"`
			State                 int    `json:"state"`
			Mode                  int    `json:"mode"`
			GoodsID               int    `json:"goods_id"`
			CreatedAt             int    `json:"created_at"`
			UpdatedAt             int    `json:"updated_at"`
			RecentDeliverRate     any    `json:"recent_deliver_rate"`
			RecentAverageDuration any    `json:"recent_average_duration"`
			Featured              int    `json:"featured"`
			Description           string `json:"description"`
			TradableCooldown      any    `json:"tradable_cooldown"`
			AllowBargain          bool   `json:"allow_bargain"`
			Income                string `json:"income"`
			AllowBargainChat      bool   `json:"allow_bargain_chat"`
			StickerPremium        any    `json:"sticker_premium"`
			OrderType             int    `json:"order_type"`
			AllowBargainSwap      bool   `json:"allow_bargain_swap"`
			StateText             string `json:"state_text"`
			Bookmarked            bool   `json:"bookmarked"`
			CanBargain            bool   `json:"can_bargain"`
			CannotBargainReason   string `json:"cannot_bargain_reason"`
			SupportedPayMethods   []int  `json:"supported_pay_methods"`
			CanBargainSwap        bool   `json:"can_bargain_swap"`
			CanBargainChat        bool   `json:"can_bargain_chat"`
			CanUseInspectTrnURL   bool   `json:"can_use_inspect_trn_url"`
			BackgroundImageURL    string `json:"background_image_url"`
			ImgSrc                string `json:"img_src"`
		} `json:"items"`
		GoodsInfos struct {
			Num886606 struct {
				Appid          int    `json:"appid"`
				Game           string `json:"game"`
				GoodsID        int    `json:"goods_id"`
				IconURL        string `json:"icon_url"`
				ItemID         any    `json:"item_id"`
				MarketHashName string `json:"market_hash_name"`
				Name           string `json:"name"`
				SteamPrice     string `json:"steam_price"`
				SteamPriceCny  string `json:"steam_price_cny"`
				Tags           struct {
					Type struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"type"`
					Rarity struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"rarity"`
					Itemset struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"itemset"`
					Quality struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"quality"`
					Category struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"category"`
					CategoryGroup struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"category_group"`
					Weaponcase struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"weaponcase"`
					WeaponcaseType struct {
						ID            int    `json:"id"`
						Category      string `json:"category"`
						InternalName  string `json:"internal_name"`
						LocalizedName string `json:"localized_name"`
					} `json:"weaponcase_type"`
				} `json:"tags"`
				OriginalIconURL    string `json:"original_icon_url"`
				MarketMinPrice     string `json:"market_min_price"`
				Description        any    `json:"description"`
				ShortName          string `json:"short_name"`
				SellMinPrice       string `json:"sell_min_price"`
				SellReferencePrice string `json:"sell_reference_price"`
				IsCharm            bool   `json:"is_charm"`
				KeychainColorImg   any    `json:"keychain_color_img"`
				CanInspect         bool   `json:"can_inspect"`
				Can3DInspect       bool   `json:"can_3d_inspect"`
			} `json:"886606"`
		} `json:"goods_infos"`
		UserInfos struct {
			U1076702943 struct {
				UserID       string `json:"user_id"`
				Avatar       string `json:"avatar"`
				ShopID       string `json:"shop_id"`
				AvatarSafe   string `json:"avatar_safe"`
				Nickname     string `json:"nickname"`
				IsAutoAccept bool   `json:"is_auto_accept"`
				SellerLevel  int    `json:"seller_level"`
				VTypes       any    `json:"v_types"`
				IsPremiumVip bool   `json:"is_premium_vip"`
			} `json:"U1076702943"`
			U1080178344 struct {
				UserID       string `json:"user_id"`
				Avatar       string `json:"avatar"`
				ShopID       string `json:"shop_id"`
				AvatarSafe   string `json:"avatar_safe"`
				Nickname     string `json:"nickname"`
				IsAutoAccept bool   `json:"is_auto_accept"`
				SellerLevel  int    `json:"seller_level"`
				VTypes       any    `json:"v_types"`
				IsPremiumVip bool   `json:"is_premium_vip"`
			} `json:"U1080178344"`
			U1102576054 struct {
				UserID       string `json:"user_id"`
				Avatar       string `json:"avatar"`
				ShopID       string `json:"shop_id"`
				AvatarSafe   string `json:"avatar_safe"`
				Nickname     string `json:"nickname"`
				IsAutoAccept bool   `json:"is_auto_accept"`
				SellerLevel  int    `json:"seller_level"`
				VTypes       any    `json:"v_types"`
				IsPremiumVip bool   `json:"is_premium_vip"`
			} `json:"U1102576054"`
		} `json:"user_infos"`
		HasMarketStores struct {
			U1102576054 bool `json:"U1102576054"`
			U1080178344 bool `json:"U1080178344"`
			U1076702943 bool `json:"U1076702943"`
		} `json:"has_market_stores"`
		SortBy             string `json:"sort_by"`
		SrcURLBackground   string `json:"src_url_background"`
		PreviewScreenshots struct {
			Selling bool   `json:"selling"`
			BgImg   string `json:"bg_img"`
		} `json:"preview_screenshots"`
		FopStr            string `json:"fop_str"`
		ShowPayMethodIcon bool   `json:"show_pay_method_icon"`
		ShowGameCmsIcon   bool   `json:"show_game_cms_icon"`
		GoodsInfoIconURL  string `json:"goods_info_icon_url"`
	} `json:"data"`
	Msg any `json:"msg"`
}

func New(c *http.Client, l *zap.Logger) (markets.Market, error) {
	buff := buff163{
		client: c,
		l:      l,
	}

	if err := buff.loadItems(); err != nil {
		return nil, err
	}

	return buff, nil
}

func (buff *buff163) loadItems() error {
	b, err := os.ReadFile("../buff163ids.json")
	if err != nil {
		buff.l.Error("Cant load buff163 ids from json", zap.Error(err))
		return err
	}

	data := make(map[string]int)

	if err := json.Unmarshal(b, &data); err != nil {
		buff.l.Error("Cant unmarshal buff163 ids", zap.Error(err))
		return err
	}

	buff.items = make(map[string]int, len(data))
	for k, v := range data {
		buff.items[strings.ToLower(k)] = v
	}

	return nil
}
