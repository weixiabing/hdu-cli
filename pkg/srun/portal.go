package srun

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hduhelp/hdu-cli/internal/core"
)

func New(endpoint, acID string) *PortalServer {
	timestampStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	return &PortalServer{
		endPoint:      endpoint,
		acID:          acID,
		jsonpCallback: "jQuery112403771213770126085_" + timestampStr,
		timestampStr:  timestampStr,
		internetCheck: "http://www.baidu.com",
	}
}

func (s *PortalServer) SetUsername(username string) error {
	if username == "" {
		return errors.New("username is empty")
	}
	s.username = username
	return nil
}

func (s *PortalServer) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is empty")
	}
	s.password = password
	return nil
}

func (s *PortalServer) SetInternetCheckEndpoint(uri string) error {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return err
	}
	s.internetCheck = uri
	return nil
}

type PortalServer struct {
	endPoint string
	// AcID NasID?
	acID          string
	jsonpCallback string

	internetCheck string

	username string
	password string

	timestampStr string

	userInfo       *userInfo
	challenge      *challenge
	loginResponse  *loginResponse
	logoutResponse *logoutResponse

	adapterHooks portalAdapterHooks
}

type portalAdapterHooks struct {
	getUserInfo  func() (*userInfo, error)
	getChallenge func() (*challenge, error)
	portalLogin  func() (*loginResponse, error)
	portalLogout func() (*logoutResponse, error)
}

func (s *PortalServer) Login(ctx context.Context, username, password string) (core.Status, error) {
	_ = ctx

	if err := s.SetUsername(username); err != nil {
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, err
	}
	if err := s.SetPassword(password); err != nil {
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, err
	}
	info, err := s.getUserInfo()
	if err != nil {
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, err
	}
	if _, err := s.getChallenge(); err != nil {
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, err
	}
	resp, err := s.portalLogin()
	if err != nil {
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, mapLoginError(err)
	}
	if ok, err := resp.IsOK(); !ok {
		if err == nil {
			err = core.ErrAuthenticationFailed
		}
		return core.Status{Phase: core.PhaseFailed, Online: false, Message: err.Error()}, mapLoginError(err)
	}

	s.userInfo = info
	return core.Status{Phase: core.PhaseConnected, Online: true}, nil
}

func (s *PortalServer) Logout(ctx context.Context, username string) error {
	_ = ctx
	if username != "" {
		if err := s.SetUsername(username); err != nil {
			return err
		}
	}
	info, err := s.getUserInfo()
	if err != nil {
		return err
	}
	if isOfflineUserInfo(info) {
		return nil
	}
	resp, err := s.portalLogout()
	if err != nil {
		return err
	}
	if ok, err := resp.IsOK(); !ok {
		return err
	}
	return nil
}

func (s *PortalServer) CurrentStatus(ctx context.Context) (core.Status, error) {
	_ = ctx
	info, err := s.getUserInfo()
	if err != nil {
		return core.Status{Phase: core.PhaseDisconnected, Online: false, Message: err.Error()}, err
	}
	if isOfflineUserInfo(info) {
		return core.Status{Phase: core.PhaseDisconnected, Online: false}, nil
	}
	if ok, err := info.IsOK(); !ok {
		if err == nil {
			err = errors.New("portal status unavailable")
		}
		return core.Status{Phase: core.PhaseDisconnected, Online: false, Message: err.Error()}, err
	}
	return core.Status{Phase: core.PhaseConnected, Online: true}, nil
}

func (s *PortalServer) InternetReachable(ctx context.Context) bool {
	_ = ctx
	return s.Internet()
}

func (s PortalServer) callback() string {
	return s.jsonpCallback
}

func (s *PortalServer) SetAcID(acID string) {
	s.acID = acID
}

func (s PortalServer) AcID() string {
	return s.acID
}

func (s PortalServer) apiUri(path string) *url.URL {
	uri, err := url.ParseRequestURI(s.endPoint + path)
	if err != nil {
		panic("endpoint uri error")
	}
	return uri
}

type ResponseError struct {
	ErrorCode interface{} `json:"ecode" chinese:"错误码"`      //错误码
	Error     string      `json:"error" chinese:"错误信息"`     //错误信息
	ErrorMsg  string      `json:"error_msg" chinese:"错误信息"` //错误信息
}

func (e ResponseError) IsOK() (bool, error) {
	if e.Error != "ok" {
		return false, errors.New(e.ErrorMsg)
	}
	return true, nil
}

func mapLoginError(err error) error {
	if err == nil {
		return nil
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "password") || strings.Contains(message, "auth") {
		return core.ErrAuthenticationFailed
	}
	return err
}

func isOfflineUserInfo(info *userInfo) bool {
	return info != nil && info.Error == "not_online_error"
}

func (s *PortalServer) getUserInfo() (*userInfo, error) {
	if s.adapterHooks.getUserInfo != nil {
		info, err := s.adapterHooks.getUserInfo()
		if info != nil {
			s.userInfo = info
		}
		return info, err
	}
	return s.GetUserInfo()
}

func (s *PortalServer) getChallenge() (*challenge, error) {
	if s.adapterHooks.getChallenge != nil {
		challenge, err := s.adapterHooks.getChallenge()
		if challenge != nil {
			s.challenge = challenge
		}
		return challenge, err
	}
	return s.GetChallenge()
}

func (s *PortalServer) portalLogin() (*loginResponse, error) {
	if s.adapterHooks.portalLogin != nil {
		resp, err := s.adapterHooks.portalLogin()
		if resp != nil {
			s.loginResponse = resp
		}
		return resp, err
	}
	return s.PortalLogin()
}

func (s *PortalServer) portalLogout() (*logoutResponse, error) {
	if s.adapterHooks.portalLogout != nil {
		resp, err := s.adapterHooks.portalLogout()
		if resp != nil {
			s.logoutResponse = resp
		}
		return resp, err
	}
	return s.PortalLogout()
}
