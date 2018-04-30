package goinsta

// Media is media object representation
//
// Do not use concurrently
type Media struct {
	insta *Instagram

	ID       string
	Info     MediaInfo
	Comments MediaComments
	Likers   MediaLikers
}

func NewMedia(insta *Instagram) {
	return Media{insta: insta}
}

// Get gets all media data (Likes, Likers, Comments)
func (media *Media) Get(mediaID string) error {
	media.ID = mediaID

	media.Comments()
	media.Likers()
	media.Likes()
}

// Comments collect all possible comments from media
func (media *Media) Comments() error {
	insta := media.insta
	mediaID := media.ID

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("media/%s/comments", mediaID))
	req.args.Set("max_id", maxID)

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &media.Comments)
	return err
}

// MediaLikers return likers of a media , input is mediaid of a media
func (media *Media) Likers() error {
	body, err := insta.sendSimpleRequest("media/%s/likers/", mediaID)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &media.Likers)
	return err
}

// MediaInfo return media information
func (media *Media) Info() error {
	insta := media.insta
	result := MediaInfoResponse{}
	mediaID := media.ID

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint(fmt.Sprintf("media/%s/info", mediaID))

	data, err := insta.prepareData(map[string]interface{}{
		"media_id": mediaID,
	})
	if err != nil {
		return result, err
	}
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)

	return result, err
}