package util

import "github.com/adwitiyaio/arka/logger"

const firebaseDeepLinkApiKey = "FIREBASE_DEEPLINK_API_KEY"
const firebaseDeepLinkShortLinksUrls = "FIREBASE_DEEPLINK_SHORT_LINKS_URL"
const firebaseDeepLinkDynamicLinkDomain = "FIREBASE_DEEPLINK_DYNAMIC_LINK_DOMAIN"
const firebaseDeepLinkAndroidPackageName = "FIREBASE_DEEPLINK_ANDROID_PACKAGE_NAME"
const firebaseDeepLinkIosBundleId = "FIREBASE_DEEPLINK_IOS_BUNDLE_ID"
const firebaseDeepLinkIosStoreId = "FIREBASE_DEEPLINK_IOS_STORE_ID"
const firebaseDeepLinkSocialTitle = "FIREBASE_DEEPLINK_SOCIAL_TITLE"
const firebaseDeepLinkSocialDescription = "FIREBASE_DEEPLINK_SOCIAL_DESCRIPTION"
const firebaseDeepLinkSocialImageUrl = "FIREBASE_DEEPLINK_SOCIAL_IMAGE_URL"

type firebaseDeepLinkProvider struct {
	shortsLinkUrl      string
	apiKey             string
	dynamicLinkDomain  string
	androidPackageName string
	iosBundleId        string
	iosStoreId         string
	socialTitle        string
	socialDescription  string
	socialImageUrl     string
}

type firebaseDeepLinkLinkResponse struct {
	shortLink   string
	previewLink string
}

type firebaseDeepLinkDynamicLinkInfo struct {
	DomainURIPrefix   string                            `json:"domainUriPrefix"`
	Link              string                            `json:"link"`
	AndroidInfo       firebaseDeepLinkAndroidInfo       `json:"androidInfo"`
	IosInfo           firebaseDeepLinkIosInfo           `json:"iosInfo"`
	SocialMetaTagInfo firebaseDeepLinkSocialMetaTagInfo `json:"socialMetaTagInfo"`
}

type firebaseDeepLinkAndroidInfo struct {
	AndroidPackageName string `json:"androidPackageName"`
}

type firebaseDeepLinkIosInfo struct {
	IosBundleID   string `json:"iosBundleId"`
	IosAppStoreID string `json:"iosAppStoreId"`
}

type firebaseDeepLinkSocialMetaTagInfo struct {
	SocialTitle       string `json:"socialTitle"`
	SocialDescription string `json:"socialDescription"`
	SocialImageLink   string `json:"socialImageLink"`
}

type firebaseDeepLinkRequest struct {
	DynamicLinkInfo firebaseDeepLinkDynamicLinkInfo `json:"dynamicLinkInfo"`
}

func (mus *multiUrlManager) initializeFirebase() {
	apiKey := mus.sm.GetValueForKey(firebaseDeepLinkApiKey)
	shortsLinkUrl := mus.sm.GetValueForKey(firebaseDeepLinkShortLinksUrls)
	dynamicLinkDomain := mus.sm.GetValueForKey(firebaseDeepLinkDynamicLinkDomain)
	androidPackageName := mus.sm.GetValueForKey(firebaseDeepLinkAndroidPackageName)
	iosBundleId := mus.sm.GetValueForKey(firebaseDeepLinkIosBundleId)
	iosStoreId := mus.sm.GetValueForKey(firebaseDeepLinkIosStoreId)
	socialTitle := mus.sm.GetValueForKey(firebaseDeepLinkSocialTitle)
	socialDescription := mus.sm.GetValueForKey(firebaseDeepLinkSocialDescription)
	socialImageUrl := mus.sm.GetValueForKey(firebaseDeepLinkSocialImageUrl)
	mus.fbp = &firebaseDeepLinkProvider{
		shortsLinkUrl:      shortsLinkUrl,
		apiKey:             apiKey,
		dynamicLinkDomain:  dynamicLinkDomain,
		androidPackageName: androidPackageName,
		iosBundleId:        iosBundleId,
		iosStoreId:         iosStoreId,
		socialTitle:        socialTitle,
		socialDescription:  socialDescription,
		socialImageUrl:     socialImageUrl,
	}
}

func (mus multiUrlManager) createDeepLinkWithFirebase(url string) (string, error) {
	reqBody := firebaseDeepLinkRequest{DynamicLinkInfo: firebaseDeepLinkDynamicLinkInfo{
		DomainURIPrefix: mus.fbp.dynamicLinkDomain,
		Link:            url,
		AndroidInfo:     firebaseDeepLinkAndroidInfo{AndroidPackageName: mus.fbp.androidPackageName},
		IosInfo: firebaseDeepLinkIosInfo{
			IosBundleID:   mus.fbp.iosBundleId,
			IosAppStoreID: mus.fbp.iosStoreId,
		},
		SocialMetaTagInfo: firebaseDeepLinkSocialMetaTagInfo{
			SocialTitle:       mus.fbp.socialTitle,
			SocialDescription: mus.fbp.socialDescription,
			SocialImageLink:   mus.fbp.socialImageUrl,
		},
	}}

	var response firebaseDeepLinkLinkResponse
	resp, err := mus.client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("key", mus.fbp.apiKey).
		SetBody(reqBody).
		SetResult(&response).
		Post(mus.fbp.shortsLinkUrl)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to create deep link with firebase")
		return "", err
	}

	logger.Log.Debug().Msgf("http response -> %s", string(resp.Body()))
	logger.Log.Debug().Msgf("firebase response preview link -> %s", response.previewLink)
	return response.shortLink, nil
}
