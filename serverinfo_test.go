package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/stretchr/testify/assert"
)

func TestProcessPayload(t *testing.T) {
	t.Parallel()
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	// Create test JSON file.
	testContent := `{"content": "test content"}`
	testFilePath := filepath.Join(path, "test", "1.0.0.json")
	err = os.MkdirAll(filepath.Dir(testFilePath), 0755)
	if err != nil {
		t.Fatalf("Failed to create directory for test file: %v", err)
	}
	err = os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	if _, err := os.Stat(testFilePath); err != nil {
		fmt.Errorf("test file not found: %s", err)
	}
	// Calculate content hash.
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(`{"content": "test content"}`)))
	testPayload := fmt.Sprintf(`{"type": "test", "version": "1.0.0", "hash": "%s"}`, hash)
	tests := []struct {
		name          string
		payload       string
		expectedError bool
	}{
		{"ValidPayload", testPayload, false},                                                       // just typical behaviour
		{"InvalidPayload", `invalid_json`, true},                                                   // invalid json should return an error
		{"DifferentHash", `{"type": "test", "version": "1.0.0", "hash": "different_hash"}`, false}, // test case for different hash
		{"EmptyPayload", `{"type": "", "version": "", "hash": ""}`, false},                         // test case for empty payload
		{"NotExistingFile", `{"type": "not_exist", "version": "9.9.9", "hash": ""}`, true},         // test case for not existing file
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println("-- Start test:", tt.name, "--")
			logger := &testLogger{}
			db := &sql.DB{}           // Mocked DB
			nk := &testNakamaModule{} // Mocked nakama

			responseJSON, err := VersionChecker(context.Background(), logger, db, nk, tt.payload)
			if tt.expectedError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Expected no error")
				assert.NotNil(t, responseJSON, "Expected response JSON")

				// Parse response to check if default values are used when payload fields are empty.
				var response Response
				err := json.Unmarshal([]byte(responseJSON), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if tt.name == "ValidPayload" {
					assert.Equal(t, "test", response.Type, "Expected default value for Type")
					assert.Equal(t, "1.0.0", response.Version, "Expected default value for Version")
					assert.Equal(t, `{"content": "test content"}`, response.Content, "Expected default value for content")
					assert.Equal(t, hash, response.Hash, "Expected default value for Hash")
				}
				if tt.name == "DifferentHash" {
					assert.Empty(t, response.Content, "Expected content to be empty for different hash")
				}
				if tt.name == "EmptyPayload" {
					assert.Equal(t, "core", response.Type, "Expected default value for Type")
					assert.Equal(t, "1.0.0", response.Version, "Expected default value for Version")
				}
			}
		})
	}
}

type testLogger struct{}

// Fields implements runtime.Logger.
func (l *testLogger) Fields() map[string]interface{} {
	panic("unimplemented")
}

// Warn implements runtime.Logger.
func (l *testLogger) Warn(format string, v ...interface{}) {
	panic("unimplemented")
}

// WithField implements runtime.Logger.
func (l *testLogger) WithField(key string, v interface{}) runtime.Logger {
	panic("unimplemented")
}

// WithFields implements runtime.Logger.
func (l *testLogger) WithFields(fields map[string]interface{}) runtime.Logger {
	panic("unimplemented")
}

