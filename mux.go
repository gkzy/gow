/*
like mux router
use regexp
sam
*/

package gow

import (
	"net/url"
	"regexp"
	"strings"
)

type routerPathInfo struct {
	Path     string
	fullPath string
	handlers HandlersChain
	params   *Params
	tsr      bool
}

// RoutesInfo defines a routerPathInfo array.
type RouterPath []routerPathInfo

func (n *node) getMuxValue(path string, params *Params, unescape bool) (value nodeValue) {
	rpm := getNodeRouterPathMap(n)
	routerPath, ok := getMatchPath(path, rpm, unescape)
	if ok {
		if params != nil {
			value.params = params
			value.params = routerPath.params
		}
		value.handlers = routerPath.handlers
		value.fullPath = routerPath.fullPath
		return
	}
	return
}

var (
	intRegexp  = []byte(`(\d+)`)
	charRegexp = []byte(`(\w+)`)
)

func mathPath(path string) (regPath string, keys []string) {
	var (
		replaceRegexp []byte
		nPath         string
	)
	nSplit := strings.Split(path, "/")
	// like {uid} or {uid:int}
	wildcardRegexp := regexp.MustCompile(`{\w+}`)

	// replace {uid} to regexp
	replaceRegexp = charRegexp

	for _, n := range nSplit {
		// math {uid:int}
		if strings.Contains(n, ":int") {
			n = strings.ReplaceAll(n, ":int", "")
			replaceRegexp = intRegexp
		}
		key := wildcardRegexp.FindAllString(n, -1)
		keys = append(keys, key...)
		nPath = string(wildcardRegexp.ReplaceAll([]byte(n), replaceRegexp))
		regPath = regPath + nPath + "/"
	}
	regPath = regPath[:len(regPath)-1]
	return regPath, keys
}

// getMatchPath return routerPathInfo
// regexp match
func getMatchPath(path string, rp RouterPath, unescape bool) (*routerPathInfo, bool) {
	for _, p := range rp {
		regPath, keys := mathPath(p.Path)

		// all match
		ok, _ := regexp.MatchString("^"+regPath+"$", path)
		if ok {
			valueRegexp := regexp.MustCompile(regPath)
			if unescape {
				if v, err := url.QueryUnescape(path); err == nil {
					path = v
				}
			}
			values := valueRegexp.FindStringSubmatch(path)
			params := new(Params)
			for i, k := range keys {
				*params = append(*params, Param{
					Key:   strings.ReplaceAll(strings.ReplaceAll(k, "{", ""), "}", ""),
					Value: values[i+1],
				})
			}
			p.params = params
			return &p, ok
		}
	}
	return nil, false
}

func getNodeRouterPathMap(n *node) (rp RouterPath) {
	rp = getNodeRouterPath("", rp, n)
	return
}

func getNodeRouterPath(path string, rp RouterPath, root *node) RouterPath {
	path += root.path
	if len(root.handlers) > 0 {
		rp = append(rp, routerPathInfo{
			Path:     strings.ToLower(path),
			fullPath: root.fullPath,
			handlers: root.handlers,
		})
	}
	for _, child := range root.children {
		rp = getNodeRouterPath(path, rp, child)
	}
	return rp
}
