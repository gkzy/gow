/*
like mux router
use regexp
sam
*/

package gow

import (
	"fmt"
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

var rpm RouterPath

// getMuxValue Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getMuxValue(path string, params *Params, unescape bool) (value nodeValue) {
	if rpm == nil {
		rpm = getNodeRouterPathMap(n)
	}
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

	prefix := n.path
	if path == prefix {
		// We should have reached the node containing the handle.
		// Check if this node has a handle registered.
		if value.handlers = n.handlers; value.handlers != nil {
			value.fullPath = n.fullPath
			return
		}

		// If there is no handle for this route, but this route has a
		// wildcard child, there must be a handle for this path with an
		// additional trailing slash
		if path == "/" && n.wildChild && n.nType != root {
			value.tsr = true
			return
		}

		// No handle found. Check if a handle for this path + a
		// trailing slash exists for trailing slash recommendation
		for i, c := range []byte(n.indices) {
			if c == '/' {
				n = n.children[i]
				value.tsr = (len(n.path) == 1 && n.handlers != nil) ||
					(n.nType == catchAll && n.children[0].handlers != nil)
				return
			}
		}

		return
	}

	// Nothing found. We can recommend to redirect to the same URL with an
	// extra trailing slash if a leaf exists for that path
	value.tsr = (path == "/") ||
		(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
			path == prefix[:len(prefix)-1] && n.handlers != nil)

	return
}

var (
	intRegexp  = []byte(`(\d+)`)
	charRegexp = []byte(`(\w+)`)
	starRegexp = []byte(`(.*)`)
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
		// static /static/*filepath
		if strings.Contains(n, "*filepath") {
			n = strings.ReplaceAll(n, "*filepath", string(starRegexp))
			replaceRegexp = starRegexp
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

		fmt.Println("path:", path)
		fmt.Println("p.Path:", p.Path)
		fmt.Println("regPath:", regPath)
		fmt.Println("=====================================")

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