func (l *testLogger) Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (l *testLogger) Error(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (l *testLogger) Debug(format string, args ...interface{}) {
	// Minimal implementation of Debug for testing.
	fmt.Printf(format+"\n", args...)
}

type testNakamaModule struct{}

// AccountDeleteId implements runtime.NakamaModule.
func (t *testNakamaModule) AccountDeleteId(ctx context.Context, userID string, recorded bool) error {
	panic("unimplemented")
}

// AccountExportId implements runtime.NakamaModule.
func (t *testNakamaModule) AccountExportId(ctx context.Context, userID string) (string, error) {
	panic("unimplemented")
}

// AccountGetId implements runtime.NakamaModule.
func (t *testNakamaModule) AccountGetId(ctx context.Context, userID string) (*api.Account, error) {
	panic("unimplemented")
}

// AccountUpdateId implements runtime.NakamaModule.
func (t *testNakamaModule) AccountUpdateId(ctx context.Context, userID string, username string, metadata map[string]interface{}, displayName string, timezone string, location string, langTag string, avatarUrl string) error {
	panic("unimplemented")
}

// AccountsGetId implements runtime.NakamaModule.
func (t *testNakamaModule) AccountsGetId(ctx context.Context, userIDs []string) ([]*api.Account, error) {
	panic("unimplemented")
}

// AuthenticateApple implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateApple(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateCustom implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateCustom(ctx context.Context, id string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateDevice implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateDevice(ctx context.Context, id string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateEmail implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateEmail(ctx context.Context, email string, password string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateFacebook implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateFacebook(ctx context.Context, token string, importFriends bool, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateFacebookInstantGame implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateFacebookInstantGame(ctx context.Context, signedPlayerInfo string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateGameCenter implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateGameCenter(ctx context.Context, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateGoogle implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateGoogle(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateSteam implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateSteam(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateTokenGenerate implements runtime.NakamaModule.
func (t *testNakamaModule) AuthenticateTokenGenerate(userID string, username string, exp int64, vars map[string]string) (string, int64, error) {
	panic("unimplemented")
}

// ChannelIdBuild implements runtime.NakamaModule.
func (t *testNakamaModule) ChannelIdBuild(ctx context.Context, sender string, target string, chanType runtime.ChannelType) (string, error) {
	panic("unimplemented")
}

// ChannelMessageRemove implements runtime.NakamaModule.
func (t *testNakamaModule) ChannelMessageRemove(ctx context.Context, channelId string, messageId string, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessageSend implements runtime.NakamaModule.
func (t *testNakamaModule) ChannelMessageSend(ctx context.Context, channelID string, content map[string]interface{}, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessageUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) ChannelMessageUpdate(ctx context.Context, channelID string, messageID string, content map[string]interface{}, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessagesList implements runtime.NakamaModule.
func (t *testNakamaModule) ChannelMessagesList(ctx context.Context, channelId string, limit int, forward bool, cursor string) (messages []*api.ChannelMessage, nextCursor string, prevCursor string, err error) {
	panic("unimplemented")
}

// Event implements runtime.NakamaModule.
func (t *testNakamaModule) Event(ctx context.Context, evt *api.Event) error {
	panic("unimplemented")
}

// FriendsAdd implements runtime.NakamaModule.
func (t *testNakamaModule) FriendsAdd(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsBlock implements runtime.NakamaModule.
func (t *testNakamaModule) FriendsBlock(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsDelete implements runtime.NakamaModule.
func (t *testNakamaModule) FriendsDelete(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsList implements runtime.NakamaModule.
func (t *testNakamaModule) FriendsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.Friend, string, error) {
	panic("unimplemented")
}

// GetSatori implements runtime.NakamaModule.
func (t *testNakamaModule) GetSatori() runtime.Satori {
	panic("unimplemented")
}

// GroupCreate implements runtime.NakamaModule.
func (t *testNakamaModule) GroupCreate(ctx context.Context, userID string, name string, creatorID string, langTag string, description string, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) (*api.Group, error) {
	panic("unimplemented")
}

// GroupDelete implements runtime.NakamaModule.
func (t *testNakamaModule) GroupDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GroupUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUpdate(ctx context.Context, id string, userID string, name string, creatorID string, langTag string, description string, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) error {
	panic("unimplemented")
}

// GroupUserJoin implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUserJoin(ctx context.Context, groupID string, userID string, username string) error {
	panic("unimplemented")
}

// GroupUserLeave implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUserLeave(ctx context.Context, groupID string, userID string, username string) error {
	panic("unimplemented")
}

// GroupUsersAdd implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersAdd(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersBan implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersBan(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersDemote implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersDemote(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersKick implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersKick(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersList implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersList(ctx context.Context, id string, limit int, state *int, cursor string) ([]*api.GroupUserList_GroupUser, string, error) {
	panic("unimplemented")
}

// GroupUsersPromote implements runtime.NakamaModule.
func (t *testNakamaModule) GroupUsersPromote(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupsGetId implements runtime.NakamaModule.
func (t *testNakamaModule) GroupsGetId(ctx context.Context, groupIDs []string) ([]*api.Group, error) {
	panic("unimplemented")
}

// GroupsGetRandom implements runtime.NakamaModule.
func (t *testNakamaModule) GroupsGetRandom(ctx context.Context, count int) ([]*api.Group, error) {
	panic("unimplemented")
}

// GroupsList implements runtime.NakamaModule.
func (t *testNakamaModule) GroupsList(ctx context.Context, name string, langTag string, members *int, open *bool, limit int, cursor string) ([]*api.Group, string, error) {
	panic("unimplemented")
}

// LeaderboardCreate implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardCreate(ctx context.Context, id string, authoritative bool, sortOrder string, operator string, resetSchedule string, metadata map[string]interface{}) error {
	panic("unimplemented")
}

// LeaderboardDelete implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// LeaderboardList implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardList(categoryStart int, categoryEnd int, limit int, cursor string) (*api.LeaderboardList, error) {
	panic("unimplemented")
}

// LeaderboardRecordDelete implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardRecordDelete(ctx context.Context, id string, ownerID string) error {
	panic("unimplemented")
}

