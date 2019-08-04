package snss

import (
	"encoding/binary"
	"strings"

	"github.com/pzl/chromefile/pickle"
)

type UpdateTabNavigation struct {
	TabID                 int32
	Index                 int32
	URL                   string
	Title                 string
	PageState             []byte
	TransitionType        TransitionType // determined from int32
	TransitionMode        TransitionMode // determined from same int32
	HasPostData           bool
	ReferrerURL           string
	ReferrerPolicy        int32
	OriginalRequestURL    string
	IsOverridingUserAgent bool
	/*
		Timestamp             int64
		SearchTerms           string
		HTTPStatusCode        int32
		ReferrerPolicy2       int32
	*/
}

/*
func (u UpdateTabNavigation) String() string {
	return fmt.Sprintf("Update Tab Navigation: Tab: %d, Index: %d, URL: %s, Title: %s\n", u.TabID, u.Index, u.URL, u.Title)
}
*/

type TransitionType int

const (
	TransitionUserClickLink TransitionType = iota
	TransitionOmnibarURL
	TransitionBookmark
	TransitionSubFrameAuto // eg embedded ad
	TransitionSubFrameManual
	TransitionOmnibarSuggestion
	TransitionStartPage
	TransitionFormSubmission
	TransitionReloaded
	TransitionKeywordSearch
	TransitionKeywordHTTP
)

func (t TransitionType) String() string {
	return []string{"user clicked link", "omnibar URL", "clicked bookmark", "subframe auto navigation", "subframe manual navigation", "omnibar suggestion", "start page", "form submission", "reloaded", "keyword search", "keyword http"}[t]
}

const (
	TransitionBackFwdButton  = 1
	TransitionAddressBar     = 1 << 1
	TransitionHomepage       = 1 << 2
	TransitionBeginNavChain  = 1 << 4
	TransitionLastRedirChain = 1 << 5
	TransitionClientRedir    = 1 << 6 // JS or meta tag
	TransitionServerRedir    = 1 << 7 // HTTP Header
)

type TransitionMode int

func (m TransitionMode) String() string {
	mp := map[int]string{
		TransitionBackFwdButton:  "used browser back or fwd button",
		TransitionAddressBar:     "used address bar",
		TransitionHomepage:       "homepage",
		TransitionBeginNavChain:  "begin nav chain",
		TransitionLastRedirChain: "last redirect chain",
		TransitionClientRedir:    "client-side redirect",
		TransitionServerRedir:    "server-side redirect",
	}

	var s strings.Builder
	for k, v := range mp {
		if int(m)&k != 0 {
			s.WriteString(v + ",")
		}
	}
	return s.String()
}

func NewUpdateTabNavigation(data []byte) (UpdateTabNavigation, error) {
	plen := binary.LittleEndian.Uint32(data)
	pData := data[4 : 4+plen]

	t := UpdateTabNavigation{
		TabID: int32(binary.LittleEndian.Uint32(pData)),
		Index: int32(binary.LittleEndian.Uint32(pData[4:])),
	}
	idx := 8

	n, url, err := pickle.ReadString(pData[idx:])
	if err != nil {
		return t, err
	}
	t.URL = url
	idx += n

	n, title, err := pickle.ReadString16(pData[idx:])
	if err != nil {
		return t, err
	}
	t.Title = title
	idx += n

	n, state, err := pickle.ReadBytes(pData[idx:])
	if err != nil {
		return t, err
	}
	t.PageState = state
	idx += n

	tt := int(binary.LittleEndian.Uint32(pData[idx:]))

	t.TransitionType = TransitionType(tt & 0xff)
	t.TransitionMode = TransitionMode(tt >> 24)

	t.HasPostData = binary.LittleEndian.Uint32(pData[idx+4:]) == 1

	idx += 8
	n, refUrl, err := pickle.ReadString(pData[idx:])
	if err != nil {
		return t, err
	}
	t.ReferrerURL = refUrl
	idx += n

	t.ReferrerPolicy = int32(binary.LittleEndian.Uint32(pData[idx:]))
	idx += 4

	n, reqURL, err := pickle.ReadString(pData[idx:])
	if err != nil {
		return t, err
	}
	t.OriginalRequestURL = reqURL
	idx += n

	t.IsOverridingUserAgent = binary.LittleEndian.Uint32(pData[idx:]) == 1

	return t, nil
	/*

		// remaining length?
		err := discardBytes(f, 4)

		nav := UpdateTabNavigation{}
		err = binary.Read(f, binary.LittleEndian, &nav.UpdateTabNavigationHead)
		if err != nil {
			return err
		}

		url, err := lengthPrefixRead(f)
		if err != nil {
			return err
		}
		nav.URL = string(url)
		err = discardBytes(f, 2) // NULs to mark the end of the string

		title, err := doubleLengthPrefixRead(f)
		if err != nil {
			return err
		}
		nav.Title = string(title)

		//thing, err := stalecucumber.Unpickle(bytes.NewReader(buf))
		//fmt.Printf("what: %v\n", thing)
		fmt.Print(nav)
	*/
}
