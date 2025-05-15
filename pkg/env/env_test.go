package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDev(t *testing.T) {
	testCases := []struct {
		name     string
		env      *Env
		expected bool
	}{
		{
			name:     "Dev environment",
			env:      &Env{Dev},
			expected: true,
		},
		{
			name:     "Non-Dev environment",
			env:      &Env{Local},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.env.IsDev()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsLocal(t *testing.T) {
	testCases := []struct {
		name     string
		env      *Env
		expected bool
	}{
		{
			name:     "Local environment",
			env:      &Env{Local},
			expected: true,
		},
		{
			name:     "Non-Local environment",
			env:      &Env{Dev},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.env.IsLocal()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsProd(t *testing.T) {
	testCases := []struct {
		name     string
		env      *Env
		expected bool
	}{
		{
			name:     "Prod environment",
			env:      &Env{Prod},
			expected: true,
		},
		{
			name:     "Non-Prod environment",
			env:      &Env{Dev},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.env.IsProd()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsStage(t *testing.T) {
	testCases := []struct {
		name     string
		env      *Env
		expected bool
	}{
		{
			name:     "Stage environment",
			env:      &Env{Stage},
			expected: true,
		},
		{
			name:     "Non-Stage environment",
			env:      &Env{Dev},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.env.IsStage()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsTest(t *testing.T) {
	testCases := []struct {
		name     string
		env      *Env
		expected bool
	}{
		{
			name:     "Test environment",
			env:      &Env{Test},
			expected: true,
		},
		{
			name:     "Non-Test environment",
			env:      &Env{Dev},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.env.IsTest()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFromString_ValidValues(t *testing.T) {
	testCases := []struct {
		name   string
		envStr string
	}{
		{
			name:   "dev environment",
			envStr: "dev",
		},
		{
			name:   "local environment",
			envStr: "local",
		},
		{
			name:   "prod environment",
			envStr: "prod",
		},
		{
			name:   "stage environment",
			envStr: "stage",
		},
		{
			name:   "test environment",
			envStr: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env, err := FromString(tc.envStr)
			assert.NoError(t, err)
			assert.NotNil(t, env)

			// Check that the correct environment type is set
			switch tc.envStr {
			case "dev":
				assert.True(t, env.IsDev())
				assert.False(t, env.IsLocal())
				assert.False(t, env.IsProd())
				assert.False(t, env.IsStage())
				assert.False(t, env.IsTest())
			case "local":
				assert.False(t, env.IsDev())
				assert.True(t, env.IsLocal())
				assert.False(t, env.IsProd())
				assert.False(t, env.IsStage())
				assert.False(t, env.IsTest())
			case "prod":
				assert.False(t, env.IsDev())
				assert.False(t, env.IsLocal())
				assert.True(t, env.IsProd())
				assert.False(t, env.IsStage())
				assert.False(t, env.IsTest())
			case "stage":
				assert.False(t, env.IsDev())
				assert.False(t, env.IsLocal())
				assert.False(t, env.IsProd())
				assert.True(t, env.IsStage())
				assert.False(t, env.IsTest())
			case "test":
				assert.False(t, env.IsDev())
				assert.False(t, env.IsLocal())
				assert.False(t, env.IsProd())
				assert.False(t, env.IsStage())
				assert.True(t, env.IsTest())
			}
		})
	}
}

func TestFromString_InvalidValue(t *testing.T) {
	invalidEnvStr := "invalid"
	env, err := FromString(invalidEnvStr)

	assert.Error(t, err)
	assert.Nil(t, env)
	assert.Contains(t, err.Error(), invalidEnvStr)
}