// LeaderboardRecordWrite implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardRecordWrite(ctx context.Context, id string, ownerID string, username string, score int64, subscore int64, metadata map[string]interface{}, overrideOperator *int) (*api.LeaderboardRecord, error) {
	panic("unimplemented")
}

// LeaderboardRecordsHaystack implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardRecordsHaystack(ctx context.Context, id string, ownerID string, limit int, cursor string, expiry int64) (*api.LeaderboardRecordList, error) {
	panic("unimplemented")
}

// LeaderboardRecordsList implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardRecordsList(ctx context.Context, id string, ownerIDs []string, limit int, cursor string, expiry int64) (records []*api.LeaderboardRecord, ownerRecords []*api.LeaderboardRecord, nextCursor string, prevCursor string, err error) {
	panic("unimplemented")
}

// LeaderboardsGetId implements runtime.NakamaModule.
func (t *testNakamaModule) LeaderboardsGetId(ctx context.Context, ids []string) ([]*api.Leaderboard, error) {
	panic("unimplemented")
}

// LinkApple implements runtime.NakamaModule.
func (t *testNakamaModule) LinkApple(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// LinkCustom implements runtime.NakamaModule.
func (t *testNakamaModule) LinkCustom(ctx context.Context, userID string, customID string) error {
	panic("unimplemented")
}

// LinkDevice implements runtime.NakamaModule.
func (t *testNakamaModule) LinkDevice(ctx context.Context, userID string, deviceID string) error {
	panic("unimplemented")
}

// LinkEmail implements runtime.NakamaModule.
func (t *testNakamaModule) LinkEmail(ctx context.Context, userID string, email string, password string) error {
	panic("unimplemented")
}

// LinkFacebook implements runtime.NakamaModule.
func (t *testNakamaModule) LinkFacebook(ctx context.Context, userID string, username string, token string, importFriends bool) error {
	panic("unimplemented")
}

// LinkFacebookInstantGame implements runtime.NakamaModule.
func (t *testNakamaModule) LinkFacebookInstantGame(ctx context.Context, userID string, signedPlayerInfo string) error {
	panic("unimplemented")
}

// LinkGameCenter implements runtime.NakamaModule.
func (t *testNakamaModule) LinkGameCenter(ctx context.Context, userID string, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string) error {
	panic("unimplemented")
}

// LinkGoogle implements runtime.NakamaModule.
func (t *testNakamaModule) LinkGoogle(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// LinkSteam implements runtime.NakamaModule.
func (t *testNakamaModule) LinkSteam(ctx context.Context, userID string, username string, token string, importFriends bool) error {
	panic("unimplemented")
}

// MatchCreate implements runtime.NakamaModule.
func (t *testNakamaModule) MatchCreate(ctx context.Context, module string, params map[string]interface{}) (string, error) {
	panic("unimplemented")
}

// MatchGet implements runtime.NakamaModule.
func (t *testNakamaModule) MatchGet(ctx context.Context, id string) (*api.Match, error) {
	panic("unimplemented")
}

// MatchList implements runtime.NakamaModule.
func (t *testNakamaModule) MatchList(ctx context.Context, limit int, authoritative bool, label string, minSize *int, maxSize *int, query string) ([]*api.Match, error) {
	panic("unimplemented")
}

// MatchSignal implements runtime.NakamaModule.
func (t *testNakamaModule) MatchSignal(ctx context.Context, id string, data string) (string, error) {
	panic("unimplemented")
}

// MetricsCounterAdd implements runtime.NakamaModule.
func (t *testNakamaModule) MetricsCounterAdd(name string, tags map[string]string, delta int64) {
	panic("unimplemented")
}

// MetricsGaugeSet implements runtime.NakamaModule.
func (t *testNakamaModule) MetricsGaugeSet(name string, tags map[string]string, value float64) {
	panic("unimplemented")
}

// MetricsTimerRecord implements runtime.NakamaModule.
func (t *testNakamaModule) MetricsTimerRecord(name string, tags map[string]string, value time.Duration) {
	panic("unimplemented")
}

// MultiUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) MultiUpdate(ctx context.Context, accountUpdates []*runtime.AccountUpdate, storageWrites []*runtime.StorageWrite, walletUpdates []*runtime.WalletUpdate, updateLedger bool) ([]*api.StorageObjectAck, []*runtime.WalletUpdateResult, error) {
	panic("unimplemented")
}

// NotificationSend implements runtime.NakamaModule.
func (t *testNakamaModule) NotificationSend(ctx context.Context, userID string, subject string, content map[string]interface{}, code int, sender string, persistent bool) error {
	panic("unimplemented")
}

// NotificationSendAll implements runtime.NakamaModule.
func (t *testNakamaModule) NotificationSendAll(ctx context.Context, subject string, content map[string]interface{}, code int, persistent bool) error {
	panic("unimplemented")
}

// NotificationsDelete implements runtime.NakamaModule.
func (t *testNakamaModule) NotificationsDelete(ctx context.Context, notifications []*runtime.NotificationDelete) error {
	panic("unimplemented")
}

// NotificationsSend implements runtime.NakamaModule.
func (t *testNakamaModule) NotificationsSend(ctx context.Context, notifications []*runtime.NotificationSend) error {
	panic("unimplemented")
}

// PurchaseGetByTransactionId implements runtime.NakamaModule.
func (t *testNakamaModule) PurchaseGetByTransactionId(ctx context.Context, transactionID string) (*api.ValidatedPurchase, error) {
	panic("unimplemented")
}

// PurchaseValidateApple implements runtime.NakamaModule.
func (t *testNakamaModule) PurchaseValidateApple(ctx context.Context, userID string, receipt string, persist bool, passwordOverride ...string) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchaseValidateGoogle implements runtime.NakamaModule.
func (t *testNakamaModule) PurchaseValidateGoogle(ctx context.Context, userID string, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchaseValidateHuawei implements runtime.NakamaModule.
func (t *testNakamaModule) PurchaseValidateHuawei(ctx context.Context, userID string, signature string, inAppPurchaseData string, persist bool) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchasesList implements runtime.NakamaModule.
func (t *testNakamaModule) PurchasesList(ctx context.Context, userID string, limit int, cursor string) (*api.PurchaseList, error) {
	panic("unimplemented")
}

// ReadFile implements runtime.NakamaModule.
func (t *testNakamaModule) ReadFile(path string) (*os.File, error) {
	panic("unimplemented")
}

// SessionDisconnect implements runtime.NakamaModule.
func (t *testNakamaModule) SessionDisconnect(ctx context.Context, sessionID string, reason ...runtime.PresenceReason) error {
	panic("unimplemented")
}

// SessionLogout implements runtime.NakamaModule.
func (t *testNakamaModule) SessionLogout(userID string, token string, refreshToken string) error {
	panic("unimplemented")
}

// StorageDelete implements runtime.NakamaModule.
func (t *testNakamaModule) StorageDelete(ctx context.Context, deletes []*runtime.StorageDelete) error {
	panic("unimplemented")
}

// StorageList implements runtime.NakamaModule.
func (t *testNakamaModule) StorageList(ctx context.Context, userID string, collection string, limit int, cursor string) ([]*api.StorageObject, string, error) {
	panic("unimplemented")
}

// StorageRead implements runtime.NakamaModule.
func (t *testNakamaModule) StorageRead(ctx context.Context, reads []*runtime.StorageRead) ([]*api.StorageObject, error) {
	panic("unimplemented")
}

// StorageWrite implements runtime.NakamaModule.
func (t *testNakamaModule) StorageWrite(ctx context.Context, writes []*runtime.StorageWrite) ([]*api.StorageObjectAck, error) {
	return []*api.StorageObjectAck{}, nil
}

// StreamClose implements runtime.NakamaModule.
func (t *testNakamaModule) StreamClose(mode uint8, subject string, subcontext string, label string) error {
	panic("unimplemented")
}

// StreamCount implements runtime.NakamaModule.
func (t *testNakamaModule) StreamCount(mode uint8, subject string, subcontext string, label string) (int, error) {
	panic("unimplemented")
}

// StreamSend implements runtime.NakamaModule.
func (t *testNakamaModule) StreamSend(mode uint8, subject string, subcontext string, label string, data string, presences []runtime.Presence, reliable bool) error {
	panic("unimplemented")
}

// StreamSendRaw implements runtime.NakamaModule.
func (t *testNakamaModule) StreamSendRaw(mode uint8, subject string, subcontext string, label string, msg *rtapi.Envelope, presences []runtime.Presence, reliable bool) error {
	panic("unimplemented")
}

// StreamUserGet implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserGet(mode uint8, subject string, subcontext string, label string, userID string, sessionID string) (runtime.PresenceMeta, error) {
	panic("unimplemented")
}

// StreamUserJoin implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserJoin(mode uint8, subject string, subcontext string, label string, userID string, sessionID string, hidden bool, persistence bool, status string) (bool, error) {
	panic("unimplemented")
}

