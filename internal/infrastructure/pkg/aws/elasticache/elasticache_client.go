package elasticache

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"

	awsclient "github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/pkg/aws"
)

type ElastiCacheClient struct {
	client *elasticache.Client
}

// NewElastiCacheClient creates a new ElastiCache client using AWS client factory
func NewElastiCacheClient(awsClientFactory *awsclient.ClientFactory) (*ElastiCacheClient, error) {
	client := elasticache.NewFromConfig(awsClientFactory.GetConfig())

	return &ElastiCacheClient{
		client: client,
	}, nil
}

// DescribeCacheClusters describes cache clusters
func (e *ElastiCacheClient) DescribeCacheClusters(ctx context.Context, clusterID *string) (*elasticache.DescribeCacheClustersOutput, error) {
	input := &elasticache.DescribeCacheClustersInput{}

	if clusterID != nil {
		input.CacheClusterId = clusterID
		input.ShowCacheNodeInfo = aws.Bool(true)
	}

	result, err := e.client.DescribeCacheClusters(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe cache clusters: %w", err)
	}

	return result, nil
}

// DescribeReplicationGroups describes replication groups
func (e *ElastiCacheClient) DescribeReplicationGroups(ctx context.Context, replicationGroupID *string) (*elasticache.DescribeReplicationGroupsOutput, error) {
	input := &elasticache.DescribeReplicationGroupsInput{}

	if replicationGroupID != nil {
		input.ReplicationGroupId = replicationGroupID
	}

	result, err := e.client.DescribeReplicationGroups(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe replication groups: %w", err)
	}

	return result, nil
}

// CreateCacheCluster creates a new cache cluster
func (e *ElastiCacheClient) CreateCacheCluster(ctx context.Context, clusterID, nodeType string, numCacheNodes int32, engine string) (*elasticache.CreateCacheClusterOutput, error) {
	input := &elasticache.CreateCacheClusterInput{
		CacheClusterId: aws.String(clusterID),
		CacheNodeType:  aws.String(nodeType),
		NumCacheNodes:  aws.Int32(numCacheNodes),
		Engine:         aws.String(engine),
	}

	result, err := e.client.CreateCacheCluster(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache cluster %s: %w", clusterID, err)
	}

	return result, nil
}

// DeleteCacheCluster deletes a cache cluster
func (e *ElastiCacheClient) DeleteCacheCluster(ctx context.Context, clusterID string) (*elasticache.DeleteCacheClusterOutput, error) {
	input := &elasticache.DeleteCacheClusterInput{
		CacheClusterId: aws.String(clusterID),
	}

	result, err := e.client.DeleteCacheCluster(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete cache cluster %s: %w", clusterID, err)
	}

	return result, nil
}

// CreateReplicationGroup creates a new replication group
func (e *ElastiCacheClient) CreateReplicationGroup(ctx context.Context, replicationGroupID, description, nodeType string, numCacheClusters int32) (*elasticache.CreateReplicationGroupOutput, error) {
	input := &elasticache.CreateReplicationGroupInput{
		ReplicationGroupId:          aws.String(replicationGroupID),
		ReplicationGroupDescription: aws.String(description),
		CacheNodeType:               aws.String(nodeType),
		NumCacheClusters:            aws.Int32(numCacheClusters),
	}

	result, err := e.client.CreateReplicationGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create replication group %s: %w", replicationGroupID, err)
	}

	return result, nil
}

// DeleteReplicationGroup deletes a replication group
func (e *ElastiCacheClient) DeleteReplicationGroup(ctx context.Context, replicationGroupID string) (*elasticache.DeleteReplicationGroupOutput, error) {
	input := &elasticache.DeleteReplicationGroupInput{
		ReplicationGroupId: aws.String(replicationGroupID),
	}

	result, err := e.client.DeleteReplicationGroup(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete replication group %s: %w", replicationGroupID, err)
	}

	return result, nil
}

// WaitForClusterAvailable waits for a cache cluster to become available
func (e *ElastiCacheClient) WaitForClusterAvailable(ctx context.Context, clusterID string, maxWaitTime time.Duration) error {
	waiter := elasticache.NewCacheClusterAvailableWaiter(e.client)

	input := &elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String(clusterID),
	}

	return waiter.Wait(ctx, input, maxWaitTime)
}

// WaitForReplicationGroupAvailable waits for a replication group to become available
func (e *ElastiCacheClient) WaitForReplicationGroupAvailable(ctx context.Context, replicationGroupID string, maxWaitTime time.Duration) error {
	waiter := elasticache.NewReplicationGroupAvailableWaiter(e.client)

	input := &elasticache.DescribeReplicationGroupsInput{
		ReplicationGroupId: aws.String(replicationGroupID),
	}

	return waiter.Wait(ctx, input, maxWaitTime)
}

// GetClient returns the underlying ElastiCache client
func (e *ElastiCacheClient) GetClient() *elasticache.Client {
	return e.client
}
