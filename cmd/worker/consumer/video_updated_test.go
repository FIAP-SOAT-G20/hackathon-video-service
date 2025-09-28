package main

import (
	"encoding/json"
	"testing"

	valueobject "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/core/domain/value_object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoUpdated_UnmarshalJSON(t *testing.T) {
	// Test JSON structure provided in the requirement
	jsonStr := `{
  "Type" : "Notification",
  "MessageId" : "2a33bf6b-bb93-57c3-afd0-de38c7b6f234",
  "TopicArn" : "arn:aws:sns:us-east-1:905417995957:video-status-updated",
  "Message" : "{\"video_id\":25,\"user_id\":5,\"status\":\"FINISHED\",\"occurred_at\":\"2025-09-28T18:09:41.110973622Z\"}",
  "Timestamp" : "2025-09-28T18:09:41.115Z",
  "SignatureVersion" : "1",
  "Signature" : "n2sK9472MGBlYH6D58MSJjo64pxWlpevdXgJxqmPLhkKf2Aox+90cADrCmycfQaHpRVqFCJwbMvKl2JSofOBjtpdw33LQzyJi9KsQQ6IbjYiiIsgf2SVTqJZdeZeJbBAZ533iFyfOhK5lVM//nLiRSVrz5zHYHQfmzKLYfY/B6KxvE8S3X5nxYG3sAg7bk3gnp92kpLAVRojwNif+XUnDYrliCyBNmEPQg/z9Y7hR+LT+K5OPiwKjZ/u6wLB7ht0E4c+uRU6+l7WONIAshM95HFh4tpO8g7UuVKYXPQ8C9XnLNsTAxminr8vTnHiD4Mewfh3N9WgA2eAXF/N1bh8Ww==",
  "SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-6209c161c6221fdf56ec1eb5c821d112.pem",
  "UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:905417995957:video-status-updated:c9aa674a-a2e1-4d99-9871-327425c720d1"
}`

	t.Run("should unmarshal SNS notification and extract video data", func(t *testing.T) {
		// First unmarshal the SNS notification structure
		var snsNotification SNSNotification
		err := json.Unmarshal([]byte(jsonStr), &snsNotification)
		require.NoError(t, err)

		// Verify SNS notification fields
		assert.Equal(t, "Notification", snsNotification.Type)
		assert.Equal(t, "2a33bf6b-bb93-57c3-afd0-de38c7b6f234", snsNotification.MessageID)
		assert.Equal(t, "arn:aws:sns:us-east-1:905417995957:video-status-updated", snsNotification.TopicArn)
		assert.NotEmpty(t, snsNotification.Message)

		// Then unmarshal the nested Message field to get the video update data
		var updatedVideo VideoUpdated
		err = json.Unmarshal([]byte(snsNotification.Message), &updatedVideo)
		require.NoError(t, err)

		// Verify video update fields
		assert.Equal(t, uint64(25), updatedVideo.VideoID)
		assert.Equal(t, uint64(5), updatedVideo.UserID)
		assert.Equal(t, valueobject.FINISHED, updatedVideo.Status)
		assert.Equal(t, "2025-09-28T18:09:41.110973622Z", updatedVideo.OccurredAt)
	})

	t.Run("should handle different video status values", func(t *testing.T) {
		testCases := []struct {
			status   string
			expected valueobject.VideoStatus
		}{
			{"CREATED", valueobject.CREATED},
			{"UPLOADED", valueobject.UPLOADED},
			{"PROCESSING", valueobject.PROCESSING},
			{"FINISHED", valueobject.FINISHED},
			{"FAILED", valueobject.FAILED},
		}

		for _, tc := range testCases {
			t.Run(tc.status, func(t *testing.T) {
				// Properly escape the message string for JSON
				messageStr := `{\"video_id\":1,\"user_id\":1,\"status\":\"` + tc.status + `\",\"occurred_at\":\"2025-09-28T18:09:41.110973622Z\"}`
				snsStr := `{
				  "Type" : "Notification",
				  "MessageId" : "test-id",
				  "TopicArn" : "arn:aws:sns:us-east-1:905417995957:video-status-updated",
				  "Message" : "` + messageStr + `",
				  "Timestamp" : "2025-09-28T18:09:41.115Z"
				}`

				var snsNotification SNSNotification
				err := json.Unmarshal([]byte(snsStr), &snsNotification)
				require.NoError(t, err)

				var updatedVideo VideoUpdated
				err = json.Unmarshal([]byte(snsNotification.Message), &updatedVideo)
				require.NoError(t, err)

				assert.Equal(t, tc.expected, updatedVideo.Status)
			})
		}
	})
}