// StreamUserKick implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserKick(mode uint8, subject string, subcontext string, label string, presence runtime.Presence) error {
	panic("unimplemented")
}

// StreamUserLeave implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserLeave(mode uint8, subject string, subcontext string, label string, userID string, sessionID string) error {
	panic("unimplemented")
}

// StreamUserList implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserList(mode uint8, subject string, subcontext string, label string, includeHidden bool, includeNotHidden bool) ([]runtime.Presence, error) {
	panic("unimplemented")
}

// StreamUserUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) StreamUserUpdate(mode uint8, subject string, subcontext string, label string, userID string, sessionID string, hidden bool, persistence bool, status string) error {
	panic("unimplemented")
}

// SubscriptionGetByProductId implements runtime.NakamaModule.
func (t *testNakamaModule) SubscriptionGetByProductId(ctx context.Context, userID string, productID string) (*api.ValidatedSubscription, error) {
	panic("unimplemented")
}

// SubscriptionValidateApple implements runtime.NakamaModule.
func (t *testNakamaModule) SubscriptionValidateApple(ctx context.Context, userID string, receipt string, persist bool, passwordOverride ...string) (*api.ValidateSubscriptionResponse, error) {
	panic("unimplemented")
}

// SubscriptionValidateGoogle implements runtime.NakamaModule.
func (t *testNakamaModule) SubscriptionValidateGoogle(ctx context.Context, userID string, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidateSubscriptionResponse, error) {
	panic("unimplemented")
}

