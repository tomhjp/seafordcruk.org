package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmailAddress(t *testing.T) {
	emailAddress := getEmailAddress()
	assert.True(t, strings.HasSuffix(emailAddress, "@gmail.com"), "expected valid email address suffix")
	assert.Len(t, emailAddress, 16, "expected email address to be 16 characters long")
}
