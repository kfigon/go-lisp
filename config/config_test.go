package config

import (
	"go-lisp/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurations(t *testing.T) {
	testCases := []struct {
		desc        string
		code        string
		assertionFn func(*testing.T, *ConfigurationStore)
	}{
		{
			desc: "basic strings",
			code: `(set port 8080)
	(set host "localhost")`,
			assertionFn: func(t *testing.T, c *ConfigurationStore) {
				portV, err := c.Get("port")
				assert.NoError(t, err)
				assert.Equal(t, models.Number(8080), portV)

				hostV, err := c.Get("host")
				assert.NoError(t, err)
				assert.Equal(t, models.String("localhost"), hostV)
			},
		},
		{
			desc: "feature flag with param",
			code: `(lambda flag (market_id) (
				(if (= market_id 123) true false)))`,
			assertionFn: func(t *testing.T, c *ConfigurationStore) {
				v, _ := c.Get("flag", models.Number(123))
				bv := v.(models.Bool)
				assert.True(t, bool(bv))

				v, _ = c.Get("flag", models.Number(1))
				bv = v.(models.Bool)
				assert.False(t, bool(bv))
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c, err := New(tC.code)
			assert.NoError(t, err)
			tC.assertionFn(t, c)
		})
	}
}