// SubscriptionsList implements runtime.NakamaModule.
func (t *testNakamaModule) SubscriptionsList(ctx context.Context, userID string, limit int, cursor string) (*api.SubscriptionList, error) {
	panic("unimplemented")
}

// TournamentAddAttempt implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentAddAttempt(ctx context.Context, id string, ownerID string, count int) error {
	panic("unimplemented")
}

// TournamentCreate implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentCreate(ctx context.Context, id string, authoritative bool, sortOrder string, operator string, resetSchedule string, metadata map[string]interface{}, title string, description string, category int, startTime int, endTime int, duration int, maxSize int, maxNumScore int, joinRequired bool) error {
	panic("unimplemented")
}

// TournamentDelete implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// TournamentJoin implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentJoin(ctx context.Context, id string, ownerID string, username string) error {
	panic("unimplemented")
}

// TournamentList implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentList(ctx context.Context, categoryStart int, categoryEnd int, startTime int, endTime int, limit int, cursor string) (*api.TournamentList, error) {
	panic("unimplemented")
}

// TournamentRecordDelete implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentRecordDelete(ctx context.Context, id string, ownerID string) error {
	panic("unimplemented")
}

// TournamentRecordWrite implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentRecordWrite(ctx context.Context, id string, ownerID string, username string, score int64, subscore int64, metadata map[string]interface{}, operatorOverride *int) (*api.LeaderboardRecord, error) {
	panic("unimplemented")
}

