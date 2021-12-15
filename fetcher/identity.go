package fetcher

import (
	"encoding/json"
	"fmt"

	"github.com/StevenACoffman/errgroup"
	"go.uber.org/zap"
)

type identityFetcherController struct {
	entry chan IdentityEntry
	done  chan int
}

func NewIdentityFetcherController() identityFetcherController {
	return identityFetcherController{make(chan IdentityEntry), make(chan int)}
}

func (ifc identityFetcherController) Add(identityEntry *IdentityEntry) {
	ifc.entry <- *identityEntry
}

func (ifc identityFetcherController) AddAndDone(identityEntry *IdentityEntry) {
	ifc.Add(identityEntry)
	ifc.Done()
}

func (ifc identityFetcherController) Done() {
	ifc.done <- 1
}

func (ifc identityFetcherController) Close() {
	close(ifc.entry)
	close(ifc.done)
}

func (f *fetcher) FetchIdentity(address string) (IdentityEntryList, error) {

	var identityArr IdentityEntryList
	fetcherController := NewIdentityFetcherController()

	var fetchers = []func(string, *identityFetcherController){
		f.processContext,
		f.processSuperrare,
		f.processSybil,
		f.processPoap,
	}
	for _, identityFetcher := range fetchers {
		go identityFetcher(address, &fetcherController)
	}

	var done int
	for {
		select {
		case entry := <-fetcherController.entry:
			if entry.Err != nil {
				zap.L().With(zap.Error(entry.Err)).Error("identity api error: " + entry.Msg)
				continue
			}
			if entry.OpenSea != nil {
				identityArr.OpenSea = append(identityArr.OpenSea, *entry.OpenSea)
			}
			if entry.Twitter != nil {
				entry.Twitter.Handle = convertTwitterHandle(entry.Twitter.Handle)
				identityArr.Twitter = append(identityArr.Twitter, *entry.Twitter)
			}
			if entry.Superrare != nil {
				identityArr.Superrare = append(identityArr.Superrare, *entry.Superrare)
			}
			if entry.Rarible != nil {
				identityArr.Rarible = append(identityArr.Rarible, *entry.Rarible)
			}
			if entry.Context != nil {
				identityArr.Context = append(identityArr.Context, *entry.Context)
			}
			if entry.Zora != nil {
				identityArr.Zora = append(identityArr.Zora, *entry.Zora)
			}
			if entry.Foundation != nil {
				identityArr.Foundation = append(identityArr.Foundation, *entry.Foundation)
			}
			if entry.Showtime != nil {
				identityArr.Showtime = append(identityArr.Showtime, *entry.Showtime)
			}
			if entry.Ens != nil {
				identityArr.Ens = entry.Ens.Ens
			}
			if entry.Poap != nil {
				identityArr.Poaps = append(identityArr.Poaps, *entry.Poap)
			}
		case <-fetcherController.done:
			done += 1
		default:
			if done == len(fetchers) {
				fetcherController.Close()
				return identityArr, nil
			}
		}
	}
}

func (f *fetcher) processContext(address string, controller *identityFetcherController) {
	var result IdentityEntry

	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    fmt.Sprintf(ContextUrl, address),
		method: "GET",
	})
	if err != nil {
		result.Err = err
		result.Msg = "[processContext] fetch identity failed"
		controller.AddAndDone(&result)
		return
	}
	contextProfile := ContextAppResp{}
	err = json.Unmarshal(body, &contextProfile)
	if err != nil {
		result.Err = err
		result.Msg = "[processContext] identity response json unmarshal failed"
		controller.AddAndDone(&result)
		return
	}

	if value, ok := contextProfile.Ens[address]; ok {
		result.Ens = &UserEnsIdentity{
			Ens:        value,
			DataSource: CONTEXT,
		}
	}

	for _, profileList := range contextProfile.Profiles {
		for _, entry := range profileList {
			switch entry.Contract {
			case SuperrareContractAddress:
				result.Superrare = &UserSuperrareIdentity{
					Homepage:   entry.Url,
					Username:   entry.Username,
					DataSource: CONTEXT,
				}
			case OpenSeaContractAddress:
				result.OpenSea = &UserOpenSeaIdentity{
					Homepage:   entry.Url,
					Username:   entry.Username,
					DataSource: CONTEXT,
				}
			case RaribleContractAddress:
				result.Rarible = &UserRaribleIdentity{
					Homepage:   entry.Url,
					Username:   entry.Username,
					DataSource: CONTEXT,
				}
			case FoundationContractAddress:
				result.Foundation = &UserFoundationIdentity{
					Website:    entry.Website,
					Username:   entry.Username,
					DataSource: CONTEXT,
				}
			case ZoraContractAddress:
				result.Zora = &UserZoraIdentity{
					Website:    entry.Website,
					Username:   entry.Username,
					DataSource: CONTEXT,
				}
			case ContextContractAddress:
				result.Context = &UserContextIdentity{
					Username:      entry.Username,
					Website:       entry.Website,
					FollowerCount: contextProfile.FollowerCount,
					DataSource:    CONTEXT,
				}
			default:
			}
		}
	}

	controller.AddAndDone(&result)
}

