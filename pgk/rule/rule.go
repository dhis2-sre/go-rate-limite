package rule

import (
	"net/http"
	"regexp"

	"github.com/dhis2-sre/go-rate-limite/pgk/config"
	"github.com/dhis2-sre/go-rate-limite/pgk/proxy"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
)

type Rule struct {
	config.Rule
	Handler http.Handler
}

func (r *Rule) pathMatch(path string) bool {
	match, err := regexp.MatchString(r.PathPattern, path)
	if err != nil {
		return false
	}
	return match
}

func NewRules(c *config.Config) *Rules {
	var rules []*Rule
	for _, rule := range c.Rules {
		lmt := newLimiter(rule)

		rules = append(rules, &Rule{
			Rule:    rule,
			Handler: tollbooth.LimitFuncHandler(lmt, proxy.Transparently(c.Backend).ServeHTTP),
		})
	}

	return &Rules{
		Rules: rules,
	}
}

func newLimiter(rule config.Rule) *limiter.Limiter {
	lmt := tollbooth.NewLimiter(rule.RequestPerSecond, nil)
	lmt.SetMethods([]string{rule.Method})
	lmt.SetBurst(rule.Burst)
	return lmt
}

type Rules struct {
	Rules []*Rule
}

func (r Rules) Match(req *http.Request) (bool, http.Handler) {
	for _, rule := range r.Rules {
		if rule.pathMatch(req.URL.Path) {
			return true, rule.Handler
		}
	}
	return false, nil
}
