package request

type AuthZaloBaseReq struct {
	Channel     string `json:"channel"      binding:"required"` // Game distribution channel. ( "h5zl": open in H5, "inapp": open game inapp )
	UID         string `json:"uid"          binding:"required"` // ID of user in Zalo platform.
	Timestamp   string `json:"timestamp"    binding:"required"` // The time when the user logs into the game. Partners can use this parameter to set the expire time.
	UtmSource   string `json:"utm_source"`                      // The source of the traffic, e.g., "google", "facebook", "newsletter".
	UtmMedium   string `json:"utm_medium"`                      // The medium or channel, e.g., "cpc", "email", "social".
	UtmCampaign string `json:"utm_campaign"`                    // The specific campaign name or identifier, e.g., "spring_sale", "launch_offer".
	Sign        string `json:"sign"         binding:"required"` // The md5 signature to verify the request is from Zalo platform.
}

type AuthZaloReq struct {
	AuthZaloBaseReq

	DisplayName string `json:"display_name"` // Display name of user in Zalo platform.
	AvatarURL   string `json:"avatar_url"`   // Avatar URL of user in Zalo platform.
}

func (req *AuthZaloBaseReq) ToMap() map[string]string {
	return map[string]string{
		"channel":      req.Channel,
		"uid":          req.UID,
		"timestamp":    req.Timestamp,
		"utm_source":   req.UtmSource,
		"utm_medium":   req.UtmMedium,
		"utm_campaign": req.UtmCampaign,
		"sign":         req.Sign,
	}
}