// TournamentRecordsHaystack implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentRecordsHaystack(ctx context.Context, id string, ownerID string, limit int, cursor string, expiry int64) (*api.TournamentRecordList, error) {
	panic("unimplemented")
}

// TournamentRecordsList implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentRecordsList(ctx context.Context, tournamentId string, ownerIDs []string, limit int, cursor string, overrideExpiry int64) (records []*api.LeaderboardRecord, ownerRecords []*api.LeaderboardRecord, prevCursor string, nextCursor string, err error) {
	panic("unimplemented")
}

// TournamentsGetId implements runtime.NakamaModule.
func (t *testNakamaModule) TournamentsGetId(ctx context.Context, tournamentIDs []string) ([]*api.Tournament, error) {
	panic("unimplemented")
}

// UnlinkApple implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkApple(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkCustom implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkCustom(ctx context.Context, userID string, customID string) error {
	panic("unimplemented")
}

// UnlinkDevice implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkDevice(ctx context.Context, userID string, deviceID string) error {
	panic("unimplemented")
}

// UnlinkEmail implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkEmail(ctx context.Context, userID string, email string) error {
	panic("unimplemented")
}

// UnlinkFacebook implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkFacebook(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkFacebookInstantGame implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkFacebookInstantGame(ctx context.Context, userID string, signedPlayerInfo string) error {
	panic("unimplemented")
}

// UnlinkGameCenter implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkGameCenter(ctx context.Context, userID string, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string) error {
	panic("unimplemented")
}

// UnlinkGoogle implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkGoogle(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkSteam implements runtime.NakamaModule.
func (t *testNakamaModule) UnlinkSteam(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UserGroupsList implements runtime.NakamaModule.
func (t *testNakamaModule) UserGroupsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.UserGroupList_UserGroup, string, error) {
	panic("unimplemented")
}

// UsersBanId implements runtime.NakamaModule.
func (t *testNakamaModule) UsersBanId(ctx context.Context, userIDs []string) error {
	panic("unimplemented")
}

// UsersGetId implements runtime.NakamaModule.
func (t *testNakamaModule) UsersGetId(ctx context.Context, userIDs []string, facebookIDs []string) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersGetRandom implements runtime.NakamaModule.
func (t *testNakamaModule) UsersGetRandom(ctx context.Context, count int) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersGetUsername implements runtime.NakamaModule.
func (t *testNakamaModule) UsersGetUsername(ctx context.Context, usernames []string) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersUnbanId implements runtime.NakamaModule.
func (t *testNakamaModule) UsersUnbanId(ctx context.Context, userIDs []string) error {
	panic("unimplemented")
}

// WalletLedgerList implements runtime.NakamaModule.
func (t *testNakamaModule) WalletLedgerList(ctx context.Context, userID string, limit int, cursor string) ([]runtime.WalletLedgerItem, string, error) {
	panic("unimplemented")
}

// WalletLedgerUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) WalletLedgerUpdate(ctx context.Context, itemID string, metadata map[string]interface{}) (runtime.WalletLedgerItem, error) {
	panic("unimplemented")
}

// WalletUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) WalletUpdate(ctx context.Context, userID string, changeset map[string]int64, metadata map[string]interface{}, updateLedger bool) (updated map[string]int64, previous map[string]int64, err error) {
	panic("unimplemented")
}

// WalletsUpdate implements runtime.NakamaModule.
func (t *testNakamaModule) WalletsUpdate(ctx context.Context, updates []*runtime.WalletUpdate, updateLedger bool) ([]*runtime.WalletUpdateResult, error) {
	panic("unimplemented")
}
