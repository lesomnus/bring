package cmd_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lesomnus/bring/cmd"
	"github.com/stretchr/testify/require"
)

func TestCmdVersion(t *testing.T) {
	require := require.New(t)

	b := &bytes.Buffer{}
	c := cmd.NewCmdVersion()
	c.Writer = b
	err := c.Run(context.Background(), []string{})
	require.NoError(err)

	v, err := godotenv.Parse(b)
	require.NoError(err)

	require.Contains(v, "BRING_VERSION")
	require.Contains(v, "BRING_TIME_BUILD")
	require.Contains(v, "BRING_GIT_REV")

	_, err = time.Parse(time.RFC3339, v["BRING_TIME_BUILD"])
	require.NoError(err)
}
