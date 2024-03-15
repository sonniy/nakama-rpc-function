package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"crypto/sha256"

	"github.com/heroiclabs/nakama-common/runtime"
)

// Payload represents the payload structure.
type Payload struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

// Response represents the response structure.
type Response struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

const collectionName string = "ZeptoLabVersionChecker"

// RPC function to process payload with optional parameters.
func VersionChecker(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	// Parse payload JSON.
	var p Payload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.Error("problem with unmarshal: %s", err)
		return "", err
	}

	// Set default values if not provided.
	if p.Type == "" {
		p.Type = "core"
	}
	if p.Version == "" {
		p.Version = "1.0.0"
	}

	// Construct file path.
	filePath := filepath.Join(p.Type, p.Version+".json")

	// Check if file exists.
	if _, err := os.Stat(filePath); err != nil {
		logger.Error("file not found: %s", err)
		return "", fmt.Errorf("file not found: %s", err)
	}

	// Read file content.
	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("failed to read file: %s", err)
		return "", fmt.Errorf("failed to read file: %s", err)
	}

	// Calculate content hash.
	hash := fmt.Sprintf("%x", sha256.Sum256(content))
	logger.Info("File hash is: %s", hash)

	// Construct response.
	response := Response{
		Type:    p.Type,
		Version: p.Version,
		Hash:    hash,
		Content: string(content),
	}
	logger.Info("responce: %s", response)
	// If hashes are not equal, set content to null.
	// c746686a45ad8d1a06fad5502596466e9de877217a9a32f2253c542a71ee10e2
	if p.Hash == "" || p.Hash != hash {
		response.Content = ""
	}

	// Convert response to JSON string.
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error("failed to marshal response: %s", err)
		return "", fmt.Errorf("failed to marshal response: %s", err)
	}
	responseJSONString := string(responseJSON)
	saveToDB(ctx, logger, nk, p, responseJSONString)

	return string(responseJSON), nil
}

func saveToDB(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, p Payload, response string) {
	userID := "00000000-0000-0000-0000-000000000000"
	key := fmt.Sprintf("%s/%s", p.Type, p.Version)
	objectIDs := []*runtime.StorageWrite{&runtime.StorageWrite{
		Collection: collectionName,
		Key:        key,
		UserID:     userID,
		Value:      string(response),
	},
	}
	_, err := nk.StorageWrite(ctx, objectIDs)
	if err != nil {
		logger.WithField("err", err).Error("Storage write error.")
	} else {
		logger.Info("Write data to storage successfully: [Collection: %s, Key:%s, Value: %s]", collectionName, key, response)
	}
}
