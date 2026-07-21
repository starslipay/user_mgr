package logic

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenMD5(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "D41D8CD98F00B204E9800998ECF8427E",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "5EB63BBBE01EEED093CB22BB8F5ACDC3",
		},
		{
			name:     "test 123",
			input:    "test123",
			expected: "CC03E747A6AFBBCBF8BE7668ACFEBEE5",
		},
		{
			name:     "pay100011718901234",
			input:    "pay100011718901234",
			expected: "FB6E32932BB7AFF0DFF5948914C56A02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenMD5(tc.input)
			assert.Equal(t, tc.expected, result)
			assert.Len(t, result, 32)
			assert.True(t, strings.ToUpper(result) == result, "MD5 should be uppercase")
		})
	}
}

func TestGenUserToken(t *testing.T) {
	testCases := []struct {
		name         string
		userID       string
		businessInfo string
	}{
		{
			name:         "normal case",
			userID:       "10001",
			businessInfo: "pay",
		},
		{
			name:         "empty business info",
			userID:       "10002",
			businessInfo: "",
		},
		{
			name:         "special characters",
			userID:       "user@test",
			businessInfo: "trade",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token := GenUserToken(tc.userID, tc.businessInfo)

			assert.Len(t, token, 42, "Token should be 42 characters (32 MD5 + 10 timestamp)")

			md5Str := token[:32]
			timestampStr := token[32:]

			assert.Len(t, md5Str, 32, "MD5 part should be 32 characters")
			assert.Len(t, timestampStr, 10, "Timestamp part should be 10 characters")

			assert.True(t, strings.ToUpper(md5Str) == md5Str, "MD5 part should be uppercase")

			expectedMD5 := GenMD5(tc.businessInfo + tc.userID + timestampStr)
			assert.Equal(t, expectedMD5, md5Str, "MD5 part should match expected")
		})
	}
}

func TestCheckUserToken(t *testing.T) {
	userID := "10001"
	businessInfo := "pay"

	token := GenUserToken(userID, businessInfo)

	validCases := []struct {
		name         string
		token        string
		userID       string
		businessInfo string
		expected     bool
	}{
		{
			name:         "valid token",
			token:        token,
			userID:       userID,
			businessInfo: businessInfo,
			expected:     true,
		},
		{
			name:         "invalid user id",
			token:        token,
			userID:       "10002",
			businessInfo: businessInfo,
			expected:     false,
		},
		{
			name:         "invalid business info",
			token:        token,
			userID:       userID,
			businessInfo: "trade",
			expected:     false,
		},
		{
			name:         "invalid token format",
			token:        "invalid_token",
			userID:       userID,
			businessInfo: businessInfo,
			expected:     false,
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "invalid token format" {
				assert.Panics(t, func() {
					CheckUserToken(tc.token, tc.userID, tc.businessInfo)
				}, "Should panic on invalid token format")
			} else {
				result := CheckUserToken(tc.token, tc.userID, tc.businessInfo)
				assert.Equal(t, tc.expected, result)
			}
		})
	}

	t.Run("expired token", func(t *testing.T) {
		pastTime := time.Now().Unix() - 120
		timestampStr := strconv.FormatInt(pastTime, 10)
		expiredToken := GenMD5(businessInfo+userID+timestampStr) + timestampStr

		result := CheckUserToken(expiredToken, userID, businessInfo)
		assert.False(t, result, "Expired token should be invalid")
	})
}