func (f *fetcher) processSuperrare(address string, controller *identityFetcherController) {
	var result IdentityEntry

	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    fmt.Sprintf(SuperrareUrl, address),
		method: "GET",
	})
	if err != nil {
		result.Err = err
		result.Msg = "[processSuperrare] fetch identity failed"
		controller.AddAndDone(&result)
		return
	}

	sprProfile := SuperrareProfile{}
	err = json.Unmarshal(body, &sprProfile)
	if err != nil {
		result.Err = err
		result.Msg = "[processSuperrare] identity response json unmarshal failednti"
		controller.AddAndDone(&result)
		return
	}

	newSprRecord := UserSuperrareIdentity{
		Username:       sprProfile.Result.Username,
		Location:       sprProfile.Result.Location,
		Bio:            sprProfile.Result.Bio,
		InstagramLink:  sprProfile.Result.InstagramLink,
		TwitterLink:    sprProfile.Result.TwitterLink,
		SteemitLink:    sprProfile.Result.SteemitLink,
		Website:        sprProfile.Result.Website,
		SpotifyLink:    sprProfile.Result.SpotifyLink,
		SoundCloudLink: sprProfile.Result.SoundCloudLink,
		DataSource:     SUPERRARE,
	}

	if newSprRecord.Username != "" || newSprRecord.Location != "" || newSprRecord.Bio != "" || newSprRecord.InstagramLink != "" ||
		newSprRecord.TwitterLink != "" || newSprRecord.SteemitLink != "" || newSprRecord.Website != "" ||
		newSprRecord.SpotifyLink != "" || newSprRecord.SoundCloudLink != "" {
		result.Superrare = &newSprRecord
	}

	controller.AddAndDone(&result)
}

func (f *fetcher) processSybil(address string, controller *identityFetcherController) {

	var result IdentityEntry

	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    SybilUrl,
		method: "GET",
	})

	if err != nil {
		result.Err = err
		result.Msg = "[processSybil] fetch identity failed"
		controller.AddAndDone(&result)
		return
	}

	var sybilVerifiedProfiles = SybilVerifiedList{}
	err = json.Unmarshal(body, &sybilVerifiedProfiles)
	if err != nil {
		result.Err = err
		result.Msg = "[processSybil] unmarshal json failed"
		controller.AddAndDone(&result)
		return
	}

	entry, exists := sybilVerifiedProfiles[address]

	if exists {
		twitterIdentity := UserTwitterIdentity{entry.Twitter.Handle, SYBIL}
		result.Twitter = &twitterIdentity
	}

	controller.AddAndDone(&result)
}

func (f *fetcher) processPoap(address string, controller *identityFetcherController) {
	var result IdentityEntry
	// call poap actions/scan API
	body, err := sendRequest(f.httpClient, RequestArgs{
		url:    fmt.Sprintf(PoapTokensFetchUrl, address),
		method: "GET",
	})
	if err != nil {
		result.Err = err
		result.Msg = "[processPoap] request failed"
		controller.AddAndDone(&result)
		return
	}

	// json decode and convert all identities to IdentityEntry
	var poapResults []PoapActionScanResultEntry
	err = json.Unmarshal(body, &poapResults)
	if err != nil {
		result.Err = err
		result.Msg = "[processPoap] unmarshal json failed"
		controller.AddAndDone(&result)
		return
	}

	// send all IdentityEntry-ies to channel
	g := new(errgroup.Group)
	for _, poapResult := range poapResults {
		eventId := fmt.Sprintf("%d", poapResult.Event.Id)
		tokenId := poapResult.TokenID
		g.Go(func() error {
			recommendations := f.processPoapRecommendations(eventId, tokenId)
			poapIdentity := UserPoapIdentity{
				EventID:         eventId,
				EventDesc:       poapResult.Event.Description,
				TokenID:         tokenId,
				Recommendations: recommendations,
			}
			result.Poap = &poapIdentity
			controller.Add(&result)
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		result.Err = err
		result.Msg = "[processPoap] some recommendations failed!"
		controller.AddAndDone(&result)
	}

	controller.Done()
}

func (f *fetcher) processPoapRecommendations(eventId string, userPoapTokenId string) []PoapRecommendation {
	var result []PoapRecommendation
	graphQuery := map[string]string{
		"query": fmt.Sprintf(`
		{
		  event(id: "%s") {
			tokens {
			  id
			  owner {
				id
			  }
			}
		  }
		}		  
	  `, eventId),
	}
	requestBody, err := json.Marshal(graphQuery)
	if err != nil {
		// ideally should not occur
		fmt.Printf("Invalid request body for event %s\n", eventId)
		return result
	}
	responseBody, err := sendRequest(f.httpClient, RequestArgs{
		url:    PoapXyzGraphUrl,
		method: "POST",
		body:   requestBody,
	})
	if err != nil {
		fmt.Printf("Unable to locate POAP recommendations for event %s\n", eventId)
		return result
	}

	var attendees PoapTokenAndOwnerGraphQueryResult
	err = json.Unmarshal(responseBody, &attendees)
	if err != nil {
		fmt.Printf("Unable to JSON decode POAP recommendations for event %s\n", eventId)
		return result
	}
	for _, attendeeToken := range attendees.Data.Event.Tokens {
		if userPoapTokenId == attendeeToken.Id {
			// skip self
			continue
		}
		recommendation := PoapRecommendation{
			EventID: eventId,
			TokenID: attendeeToken.Id,
			Address: attendeeToken.Owner.Id,
		}
		result = append(result, recommendation)
	}
	return result
}
