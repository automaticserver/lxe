package cli

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var (
	_, cmd = New(Options{
		Type: TypeTool,
	})
)

func init() {
	cmd.Use = "prog"
	cmd.Short = "prog has a short description"
	cmd.Long = "prog has a very long description. Phasellus tincidunt ac leo a malesuada. Suspendisse posuere libero sit amet augue facilisis consectetur. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Sed sit amet tellus dolor. Nam nec eros in orci molestie efficitur non id justo." + "\n\n" + "Sed mollis, lorem nec consectetur egestas, libero nisl gravida ex, non semper turpis justo lacinia mauris. Sed consectetur tellus nisi, in bibendum purus pulvinar ac. Curabitur maximus massa at est dignissim, quis bibendum augue tempor."
	cmd.Example = "prog -s foo -r bar -k baz"

	pflags := cmd.PersistentFlags()

	pflags.StringP("short", "s", "", "A pretty normal short flag. Except this usage description is made exceptionally long so it should word-wrap in configuration files, depending on if they are told to do so. Do usage flags have an ending punctuation or not?")
	pflags.StringP("remote.first", "", "", "A flag which is in in a subtree")
	pflags.StringP("remote.second", "", "", "The other part of the subtree flag so we can see what this means")
	pflags.StringP("store-dir.dir", "S", "store", "A flag which has a dash and a subtree. The dash should is part of the main key, and not a delimititer for the subtree")
	pflags.StringP("store-dir.log-level", "L", "debug", "The other subtree element has a dash as well")
	pflags.StringP("store-dir.another.sub-level", "", "foo", "A flag with a second sublevel, sometimes with dashes")
	pflags.Bool("abool", false, "A bool flag")
	pflags.BytesBase64("abytes", []byte{'\n'}, "A bytes base64 flag")
	pflags.Duration("aduration", 30*time.Second, "A duration flag")
	pflags.Float64("afloat", 3.1415, "A float flag. THIS IS A VIPER BUG! Gets transformed to a string!")
	pflags.IP("anip", net.ParseIP("127.0.0.1"), "An IP flag")
	_, ipnet, _ := net.ParseCIDR("127.0.224.225/24")
	pflags.IPNet("anipnet", *ipnet, "An IPNet flag")
	pflags.Int("anint", 42, "An int flag")
	pflags.StringSlice("astringslice", []string{"o h", "sole", "mio"}, "A string slice flag")

	// Bind pflags as late as possible so all imports were able to set their flags
	err := venom.BindPFlags(pflags)
	if err != nil {
		panic(err)
	}
}

func runC(t *testing.T, c *cobra.Command, a []string, w io.Writer) {
	for c != rootCmd {
		a = append([]string{c.Name()}, a...)
		c = c.Parent()
	}

	cmd.SetArgs(a)
	cmd.SetOutput(w)

	_, err := cmd.ExecuteC()
	assert.NoError(t, err)
}

func compareGoldenFile(t *testing.T, c *cobra.Command, a []string, g string) {
	exp, err := ioutil.ReadFile(g)
	assert.NoError(t, err)

	act := &bytes.Buffer{}

	runC(t, c, a, act)

	assert.Equal(t, string(exp), act.String())
}

func Test_Is(t *testing.T) {
	assert.True(t, Is(rootCmd, TypeTool))
	assert.False(t, Is(rootCmd, TypeService))
}
