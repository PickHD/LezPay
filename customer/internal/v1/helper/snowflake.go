package helper

import (
	"github.com/bwmarrin/snowflake"
)

// GenerateSnowflakeID will generating unique ID
func GenerateSnowflakeID() (int64, error) {
	// Create a new Node with a Node number of 19
	node, err := snowflake.NewNode(19)
	if err != nil {
		return 0, err
	}

	// Generate a snowflake ID.
	return node.Generate().Int64(), nil
}
